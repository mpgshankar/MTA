# Movie Ticket Application (MTA)
An application that can be used to book movie tickets and record each transaction over the blockchain ledger.
```sh
# Step 1 :
## Onboard Theatre
This is the first step inorder to run MTA. 
To onboard theatre we need to invoke `add_theatre` function which takes 
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
To add shows we need to invoke `add_shows` function which takes only 1 argument of JSON Object.
Sample :- {"showId":"value1","showTiming":"value2", "movieId":"value3","docType":"value4"}
Here showId can be any unique Id to distinguish between Shows for Movies
```