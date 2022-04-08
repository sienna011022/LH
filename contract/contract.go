/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}
type ContractProcess struct {
	Key          string `json:"key"`
	Name_A       string `json:"name_A"`
	Id_number_A  uint64 `json:"id_number_A"`
	Name_B       string `json:"name_B"`
	Id_number_B  uint64 `json:"id_number_B"`
	State        string `json:state`
	Contract_id  uint64 `json:"contract_id"`
	ContractHash string `json:"contractHash"`
}

//fuction
func (s *SmartContract) InitContract(ctx contractapi.TransactionContextInterface, key string) error {

	//marshal
	var contract = ContractProcess{Key: key}
	contractAsBytes, _ := json.Marshal(contract)
	return ctx.GetStub().PutState(key, contractAsBytes)

}

func (s *SmartContract) UpdateState(ctx contractapi.TransactionContextInterface, key string, newstate string) error {

	contractAsBytes, err := ctx.GetStub().GetState(key)

	if err != nil {
		return err
	} else if contractAsBytes == nil {
		return fmt.Errorf("User does not exist " + key + "/")
	}

	contract := ContractProcess{}
	err = json.Unmarshal(contractAsBytes, &contract)

	if err != nil {
		return err
	}

	contract.State = newstate

	contractAsBytes, err = json.Marshal(contract)

	return ctx.GetStub().PutState(key, contractAsBytes)

	if err != nil {
		return fmt.Errorf("failed to Marshaling:%v", err)

	}

	err = ctx.GetStub().PutState(key, contractAsBytes)

	if err != nil {
		return fmt.Errorf("failed to AddContract %v", err)

	}

	return nil
}

func (s *SmartContract) UpdateContract(ctx contractapi.TransactionContextInterface, key string, name_A string, name_B string, id_number_A uint64, id_number_B uint64, state string, contract_id uint64, contractHash string) error {
	contractAsBytes, err := ctx.GetStub().GetState(key)

	if err != nil {
		return err
	} else if contractAsBytes == nil {
		return fmt.Errorf("User does not exist " + key + "/")
	}

	contract := ContractProcess{}
	err = json.Unmarshal(contractAsBytes, &contract)

	if err != nil {
		return err
	}

	//contract_id64, _ := strconv.ParseInt(contract_id, 10, 64)

	contract.Key = key
	contract.Name_A = name_A
	contract.Name_B = name_B
	contract.Id_number_A = id_number_A
	contract.Id_number_B = id_number_B
	contract.Contract_id = contract_id
	contract.ContractHash = contractHash

	contractAsBytes, err = json.Marshal(contract)

	if err != nil {
		return fmt.Errorf("failed to Marshaling:%v", err)

	}

	err = ctx.GetStub().PutState(key, contractAsBytes)

	if err != nil {
		return fmt.Errorf("failed to AddContract %v", err)

	}

	return nil

}
func (s *SmartContract) ReadContract(ctx contractapi.TransactionContextInterface, key string) (string, error) {
	//get value from ctx
	contractAsBytes, err := ctx.GetStub().GetState(key)

	if err != nil {
		return "", fmt.Errorf("failed to read from world state,%s", err.Error())
	}

	if contractAsBytes == nil {
		return "", fmt.Errorf("%s  does not exist", key)

	}

	return string(contractAsBytes[:]), nil

}

func (s *SmartContract) ContractState(ctx contractapi.TransactionContextInterface, Key string) (string, error) {
	//get value from ctx
	contractAsBytes, err := ctx.GetStub().GetState(Key)

	if err != nil {
		return "", fmt.Errorf("failed to read from world state,%s", err.Error())
	}

	if contractAsBytes == nil {
		return "", fmt.Errorf("%s  does not exist", Key)

	}

	return string(contractAsBytes[:]), nil

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
