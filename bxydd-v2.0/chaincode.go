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
		// put the balance
		err = stub.PutState(key, []byte(strconv.Itoa(balance)))
		if err != nil {
			fmt.Printf("Error putting state %s", err)
			return nil, fmt.Errorf("create operation failed. Error updating state: %s", err)
		}
		// put the num of tx
		key = id + sp + "numoftx"
		err = stub.PutState(key, []byte(strconv.Itoa(1)))
		if err != nil {
			fmt.Printf("Error putting state %s", err)
			return nil, fmt.Errorf("create operation failed. Error updating state: %s", err)
		}

		//store history
		key = id + sp + strconv.Itoa(1)
		value := "R" + sp + "admin" + sp + strconv.Itoa(balance) + timestamp
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
		//get the txs num of A and B
		Anum, err := stub.GetState(A + sp + "numoftx")
		if err != nil {
			return nil, fmt.Errorf("get operation failed. Error accessing state: %s", err)
		}
		Bnum, err := stub.GetState(B + sp + "numoftx")
		if err != nil {
			return nil, fmt.Errorf("get opereation failed. Error accessing state: %s", err)
		}
		A_num, _ := strconv.Atoi(string(Anum))
		B_num, _ := strconv.Atoi(string(Bnum))
		A_num += 1
		B_num += 1
		// put the record to the db
		key := A + sp + strconv.Itoa(A_num)
		value := "S" + sp + B + sp + amount + sp + timestamp
		err = stub.PutState(key, []byte(value))
		if err != nil {
			return nil, err
		}
		err = stub.PutState(A + sp + "numoftx", []byte(strconv.Itoa(A_num)))
		if err != nil {
			return nil, err
		}

		key = B + sp + strconv.Itoa(B_num)
		value = "R" + sp + A + sp + amount + sp + timestamp
		err = stub.PutState(key, []byte(value))
		if err != nil {
			return nil, err
		}
		err = stub.PutState(B + sp + "numoftx", []byte(strconv.Itoa(B_num)))
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
		if len(args) < 2{
			return nil, errors.New("myhistory operation must include at last two arguments, id and a No")
		}
		id := args[0]
		num := args[1]
		//Todo: some check for the num

		val, err := stub.GetState(id + sp + num)
		if err != nil {
			return nil, fmt.Errorf("history failed. Error accessing state: %s", err)
		}
		if val == nil {
			return nil, fmt.Errorf("get nothing!")
		}
		return val, nil

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

	case "gettxnum":
		if len(args) < 1 {
			return nil, errors.New("gettxnum operation must include  at last one argument, id")
		}
		id := args[0]

		val, err := stub.GetState(id + sp + "numoftx")
		if err != nil {
			return nil, fmt.Errorf("get txnum failed. Error accessing state: %s", err)
		}
		if val == nil {
			return nil, fmt.Errorf("the id not exists!")
		}

		return val, nil

	default:
		return nil, errors.New("Unsupported operation")
	}
}
