/*
Li
censed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding  donorship.  The ASF licenses this file
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
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// MTA Chaincode implementation
type MTA struct {
}

// Theatre Struct
type Theatre struct {
	ObjectType       string   `json:"docType"` // field defined for couchdb
	TheatreRegNo     string   `json:"theatreRegNo"`
	TheatreName      string   `json:"theatreName"`
	TheatreLocation  string   `json:"theatreLocation"`
	MoviesRunning    []Movies `json:"moviesRunning"`
	MoviesComingSoon []Movies `json:"moviesComingSoon"`
}

// Movies Struct
type Movies struct {
	ObjectType       string  `json:"docType"` // field defined for couchdb
	MovieId          string  `json:"movieId"`
	MovieName        string  `json:"movieName"`
	MovieReleaseDate string  `json:"movieReleaseDate"`
	MovieDuration    string  `json:"movieDuration"`
	ShowTimings      []Shows `json:"showTimings"`
	TheatreRegNo     string  `json:"theatreRegNo"`
	Status           string  `json:"status"`
}

type Shows struct {
	ShowTiming    string `json:showTiming`
	TotalSeat     int    `json:totalSeat`
	AvailableSeat int    `json:availableSeat`
	BookedSeat    int    `json:bookedSeat`
}

type Tickets struct {
	ObjectType string `json:"docType"` // field defined for couchdb
	TicketId   string `json:"ticketId"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(MTA))
	if err != nil {
		fmt.Printf("Error starting MTA chaincode - %s", err)
	}
}

// ============================================================================================================================
// Init - initialize the chaincode - MTA don’t need anything initlization, so let's run a dead simple test instead
// ============================================================================================================================
func (t *MTA) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("MTA Is Starting Up")
	_, args := stub.GetFunctionAndParameters()
	var Aval int
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// convert numeric string to integer
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return shim.Error("Expecting a numeric string argument to Init()")
	}

	// store compaitible projects application version
	err = stub.PutState("projects_ui", []byte("3.5.0"))
	if err != nil {
		return shim.Error(err.Error())
	}

	// this is a very simple dumb test.  let's write to the ledger and error on any errors
	err = stub.PutState("selftest", []byte(strconv.Itoa(Aval))) //making a test var "selftest", its handy to read this right away to test the network
	if err != nil {
		return shim.Error(err.Error()) //self-test fail
	}

	fmt.Println(" - ready for action") //self-test pass
	return shim.Success(nil)
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *MTA) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println(" ")
	fmt.Println("starting invoke, for - " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub)
	} else if function == "read" { //generic read ledger
		return read(stub, args)
	} else if function == "generic_query" { //generic_query read from ledger using couchdb
		return t.generic_query(stub, args)
	} else if function == "generic_query_pagination" { //generic_query_pagination read from ledger using couchdb
		return t.generic_query_pagination(stub, args)
	} else if function == "getHistory" { //read history of a key (audit)
		return getHistory(stub, args)
	} else if function == "invoke_transaction_insert_update" { //generic insert on ledger
		return invoke_transaction_insert_update(stub, args)
	} else if function == "add_theatre" { //add theatre details on ledger
		return add_theatre(stub, args)
	} else if function == "add_movies" { //add movie details on ledger
		return add_movies(stub, args)
	} else if function == "book_tickets" { //book movie tickets
		return add_movies(stub, args)
	}

	// error out
	fmt.Println("Received unknown invoke function name - " + function)
	return shim.Error("Received unknown invoke function name - '" + function + "'")
}

// ============================================================================================================================
// Query - legacy function
// ============================================================================================================================
func (t *MTA) Query(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Error("Unknown supported call - Query()")
}
