/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}
type UserRating struct{
	User string `json:"user"`
	Average float64 `json:"average"`
	Rates []Rate `json:"rates"`

}

type Rate struct{
	projectTitle string `json:"projecttitle"`
	Score float64 `json:"score"`
}

//fuction

func (s *SmartContract)AddUser(ctx contractapi.TransactionContextInterface,username string) error {
	//error check



	
	//marshal
	var user = UserRating{User:username,Average : 0 }
	userAsBytes,_ := json.Marshal(user)
	return  ctx.GetStub().PutState(username,userAsBytes)


}
func (s *SmartContract)AddRating(ctx contractapi.TransactionContextInterface,username string,prjTitle string,prjscore string) error {
	
  userAsBytes,err:= ctx.GetStub().GetState(username)


  if err != nil{
	  return err
  }else if userAsBytes == nil{
	  return fmt.Errorf("User does not exist "+username+"/")
  }

  user  := UserRating{}
  err= json.Unmarshal(userAsBytes,&user)
  if err != nil{
	  return err
  }

  newRate ,_ := strconv.ParseFloat(prjscore,64)
  var Rate = Rate{projectTitle:prjTitle,Score : newRate}

  rateCount := float64(len(user.Rates))

  user.Rates= append(user.Rates,Rate)

  user.Average = (rateCount * user.Average + newRate)/(rateCount+1)
  userAsBytes,err = json.Marshal(user);

  if err !=nil{
	  return fmt.Errorf("failed to Marshaling:%v",err)

  }

  err = ctx.GetStub().PutState(username,userAsBytes)

  if err != nil{
	  return fmt.Errorf("failed to AddRating %v",err)

  }

  return nil




}
func (s *SmartContract)ReadRating(ctx contractapi.TransactionContextInterface,username string) (string,error) {
	//get value from ctx
	userAsBytes, err := ctx.GetStub().GetState(username)

	if(err != nil){
		return "",fmt.Errorf("failed to read from world state,%s", err.Error())
	}
	
	if (userAsBytes == nil) {
		return "",fmt.Errorf("%s  does not exist",username)

	}

	return string(userAsBytes[:]),nil





}

// Car describes basic details of what makes up a car
type Car struct {
	Make   string `json:"make"`
	Model  string `json:"model"`
	Colour string `json:"colour"`
	Owner  string `json:"owner"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Car
}

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	cars := []Car{
		Car{Make: "Toyota", Model: "Prius", Colour: "blue", Owner: "Tomoko"},
		Car{Make: "Ford", Model: "Mustang", Colour: "red", Owner: "Brad"},
		Car{Make: "Hyundai", Model: "Tucson", Colour: "green", Owner: "Jin Soo"},
		Car{Make: "Volkswagen", Model: "Passat", Colour: "yellow", Owner: "Max"},
		Car{Make: "Tesla", Model: "S", Colour: "black", Owner: "Adriana"},
		Car{Make: "Peugeot", Model: "205", Colour: "purple", Owner: "Michel"},
		Car{Make: "Chery", Model: "S22L", Colour: "white", Owner: "Aarav"},
		Car{Make: "Fiat", Model: "Punto", Colour: "violet", Owner: "Pari"},
		Car{Make: "Tata", Model: "Nano", Colour: "indigo", Owner: "Valeria"},
		Car{Make: "Holden", Model: "Barina", Colour: "brown", Owner: "Shotaro"},
	}

	for i, car := range cars {
		carAsBytes, _ := json.Marshal(car)
		err := ctx.GetStub().PutState("CAR"+strconv.Itoa(i), carAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateCar adds a new car to the world state with given details
func (s *SmartContract) CreateCar(ctx contractapi.TransactionContextInterface, carNumber string, make string, model string, colour string, owner string) error {
	car := Car{
		Make:   make,
		Model:  model,
		Colour: colour,
		Owner:  owner,
	}

	carAsBytes, _ := json.Marshal(car)

	return ctx.GetStub().PutState(carNumber, carAsBytes)
}

// QueryCar returns the car stored in the world state with given id
func (s *SmartContract) QueryCar(ctx contractapi.TransactionContextInterface, carNumber string) (*Car, error) {
	carAsBytes, err := ctx.GetStub().GetState(carNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if carAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", carNumber)
	}

	car := new(Car)
	_ = json.Unmarshal(carAsBytes, car)

	return car, nil
}

// QueryAllCars returns all cars found in world state
func (s *SmartContract) QueryAllCars(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		car := new(Car)
		_ = json.Unmarshal(queryResponse.Value, car)

		queryResult := QueryResult{Key: queryResponse.Key, Record: car}
		results = append(results, queryResult)
	}

	return results, nil
}

// ChangeCarOwner updates the owner field of car with given id in world state
func (s *SmartContract) ChangeCarOwner(ctx contractapi.TransactionContextInterface, carNumber string, newOwner string) error {
	car, err := s.QueryCar(ctx, carNumber)

	if err != nil {
		return err
	}

	car.Owner = newOwner

	carAsBytes, _ := json.Marshal(car)

	return ctx.GetStub().PutState(carNumber, carAsBytes)
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}