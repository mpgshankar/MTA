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
	"strings"
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
	var key, theatreRegNo, value string
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

	// Create Movie Object
	var mov Movies
	mov.ObjectType = "Movies"
	mov.MovieId = key
	mov.MovieName = movieName
	mov.TheatreRegNo = theatreRegNo

	//check if theatre exists or not
	theatreAsBytes, _ := stub.GetState(theatreRegNo)
	if theatreAsBytes == nil {
		fmt.Println("This theatre does not exists - " + theatreRegNo)
		return shim.Error("This theatre does not exists - " + theatreRegNo)
	}
	theatre := Theatre{}
	json.Unmarshal(theatreAsBytes, &theatre) //un stringify it aka JSON.parse()

	// check movies when it will be releasing
	mov.Status = "Running"
	if len(theatre.MoviesRunning) < theatre.NumberOfScreens {
		theatre.MoviesRunning = append(theatre.MoviesRunning, mov)
	} else {
		return shim.Error("Only " + string(theatre.NumberOfScreens) + " movies can run for this theatre " + theatreRegNo)
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

// ============================================================================================================================
// add_shows() - add shows for a movie into ledger
//
// Shows Off PutState() - writting a key/value into the ledger
//
// Inputs - JSON Object
//    0
//   json_object
//  {"theatreRegNo":"value1","theatreLocation":"value2"}
// ============================================================================================================================
func add_shows(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var theatreRegNo, value string
	// var err error
	fmt.Println("starting add_shows")

	certname, err := get_cert(stub)
	if err != nil {
		fmt.Printf("INVOKE: Error retrieving cert: %s", err)
		return shim.Error("Error retrieving cert")
	}

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting Minimum 1. arguments of the variable and value to set")
	}
	theatreRegNo = string(certname)
	//check if theatre exists or not
	theatreAsBytes, _ := stub.GetState(theatreRegNo)
	if theatreAsBytes == nil {
		fmt.Println("Only theatres can add shows for a movie - " + theatreRegNo)
		return shim.Error("Only theatres can add shows for a movie - " + theatreRegNo)
	}
	ttr := Theatre{}
	json.Unmarshal(theatreAsBytes, &ttr)

	value = args[0]
	var show Shows
	json.Unmarshal([]byte(value), &show)
	show.TotalSeat = 100
	show.AvailableSeat = 100
	show.BookedSeat = 0
	show.ShowStatus = "Running"
	if strings.HasSuffix(show.ShowTiming, "am") {
		show.PricePerTicket = 100
	} else {
		show.PricePerTicket = 180
	}

	//check if theatre exists or not
	movieAsBytes, _ := stub.GetState(show.MovieId)
	if movieAsBytes == nil {
		fmt.Println("Only theatres can add shows for a movie - " + theatreRegNo)
		return shim.Error("Only theatres can add shows for a movie - " + theatreRegNo)
	}

	mov := Movies{}
	json.Unmarshal(movieAsBytes, &mov)
	if theatreRegNo != mov.TheatreRegNo {
		fmt.Println("You cannot add a show for a movie which is not running in - " + theatreRegNo)
		return shim.Error("You cannot add a show for a movie which is not running in - " + theatreRegNo)
	}

	screenNumber := screenAvailable(ttr.NumberOfScreens, show.ShowTiming, stub)
	if screenNumber == 0 {
		fmt.Println("All the screens are full for this show timing. Please select different time for show")
		return shim.Error("All the screens are full for this show timing. Please select different time for show")
	} else {
		show.ScreenNumber = screenNumber
	}

	showAsBytes, _ := json.Marshal(show)

	errTr := stub.PutState(show.ShowId, showAsBytes) // update the theatre details into the ledger
	if errTr != nil {
		return shim.Error("Failed to add shows : " + errTr.Error())
	}

	fmt.Println("- end add_shows")
	return shim.Success(nil)
}

