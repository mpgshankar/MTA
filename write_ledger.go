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
	"math/rand"
	"strconv"
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
//  {"theatreRegNo":"value1","theatreLocation":"value2","theatreName":"value3","numberOfScreens":"value4","docType":"value5"}
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

	//check if theatre already exists
	tr, _ := stub.GetState(key)
	if tr != nil {
		fmt.Println("This theatre already exists - " + key)
		return shim.Error("This theatre already exists - " + key)
	}

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
//  {"movieId":"value1","movieName":"value2","docType":"value3"}
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
		return shim.Error("Only " + strconv.Itoa(theatre.NumberOfScreens) + " movies can run for this theatre " + theatreRegNo)
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
//  {"showId":"value1","showTiming":"value2", "movieId":"value3","docType":"value4"}
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
	// show.ObjectType = "Shows"
	show.TotalSeat = 100
	show.AvailableSeat = 100
	show.BookedSeat = 0
	show.ShowStatus = "Running"
	show.TheatreRegNo = theatreRegNo
	show.ShowDate = show.ShowTiming[:10]
	if strings.HasSuffix(show.ShowTiming, "am") {
		show.PricePerTicket = 100
	} else {
		show.PricePerTicket = 180
	}

	//check if show already exists
	sw, _ := stub.GetState(show.ShowId)
	if sw != nil {
		fmt.Println("This show already exists - " + show.ShowId)
		return shim.Error("This show already exists - " + show.ShowId)
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

	screenNumber := screenAvailable(ttr.NumberOfScreens, show.ShowTiming, show.ShowDate, show.MovieId, stub)
	if screenNumber == 0 {
		fmt.Println("All the screens are full for this show timing or this movie is already running for the show timing on another screen. Please select different time for show")
		return shim.Error("All the screens are full for this show timing or this movie is already running for the show timing on another screen. Please select different time for show")
	} else if screenNumber == 20 {
		fmt.Println("Only 4 shows are allowed for a day for a particular movie")
		return shim.Error("Only 4 shows are allowed for a day for a particular movie")
	}
	show.ScreenNumber = screenNumber
	showAsBytes, _ := json.Marshal(show)

	errShw := stub.PutState(show.ShowId, showAsBytes) // update the theatre details into the ledger
	if errShw != nil {
		return shim.Error("Failed to add shows : " + errShw.Error())
	}

	var acc Accessories
	acc.ObjectType = "Accessories"
	acc.Asset = "Soda"
	acc.TotalQty = 200
	acc.ForDate = show.ShowDate
	acc.AvailableQty = 200

	//check if Accessories already exists
	access, _ := stub.GetState(acc.ForDate)
	if access == nil {
		accAsBytes, _ := json.Marshal(acc)
		errAcc := stub.PutState(acc.ForDate, accAsBytes) // update the theatre details into the ledger
		if errAcc != nil {
			return shim.Error("Failed to add shows : " + errAcc.Error())
		}
	}

	fmt.Println("- end add_shows")
	return shim.Success(nil)
}

// ============================================================================================================================
// book_tickets() - Buy Movie Tickets and record into ledger
//
// Shows Off PutState() - writting a key/value into the ledger
//
// Inputs - JSON Object
//    0
//   json_object
//  {"showId":"value1","numberOfTickets":"value2"}
// ============================================================================================================================
func book_tickets(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var value string
	var ticketAsBytes []byte
	// var err error
	fmt.Println("starting book_tickets")

	value = args[0]
	var ticket Tickets
	json.Unmarshal([]byte(value), &ticket)
	shAsBytes, _ := stub.GetState(ticket.ShowId)
	show := Shows{}
	json.Unmarshal(shAsBytes, &show)
	if show.AvailableSeat == 0 {
		return shim.Error("Failed to book tickets for show as no seats are available.")
	} else if ticket.NumberOfTickets <= show.AvailableSeat {
		movieAsBytes, _ := stub.GetState(show.MovieId)
		mov := Movies{}
		json.Unmarshal(movieAsBytes, &mov)
		rand.Seed(time.Now().UnixNano())

		ticket.ObjectType = "Tickets"
		ticket.TicketId = "T" + show.TheatreRegNo + show.ShowId + strconv.Itoa(rand.Intn(1000000))
		ticket.MovieName = mov.MovieName
		ticket.ShowTiming = show.ShowTiming
		ticket.TotalPrice = show.PricePerTicket * ticket.NumberOfTickets
		ticket.ScreenNumber = show.ScreenNumber
		show.BookedSeat += ticket.NumberOfTickets
		show.AvailableSeat -= ticket.NumberOfTickets

		for i := 1; i <= ticket.NumberOfTickets; i++ {
			var amn Amenities
			amn.SeatNumber = "S" + show.TheatreRegNo + show.ShowId + strconv.Itoa(rand.Intn(100000))
			amn.PopCorn = 1
			amn.Water = 1
			ticket.Amenities = append(ticket.Amenities, amn)
		}

		ticketAsBytes, _ = json.Marshal(ticket)
		errTkt := stub.PutState(ticket.TicketId, ticketAsBytes) // update the theatre details into the ledger
		if errTkt != nil {
			return shim.Error("Failed to book tickets : " + errTkt.Error())
		}

		showAsBytes, _ := json.Marshal(show)
		errShow := stub.PutState(show.ShowId, showAsBytes) // update the theatre details into the ledger
		if errShow != nil {
			return shim.Error("Failed to book tickets : " + errShow.Error())
		}

	} else {
		return shim.Error("Failed to book tickets for shows as only " + strconv.Itoa(show.AvailableSeat) + " seats are available.")
	}

	fmt.Println("- end book_tickets")
	return shim.Success(ticketAsBytes)
}

