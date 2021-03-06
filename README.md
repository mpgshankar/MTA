# Movie Ticket Application (MTA)
An application that can be used to book movie tickets and record each transaction over the blockchain ledger.
```sh
# Step 1 :
## Add Theatre
Here multiple theatres can be added where unique ID is theatreRegNo.
To add theatre we need to invoke `add_theatre` function which takes 
only 1 argument of JSON Object.
Sample :- {"theatreRegNo":"value1","theatreLocation":"value2","theatreName":"value3","numberOfScreens":"value4","docType":"value5"}

# Step 2 :
## Add Movies
Once the theatre is onboarded movie can be added into that Theatre. 
Adding movies will be done using credentials of Theatre. 
To add movies we need to invoke `add_movies` function which takes only 1 argument of JSON Object.
Sample :- {"movieId":"value1","movieName":"value2","docType":"value3"}
Here movieId can be any unique Id to distinguish between Movies

# Step 3 :
## Add Shows
Once the movie has been added for a Theatre. Theatre user can add shows using their credentials.
While adding shows the application will itself identify available screens on which the current
show will be running. Also sanity checks like each day max 4 shows can run for a Movie is also done.  
To add shows we need to invoke `add_shows` function which takes only 1 argument of JSON Object.
Sample :- {"showId":"value1","showTiming":"value2", "movieId":"value3","docType":"value4"}
Here showId can be any unique Id to distinguish between Shows for Movies

# Step 4 :
## Book Tickets
Once the shows are visible to buyers, now they can book tickets for any show they want to.
To book tickets we need to invoke `book_tickets` function which takes only 1 argument of JSON Object.
Sample :- {"showId":"value1","numberOfTickets":"value2"}
In response the buyer gets the ticket details along with amenities like Water Bottle and Pop Corn.
Later buyer can exchange water bottle with soda if required.

# Step 5 :
## Exchange Water
Post booking of ticket by buyer can exchange water with soda, but only 200 customers will only be 
able to avail this offer per day. 
To exchange water with soda we need to invoke `book_tickets` function which takes only 1 argument of JSON Object.
Sample :- {"ticketId":"value1"}

# Note: This application is built on CouchDB as primary database for hyperledger fabric as we can use 
# rich queries to fetch the details as required. Below mentioned functions are already available in 
# this application.

# Query 1 :
## Generic Query
This takes query input from user and gives all the desired output.
To use this we need to call `generic_query` function. 
Sample for passing query :- {"selector":{"docType":"Shows", "showTiming":"value"}}

# Query 2 :
## Generic Query with pagination
This takes query input from user along with page size and gives the desired output as a pagination result.
To use this we need to call `generic_query_pagination` function. 
Sample for passing query :- [{"selector":{"docType":"Shows", "showTiming":"value"}},10,""]


```
