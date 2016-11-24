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
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	err := createTable(stub)
	if err != nil {
		return nil, fmt.Errorf("Error creating table during init. %s", err)
	}

	return nil, nil
}
func createTable(stub shim.ChaincodeStubInterface) error {
	// Create table one
	var columnDefsTableOne []*shim.ColumnDefinition
	columnOneTableOneDef := shim.ColumnDefinition{Name: "id",
		Type: shim.ColumnDefinition_STRING, Key: true}
	columnTwoTableOneDef := shim.ColumnDefinition{Name: "productor",
		Type: shim.ColumnDefinition_STRING, Key: false}
	columnThreeTableOneDef := shim.ColumnDefinition{Name: "time",
		Type: shim.ColumnDefinition_STRING, Key: false}
	columnDefsTableOne = append(columnDefsTableOne, &columnOneTableOneDef)
	columnDefsTableOne = append(columnDefsTableOne, &columnTwoTableOneDef)
	columnDefsTableOne = append(columnDefsTableOne, &columnThreeTableOneDef)
	return stub.CreateTable("table", columnDefsTableOne)
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "insert" {
		if len(args) < 3 {
			return nil, errors.New("insert failed. Must include 3 column values")
		}

		col1Val := args[0]
		col2Val := args[1]
		col3Val := args[2]

		var columns []*shim.Column
		col1 := shim.Column{Value: &shim.Column_String_{String_: col1Val}}
		col2 := shim.Column{Value: &shim.Column_String_{String_: col2Val}}
		col3 := shim.Column{Value: &shim.Column_String_{String_: col3Val}}
		columns = append(columns, &col1)
		columns = append(columns, &col2)
		columns = append(columns, &col3)

		row := shim.Row{Columns: columns}
		ok, err := stub.InsertRow("table", row)
		if err != nil {
			return nil, fmt.Errorf("table operation failed. %s", err)
		}
		if !ok {
			return nil, errors.New("table operation failed. Row with given key already exists")
		}


	}
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries  
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "getbyid" {
		if len(args) < 1 {
			return nil, errors.New("get failed. Must have the key values")
		}
		var columns []shim.Column

		col1Val := args[0]
		col1 := shim.Column{Value: & shim.Column_String_{String_: col1Val}}
		columns = append(columns, col1)

		rowChannel, err := stub.GetRows("table", columns)
		if err != nil {
			return nil, fmt.Errorf("get operation failed. %s", err)
		}

		var rows []shim.Row
		for {
			select {
			case row, ok := <-rowChannel:
				if !ok {
					rowChannel = nil
				} else {
					rows = append(rows, row)
				}
			}
			if rowChannel == nil {
				break
			}
		}

		jsonRows, err := json.Marshal(rows)
		if err != nil {
			return nil, fmt.Errorf("get operation failed. Error marshaling JSON: %s", err)
		}

		return jsonRows, nil
	}
	fmt.Println("query did not find func: " + function)	

	return nil, errors.New("Received unknown function query: " + function)
}
