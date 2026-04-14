## Spread management service
### Overview
This is a simulated spread management service for a trading platform. It allows to set spread values for currency pairs,
read those values and records the time they were updated. I have decided to use Gin as it is a clean library to route the endpoints to the functions and business logic.
The servie maintains a `SpreadHandler` which tracks and manages spreads and update times. It is initialised from a list of currency pairs.

### API endpoints
The service has 3 endpoints:
- GET `/symbols`: returns a list of currency pairs for which spreads are managed by the service. Example request: `curl http://localhost:8080/symbols`
- GET `/spreads/<symbol>`: returns the spread if it exists and is set, alongside the last updated timestamp with 200 OK. Otherwise returns 404 Not found. Example request: `curl http://localhost:8080/spreads/BTCUSD`
- PATCH `spreads/<symbol>`: update the spread for a currency pair. Must include `"spread" : <amount>"` in request body. Returns 200 OK if succesfully updated, 404 Not found if symbol not managed by the service, 400 bad request if amount is not > 0.  
    Example request: `curl -X PATCH http://localhost:8080/spreads/BTCUSD 
  -H "Content-Type: application/json" 
  -d '{"spread": 0.1}'`

### Build instructions
It is best to build and run the service with docker:
- Build: `docker build -t spread-service .`
- Run (on port 8080): `docker run -p 8080:8080 spread-service`    

It can also be run locally as `go run ./cmd/main.go`.
The service is launched on `http://localhost:8080`   
To run tests: `go test ./...`

