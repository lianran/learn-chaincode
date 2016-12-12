/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"time"
	"fmt"
	"strconv"
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// myChaincode 
type myChaincode struct {
}
//sp
var sp = "\n"

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(myChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *myChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *myChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {

	case "create":
		if len(args) < 3{
			return nil, errors.New("create operation must include at last three arguments, a new id , the balance and timestamp")
		}
		// get the args
		id := args[0]
		//check the balance format
		balance, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, errors.New("Expecting integer value for asset holding")
		}
		//check the timestamp
		timestamp := args[2]
		ts := time.Now().Unix() 
		tm, err := strconv.ParseInt(timestamp, 10, 64) 
		if err != nil {
			return nil, fmt.Errorf("bad format of the timestamp")
		}
		if tm - ts > 3600 || ts - tm > 3600 {
			return nil, fmt.Errorf("the timestamp is bad one !")
		}

		//check for existence of id
		oldvalue, err := stub.GetState(id)
		if err != nil {
			return nil, fmt.Errorf("create operation failed. Error accessing state(check the existence of id): %s", err)
		}
		if oldvalue != nil {
			return nil, fmt.Errorf("existed id!")
		} 

		//creat the user
		key := id
		fmt.Printf("the new id is %s and the balance is %s", key, strconv.Itoa(balance))

		err = stub.PutState(key, []byte(strconv.Itoa(balance)))
		if err != nil {
			fmt.Printf("Error putting state %s", err)
			return nil, fmt.Errorf("create operation failed. Error updating state: %s", err)
		}

		//store history
		key = "R" + sp + id + sp + timestamp
		value := "admin" + sp + strconv.Itoa(balance)
		err = stub.PutState(key, []byte(value))
		if err != nil {
			fmt.Printf("Error putting state for fromid : %s", err)
		}
		return nil,nil

	case "transfer":
		if len(args) < 4{
			return nil, errors.New("transfer operation must include at last four arguments, a outid , a toid ,  amount of money and timestamp")
		}
		//get the args
		A := args[0]
		B := args[1]
		amount := args[2]
		timestamp := args[3]

		//check the args
		//check the amount format
		toval, err := strconv.Atoi(amount)
		if err != nil {
			return nil, errors.New("Expecting integer value for asset holding")
		}
		//check the timestamp
		ts := time.Now().Unix() 
		tm, err := strconv.ParseInt(timestamp, 10, 64) 
		if err != nil {
			return nil, fmt.Errorf("bad format of the timestamp")
		}
		if tm - ts > 3600 || ts - tm > 3600 {
			return nil, fmt.Errorf("the timestamp is bad one !")
		}

		//get the balance of outid
		Aval, err := stub.GetState(A)
		if err != nil {
			return nil, fmt.Errorf("get operation failed. Error accessing state: %s", err)
		}
		if Aval == nil {
			return nil, fmt.Errorf(" the user (inid: %s) not exists!", A)
		}

		//get the balance of inid
		Bval, err := stub.GetState(B)
		if err != nil {
			return nil, fmt.Errorf("get opereation failed. Error accessing state: %s", err)
		}
		if Bval == nil {
			return nil, fmt.Errorf(" the user (outid: %s) not exists", B)
		}

		// perform the execution
		Abal, _ := strconv.Atoi(string(Aval))
		Bbal, _ := strconv.Atoi(string(Bval))
		Abal = Abal - toval
		Bbal = Bbal + toval

		//Write the state back to the ledger
		err = stub.PutState(A, []byte(strconv.Itoa(Abal)))
		if err != nil {
			return nil, err
		}
		err = stub.PutState(B, []byte(strconv.Itoa(Bbal)))

		//record the history to db
		key := "S" + sp + A + sp + timestamp
		value := B + sp + amount
		err = stub.PutState(key, []byte(value))
		if err != nil {
			return nil, err
		}
		key = "R" + sp + B + sp +timestamp
		value = A + sp + amount
		err = stub.PutState(key, []byte(value))
		if err != nil {
			return nil, err
		}
		return nil, nil
	default:
		return nil, errors.New("Unsupported operation")
	}
}

// Query is our entry point for queries
func (t *myChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	switch function {

	case "history":
		if len(args) < 1{
			return nil, errors.New("myhistory operation must include at last one arguments, id (and time)")
		}
		id := args[0]

		//get the timestamp
		ts := time.Now().Unix() 
		//timestamp := strconv.FormatInt(ts, 10)

		tm := int64(3600)
		var err error
		if len(args) >= 2{
			tm, err = strconv.ParseInt(args[1], 10, 64)
		}
		if err != nil {
			return nil, fmt.Errorf("get history failed. Bad format of the time: %s", err)
		}

		//set the query of the time
		starttime := strconv.FormatInt(ts-tm, 10)
		endtime := strconv.FormatInt(ts,10)

		//do the range query for the recevied
		keysIter, err := stub.RangeQueryState("R" + sp + id + sp + starttime, "R" + sp + id + sp + endtime)
		if err != nil {
			return nil, fmt.Errorf("get history failed. Error accessing state: %s", err)
		}
		defer keysIter.Close()

		var keys []string
		for keysIter.HasNext() {
			key, _, iterErr := keysIter.Next()
			if iterErr != nil {
				return nil, fmt.Errorf("get history operation failed. Error accessing state: %s", err)
			}
			keys = append(keys, key)
		}
		//do the range query for the sent
		keysIter, err = stub.RangeQueryState("S" + sp + id + sp + starttime, "S" + sp + id + sp + endtime)
		if err != nil {
			return nil, fmt.Errorf("get history failed. Error accessing state: %s", err)
		}
		defer keysIter.Close()

		for keysIter.HasNext() {
			key, _, iterErr := keysIter.Next()
			if iterErr != nil {
				return nil, fmt.Errorf("get history operation failed. Error accessing state: %s", err)
			}
			keys = append(keys, key)
		}
		
		jsonKeys, err := json.Marshal(keys)
		if err != nil {
			return nil, fmt.Errorf("keys operation failed. Error marshaling JSON: %s", err)
		}

		return jsonKeys, nil
	case "getbalance":
		if len(args) < 1{
			return nil, errors.New("getnumsofbills operation must include at last one argument, id")
		}
		id := args[0]
		//Todo: some check for the owner?

		val, err := stub.GetState(id)
		if err != nil {
			return nil, fmt.Errorf("get balance failed. Error accessing state: %s", err)
		}
		if val == nil {
			return nil, fmt.Errorf(" the id not exists!")
		}

		return val, nil

	default:
		return nil, errors.New("Unsupported operation")
	}
}
