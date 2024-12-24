# Cell Router

## Build the Project

To build the project, run the following command:

```sh
docker-compose build
```

## Run Locally

To run the project locally, use the following command:

```sh
docker-compose up
```

This will start the `cell-router` service on port 8080 and the mock cells on ports 8081 and 8082.

## Test Locally

To test the routing, you can use `curl` or any HTTP client to send requests to the `cell-router` service:

```sh
# Test routing to cell1
curl "http://localhost:8080?client_id=50"

# Test routing to cell2
curl "http://localhost:8080?client_id=150"
```

You should receive responses from the respective mock cells based on the `client_id` provided.