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
	"fmt"
	"strings"
	"time"
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
			return nil, errors.New("create operation must include at last there arguments, a uuid , a from and a to")
		}
		// get the args
		uuid := args[0]
		fromid := args[1]
		toid := args[2]
		metadata := args[3]
		history := fromid
		owner := fromid

		//TODO: need some check for fromid ?


		//get the timestamp
		ts := time.Now().Unix() 
		timestamp := strconv.FormatInt(ts, 10) 

		key := uuid
		value := fromid + sp + toid + sp + history + sp + metadata + sp + owner
		fmt.Printf("value is %s", value)

		err := stub.PutState(key, []byte(value))
		if err != nil {
			fmt.Printf("Error putting state %s", err)
			return nil, fmt.Errorf("create operation failed. Error updating state: %s", err)
		}
		//store the from and to 
		key = fromid + sp + timestamp + sp + uuid
		err = stub.PutState(key, []byte(timestamp))
		if err != nil {
			fmt.Printf("Error putting state for fromid : %s", err)
		}
		key = toid + sp + timestamp + sp + uuid
		err = stub.PutState(key, []byte(timestamp))
		if err != nil {
			fmt.Printf("Error putting state for toid : %s", err)
		}
		return nil,nil

	case "transfer":
		if len(args) < 3{
			return nil, errors.New("transfer operation must include at last there arguments, a uuid , a owner and a toid")
		}
		//get the args
		key := args[0]
		uuid := key
		_owner := args[1]
		_toid := args[2]

		//get the  info of uuid
		value, err := stub.GetState(key)
		if err != nil {
			return nil, fmt.Errorf("get operation failed. Error accessing state: %s", err)
		}
		listValue := strings.Split(string(value), sp)
		fromid := listValue[0]
		toid := listValue[1]
		history := listValue[2]
		metadata := listValue[3]
		owner := listValue[4]
		
		//ToDo: some check for the owner?
		// the follow is just a example
		if _owner != owner {
			return []byte("don't have the right to transfer the bill"), errors.New("don't have the right to transfer")
			//return nil, errors.New("don't have the right to transfer")
		}

		//get the timestamp
		ts := time.Now().Unix() 
		timestamp := strconv.FormatInt(ts, 10) 

		history = history + "," + _toid
		owner = _toid
		newvalue := fromid + sp + toid + sp + history + sp + metadata + sp + owner
		fmt.Printf("the old value is: %s", value)
		fmt.Printf("the new value is: %s", newvalue)
		err = stub.PutState(key, []byte(newvalue))
		if err != nil {
			fmt.Printf("Error putting state %s", err)
			return nil, fmt.Errorf("transfer operation failed. Error updating state: %s", err)
		}
		//ToDo: some check for the state of puting 
		key = owner + sp + uuid
		err = stub.PutState(key, []byte(timestamp))
		if err != nil {
			fmt.Printf("Error putting state for owner : %s", err)
		}
		return nil, nil
	default:
		return nil, errors.New("Unsupported operation")
	}
}

// Query is our entry point for queries
func (t *myChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	switch function {

	case "myhistory":
		if len(args) < 1{
			return nil, errors.New("myhistory operation must include at last one argument, owner(and time s)")
		}
		/*#owner := args[0]
		time := 3600
		if len(args) >= 2{
			time, err := strconv.Atoi(args[1])
		}
		*/
		// Todo:rangequery?
		return []byte("todo"), nil
	case "getnumofbills":
		if len(args) < 1{
			return nil, errors.New("getnumsofbills operation must include at last one argument, owner(and time s)")
		}
		owner := args[0]
		//Todo: some check for the owner?

		//get the timestamp
		ts := time.Now().Unix() 
		//timestamp := strconv.FormatInt(ts, 10)

		tm := int64(3600)
		if len(args) >= 2{
			tm, _ = strconv.ParseInt(args[1], 10, 64)
		}
		starttime := strconv.FormatInt(ts-tm, 10)
		endtime := strconv.FormatInt(ts,10)

		keysIter, err := stub.RangeQueryState(owner + sp + starttime, owner + sp + endtime)
		if err != nil {
			return nil, fmt.Errorf("getnumofbills failed. Error accessing state: %s", err)
		}
		defer keysIter.Close()

		cnt := 0

		for keysIter.HasNext() {
			_, _, iterErr := keysIter.Next()
			if iterErr != nil {
				return nil, fmt.Errorf("getnumofbills operation failed. Error accessing state: %s", err)
			}
			cnt = cnt + 1
		}
		return []byte("the number is :" + strconv.Itoa(cnt)), nil

	case "getbill":
		if len(args) < 2{
			return nil, errors.New("getbill operation must include at last two arguments, uuid and owner")
		}
		uuid := args[0]
		_owner := args[1]

		//ToDo: some checks?
		key := uuid
		value, err := stub.GetState(key)
		listValue := strings.Split(string(value), sp)
		// check the ownership
		owner := listValue[4]
		if _owner != owner {
			return []byte("you don't have the right to get this bill"), nil
		}
		if err != nil {
			return nil, fmt.Errorf("get operation failed. Error accessing state: %s", err)
		}
		return value, nil
	default:
		return nil, errors.New("Unsupported operation")
	}
}
