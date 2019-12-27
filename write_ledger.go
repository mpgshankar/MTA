/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding  ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at
  http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ============================================================================================================================
// invoke_transaction_insert_update() - genric insert json object into ledger
//
// Shows Off PutState() - writting a key/value into the ledger
//
// Inputs - JSON Object
//    0
//   json_object
//  {"key1":"value1","key2":"value2","key3":"value3"}
// ============================================================================================================================
func invoke_transaction_insert_update(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, value string
	// var err error
	fmt.Println("starting invoke_transaction_insert_update")

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting Minimum 1. arguments of the variable and value to set")
	}

	// value = strings.Replace(args[0], "\"", "", -1) //rename for funsies
	value = args[0]
	var jsonValue map[string]interface{}
	json.Unmarshal([]byte(value), &jsonValue)
	key, _ = jsonValue["transactionGroupId"].(string)

	valueAsBytes, _ := json.Marshal(jsonValue)

	errPut := stub.PutState(key, valueAsBytes) //write the transaction into the ledger
	if errPut != nil {
		return shim.Error("Failed to put state : " + errPut.Error())
	}

	fmt.Println("- end invoke_transaction_insert_update")
	return shim.Success(nil)
}

// ============================================================================================================================
// add_theatre() - add theatre into ledger
//
// Shows Off PutState() - writting a key/value into the ledger
//
// Inputs - JSON Object
//    0
//   json_object
//  {"theatreRegNo":"value1","theatreLocation":"value2"}
// ============================================================================================================================
func add_theatre(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, value string
	// var err error
	fmt.Println("starting add_theatre")

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting Minimum 1. arguments of the variable and value to set")
	}

	// value = strings.Replace(args[0], "\"", "", -1) //rename for funsies
	value = args[0]
	var jsonValue map[string]interface{}
	json.Unmarshal([]byte(value), &jsonValue)
	key, _ = jsonValue["theatreRegNo"].(string)

	valueAsBytes, _ := json.Marshal(jsonValue)

	errPut := stub.PutState(key, valueAsBytes) //write the theatre details into the ledger
	if errPut != nil {
		return shim.Error("Failed to add theatre : " + errPut.Error())
	}

	fmt.Println("- end add_theatre")
	return shim.Success(nil)
}

// ============================================================================================================================
// add_movies() - add movie into ledger
//
// Shows Off PutState() - writting a key/value into the ledger
//
// Inputs - JSON Object
//    0
//   json_object
//  {"theatreRegNo":"value1","theatreLocation":"value2"}
// ============================================================================================================================
func add_movies(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, theatreRegNo, value string
	// var err error
	fmt.Println("starting add_theatre")

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting Minimum 1. arguments of the variable and value to set")
	}

	// value = strings.Replace(args[0], "\"", "", -1) //rename for funsies
	value = args[0]
	var jsonValue map[string]interface{}
	json.Unmarshal([]byte(value), &jsonValue)
	key, _ = jsonValue["movieId"].(string)
	theatreRegNo, _ = jsonValue["theatreRegNo"].(string)

	theatre, err := stub.GetState(theatreRegNo)

	valueAsBytes, _ := json.Marshal(jsonValue)

	errPut := stub.PutState(key, valueAsBytes) //write the theatre details into the ledger
	if errPut != nil {
		return shim.Error("Failed to add theatre : " + errPut.Error())
	}

	fmt.Println("- end add_theatre")
	return shim.Success(nil)
}
