# CHPTER Service

### This repository contains a user service, an order service and integration tests that test service concurrency.

The order and user services share rpc files via a symbolic link

The whole application can be started by running `docker compose up`

## Order Service

- The order service can be found here `./services/order`
- You can run the order service's unit tests by running `make order_test` from the root directory

## User Service

- The user service can be found here `./services/user`
- You can run the user service's unit tests by running `make user_test` from the root directory

## Integration Test

- The application ships an integration test that tests the connection between the order service and the user service.

- Integration tests can be run by running `make e2e_test` from the root directory

## Protobuf

- The application's protobuf files can be compiled by running `make proto_gen` from the root directory
