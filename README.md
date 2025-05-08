# Digi-Wallet

Hi! This is Geizka Rozilia Rui**costa**'s submission for KlikCair's application as Backend Engineer.

# How to run
*Make sure you have go version 21 installed*
 1. Git clone the repo
 2. go mod tidy
 3. import the postman collection to Postman
 4. run `go run main.go`
 5. hit the register API from postman
    `{"username":  "Spencer", "password":  "1"}`

 6. hit the login API from postman with body raw:
	  `{"username":  "Spencer", "password":  "1"}`
 7. Use the auth token received in auth token variable
 8. hit the get balance api with the bearer token to get balance (initially 100.000)
 9. hit the take balance api with the bearer token with body containing form-data
 key: `amount`
 value: `1000` (example)