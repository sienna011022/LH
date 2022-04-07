/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

type RightProcess struct {
	key       string     `json:"key"`
	name      string     `json:"name"`
	state     string     `json:"state"`
	Contracts []Contract `json:"Contracts"`
}

type Contract struct {
	docu_id       uint64    `json:"docu_id"`
	docu_name     string    `json:"docu_name"`
	document_hash string    `json:"document_hash"`
	timestamp     time.Time `json:"timestamp"`
}

//fuction
func (s *SmartContract) AddUser(ctx contractapi.TransactionContextInterface, key string) error {

	//marshal
	var request = RightProcess{key: key}
	requestAsBytes, _ := json.Marshal(request)
	return ctx.GetStub().PutState(key, requestAsBytes)

}

func (s *SmartContract) AddContract(ctx contractapi.TransactionContextInterface, key string, name string, state string, docu_id string, docu_name string, document_hash string) error {
	requestAsBytes, err := ctx.GetStub().GetState(key)

	if err != nil {
		return err
	} else if requestAsBytes == nil {
		return fmt.Errorf("User does not exist " + key + "/")
	}

	request := RightProcess{}
	err = json.Unmarshal(requestAsBytes, &request)

	if err != nil {
		return err
	}

	time := time.Now()

	docu_id64, _ := strconv.ParseInt(docu_id, 10, 64)

	Contract := Contract{docu_id: uint64(docu_id64), docu_name: docu_name, document_hash: document_hash, timestamp: time}

	request.Contracts = append(request.Contracts, Contract)

	requestAsBytes, err = json.Marshal(request)

	if err != nil {
		return fmt.Errorf("failed to Marshaling:%v", err)

	}

	err = ctx.GetStub().PutState(key, requestAsBytes)

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