// ============================================================================================================================
// exchange_water() - Exchange Water with Soda and record into ledger
//
// Shows Off PutState() - writting a key/value into the ledger
//
// Inputs - JSON Object
//    0
//   json_object
//  {"ticketId":"value1"}
// ============================================================================================================================
func exchange_water(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var ticketId, value string
	// var err error
	fmt.Println("starting exchange_water")
	rand.Seed(time.Now().UnixNano())
	randNo := rand.Intn(200)
	fmt.Println("randNo ====> ")
	fmt.Println(randNo)

	if randNo%2 == 0 {
		value = args[0]
		var jsonValue map[string]interface{}
		json.Unmarshal([]byte(value), &jsonValue)
		ticketId, _ = jsonValue["ticketId"].(string)

		tktAsBytes, _ := stub.GetState(ticketId)
		ticket := Tickets{}
		json.Unmarshal(tktAsBytes, &ticket)
		forDate := ticket.ShowTiming[:10]

		accAsBytes, _ := stub.GetState(forDate)
		acc := Accessories{}
		json.Unmarshal(accAsBytes, &acc)
		fmt.Println(acc)
		if ticket.NumberOfTickets <= acc.AvailableQty {
			for i, amn := range ticket.Amenities {
				if amn.Soda == 1 {
					return shim.Error("Soda has already been exchanged for this ticket.")
				}
				amn.Water = 0
				amn.Soda = 1
				ticket.Amenities[i] = amn
				// fmt.Println(ticket.Amenities[i])
			}
			acc.AvailableQty -= ticket.NumberOfTickets
			ticketAsBytes, _ := json.Marshal(ticket)
			errTkt := stub.PutState(ticketId, ticketAsBytes) // update the theatre details into the ledger
			if errTkt != nil {
				return shim.Error("Failed to exchange_water : " + errTkt.Error())
			}
			fmt.Println(acc)

			accessAsBytes, _ := json.Marshal(acc)
			errAcc := stub.PutState(forDate, accessAsBytes) // update the theatre details into the ledger
			if errAcc != nil {
				return shim.Error("Failed to exchange_water : " + errAcc.Error())
			}

		} else {
			fmt.Println("Soda is out of stock.")
			return shim.Error("Soda is out of stock.")
		}
	} else {
		fmt.Println("Exchanging Water with Soda is currently not possible.")
		return shim.Error("Exchanging Water with Soda is currently not possible.")
	}

	fmt.Println("- end exchange_water")
	return shim.Success(nil)
}

// Assigns screen number for a particular show
func screenAvailable(noOfScreen int, showTiming string, showDate string, movieId string, stub shim.ChaincodeStubInterface) int {
	// Compares whether a movie is not running more than 4 times a day.
	queryFrmDate := `{"selector":{"docType":"Shows", "showDate":"` + showDate + `"}}`
	queryStringDate := fmt.Sprintf(queryFrmDate)

	queryFrm := `{"selector":{"docType":"Shows", "showTiming":"` + showTiming + `"}}`
	queryString := fmt.Sprintf(queryFrm)
	var screenNumber int
	var arrayOfScreensUsed []int
	var totalScreens []int
	var unique []int
	showsPerDay := 1
	for i := 1; i <= noOfScreen; i++ {
		totalScreens = append(totalScreens, i)
	}

	queryResultsDate, _ := getQueryResultForQueryString(stub, queryStringDate)
	var arrayOfShowsDate []Shows
	json.Unmarshal(queryResultsDate, &arrayOfShowsDate)
	if len(arrayOfShowsDate) > 0 {
		for _, eachShowDate := range arrayOfShowsDate {
			if showDate == eachShowDate.ShowDate && movieId == eachShowDate.MovieId {
				showsPerDay += 1
			}
		}
		if showsPerDay > 4 {
			return 20 // Maximum Shows for a particular movie reached
		}
	}

	// Assigns screens for a particular show for a movie.
	queryResults, _ := getQueryResultForQueryString(stub, queryString)
	var arrayOfShows []Shows
	json.Unmarshal(queryResults, &arrayOfShows)
	if len(arrayOfShows) > 0 {
		for _, eachShow := range arrayOfShows {
			arrayOfScreensUsed = append(arrayOfScreensUsed, eachShow.ScreenNumber)

			if eachShow.MovieId == movieId && eachShow.ShowTiming == showTiming {
				return 0
			}
		}
		if len(arrayOfScreensUsed) >= noOfScreen {
			return 0
		} else {
			unique = Difference(totalScreens, arrayOfScreensUsed)
			for _, val := range unique {
				return val
			}
		}
	} else {
		screenNumber = 1
	}
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

// Set Difference: A - B
func Difference(a, b []int) (diff []int) {
	m := make(map[int]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return diff
}
