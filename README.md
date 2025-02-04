# fetch-backend-challenge
Receipt processor as per the instructions in https://github.com/fetch-rewards/receipt-processor-challenge.
Contains a webservice with 2 API endpoints, `POST /receipts/process` and `GET /receipts/{id}/points`. Additional information can be found in the above link.

## Install and run:
1. Clone the repo to your favorite folder
2. run `go mod init favorite-folder/fetch-backend-challenge`
3. run `cd fetch-backend-challenge && go get .`
4. run `go run .`

The service will then be live! The above endpoints can be accessed via base url `localhost:8080`.

## Testing:
1. With the webservice running, run `go test`.


Thank you for your time!