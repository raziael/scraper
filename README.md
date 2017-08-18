# scraper
This service performs a person lookup by phone number. If multiple records are returned then it looks for the closes match using the name (by approximation).

Once a match is found, the record is saved in local storage. 

## Searching
Searching is done be entering the full name along with the phone number. The format of the phone number is `xxx-xxx-xxxx`

## Start up the service
To start the service, execute `./start.sh` or `go run *.go`

