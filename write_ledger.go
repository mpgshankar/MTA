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
	"time"

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
	var key, theatreRegNo, releaseDate, value string
	// var err error
	fmt.Println("starting add_movies")

	certname, err := get_cert(stub)
	if err != nil {
		fmt.Printf("INVOKE: Error retrieving cert: %s", err)
		return shim.Error("Error retrieving cert")
	}

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting Minimum 1. arguments of the variable and value to set")
	}

	// value = strings.Replace(args[0], "\"", "", -1) //rename for funsies
	value = args[0]
	var jsonValue map[string]interface{}
	json.Unmarshal([]byte(value), &jsonValue)
	key, _ = jsonValue["movieId"].(string)
	theatreRegNo = string(certname)
	movieName, _ := jsonValue["movieName"].(string)
	movieDuration, _ := jsonValue["movieDuration"].(string)
	releaseDate, _ = jsonValue["movieReleaseDate"].(string)
	showTimings, _ := jsonValue["showTimings"].([]interface{})

	// Create Movie Object
	var mov Movies
	mov.ObjectType = "Movies"
	mov.MovieId = key
	mov.MovieName = movieName
	mov.MovieDuration = movieDuration
	mov.MovieReleaseDate = releaseDate

	// Add Show Timings in Movie Struct

	for _, show := range showTimings {
		shTm, _ := show.(string)
		var shw Shows
		shw.ShowTiming = shTm
		shw.TotalSeat = 100
		shw.AvailableSeat = 100
		shw.BookedSeat = 0
		mov.ShowTimings = append(mov.ShowTimings, shw)
	}

	// Get Current Date Time
	loc, _ := time.LoadLocation("Asia/Kolkata")
	currentDate := time.Now().In(loc)
	format := "2006-01-02"
	curDate, _ := time.Parse(format, currentDate.Format("2006-01-02"))
	release, _ := time.Parse(format, releaseDate)

	//check if theatre exists or not
	theatreAsBytes, _ := stub.GetState(theatreRegNo)
	if theatreAsBytes == nil {
		fmt.Println("This theatre does not exists - " + theatreRegNo)
		return shim.Error("This theatre does not exists - " + theatreRegNo)
	}
	theatre := Theatre{}
	json.Unmarshal(theatreAsBytes, &theatre) //un stringify it aka JSON.parse()

	// check movies when it will be releasing
	if greaterThanEqualCurrentDate(release, curDate) {
		if equalCurrentDate(release, curDate) {
			mov.Status = "Running"
			theatre.MoviesRunning = append(theatre.MoviesRunning, mov)
		} else {
			mov.Status = "Coming Soon"
			theatre.MoviesComingSoon = append(theatre.MoviesComingSoon, mov)
		}
	}

	trAsBytes, _ := json.Marshal(theatre)

	errTr := stub.PutState(theatreRegNo, trAsBytes) // update the theatre details into the ledger
	if errTr != nil {
		return shim.Error("Failed to add movies : " + errTr.Error())
	}

	valueAsBytes, _ := json.Marshal(mov)

	errPut := stub.PutState(key, valueAsBytes) //write the movie details into the ledger
	if errPut != nil {
		return shim.Error("Failed to add movies : " + errPut.Error())
	}

	fmt.Println("- end add_movies")
	return shim.Success(nil)
}

func greaterThanEqualCurrentDate(start, check time.Time) bool {
	return start.After(check) || start.Equal(check)
}

func equalCurrentDate(start, check time.Time) bool {
	return start.Equal(check)
}