// ============================================================================================================================
// book_tickets() - Bur Movie Tickets and record into ledger
//
// Shows Off PutState() - writting a key/value into the ledger
//
// Inputs - JSON Object
//    0
//   json_object
//  {"theatreRegNo":"value1","theatreLocation":"value2"}
// ============================================================================================================================
func book_tickets(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key, theatreRegNo, value string
	// var err error
	fmt.Println("starting book_tickets")

	value = args[0]
	var ticket Tickets
	json.Unmarshal([]byte(value), &ticket)

	//check if theatre exists or not
	theatreAsBytes, _ := stub.GetState(theatreRegNo)
	if theatreAsBytes == nil {
		fmt.Println("This theatre does not exists - " + theatreRegNo)
		return shim.Error("This theatre does not exists - " + theatreRegNo)
	}
	theatre := Theatre{}
	json.Unmarshal(theatreAsBytes, &theatre) //un stringify it aka JSON.parse()
	movieId := ticket.MovieId
	noOfTickets := ticket.NumberOfTickets
	movieTiming := ticket.ShowTiming
	// moviesRunning := theatre.MoviesRunning
	fmt.Println(noOfTickets)
	fmt.Println(movieTiming)

	//check if movie exists or not
	movieAsBytes, _ := stub.GetState(movieId)
	if theatreAsBytes == nil {
		fmt.Println("This movie does not exists - " + movieId)
		return shim.Error("This movie does not exists - " + movieId)
	}
	movies := Movies{}
	json.Unmarshal(movieAsBytes, &movies) //un stringify it aka JSON.parse()
	fmt.Println(key)

	if movies.Status == "Running" || movies.Status == "Coming Soon" {
		// for i, _ := range movies.ShowTimings {
		// 	if movies.ShowTimings[i].ShowTiming == movieTiming {
		// 		if movies.ShowTimings[i].AvailableSeat == 0 {
		// 			return shim.Error("Booking limit reached. Show is Housefull")
		// 		} else if movies.ShowTimings[i].AvailableSeat >= noOfTickets {
		// 			ticket.TotalPrice = noOfTickets * movies.ShowTimings[i].PricePerTicket
		// 			movies.ShowTimings[i].AvailableSeat -= noOfTickets
		// 			movies.ShowTimings[i].BookedSeat += noOfTickets
		// 		} else {
		// 			movies.ShowTimings[i].ShowStatus = "HouseFull"
		// 		}
		// 	}
		// }
	}

	// for i, _ := range theatre.MoviesRunning {
	// 	if movieId == theatre.MoviesRunning {
	// 		if movie.Status == "Running" || movie.Status == "Coming Soon" {
	// 			showTimings := movie.ShowTimings
	// 			for _, show := range showTimings {
	// 				if show.ShowTiming == movieTiming {
	// 					ticket.TotalPrice = noOfTickets * show.PricePerTicket
	// 				}
	// 			}

	// 		} else {
	// 			return shim.Error("Booking currently not allowed for " + movie.MovieName)
	// 		}
	// 	}

	// }

	trAsBytes, _ := json.Marshal(theatre)

	errTr := stub.PutState(theatreRegNo, trAsBytes) // update the theatre details into the ledger
	if errTr != nil {
		return shim.Error("Failed to add movies : " + errTr.Error())
	}

	fmt.Println("- end book_tickets")
	return shim.Success(nil)
}

//Check Whether Current Date greater than or equal to Relase Date
func screenAvailable(noOfScreen int, showTiming string, stub shim.ChaincodeStubInterface) int {
	queryFrm := `{"selector":{"docType":"Shows", "showTiming":"` + showTiming + `"}}`
	queryString := fmt.Sprintf(queryFrm)
	var screenNumber int
	var arrayOfScreensUsed []int
	var totalScreens []int

	for i := 1; i <= noOfScreen; i++ {
		totalScreens = append(totalScreens, i)
	}

	fmt.Println("========= screenAvailable start =========")
	queryResults, _ := getQueryResultForQueryString(stub, queryString)
	fmt.Println(queryResults)
	var arrayOfShows []Shows
	json.Unmarshal(queryResults, &arrayOfShows)
	fmt.Println(len(arrayOfShows))
	if len(arrayOfShows) > 0 {
		for _, eachShow := range arrayOfShows {
			show := Shows{}

			fmt.Println("eachShow =======> ")
			arrayOfScreensUsed = append(arrayOfScreensUsed, show.ScreenNumber)
			fmt.Println(eachShow)
		}
		if len(arrayOfScreensUsed) >= noOfScreen {
			return 0
		} else {
			for _, v := range totalScreens {
				for _, u := range arrayOfScreensUsed {
					if v != u {
						screenNumber = v
						break
					}
				}
			}
		}
	} else {
		screenNumber = 1
	}

	fmt.Println("========= screenAvailable end =========")

	return screenNumber
}

//Check Whether Current Date greater than or equal to Relase Date
func greaterThanEqualCurrentDate(start, check time.Time) bool {
	return start.After(check) || start.Equal(check)
}

//Check Whether Current Date equal to Relase Date
func equalCurrentDate(start, check time.Time) bool {
	return start.Equal(check)
}
