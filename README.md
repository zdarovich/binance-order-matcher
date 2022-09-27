# binance-order-matcher

## Project structure
### `cmd/app/main.go`
Configuration and logger initialization. Then the main function "continues" in
`internal/app/app.go`.
### `config`
Configuration. First, `config.yml` is read, then environment variables overwrite the yaml config if they match.
### `internal/app`

This is where all the main objects are created.
### `internal/model`
Entities of business logic (models) can be used in any layer.
There can also be methods, for example, for validation.

### `internal/service`
Business logic.
- Methods are grouped by area of application (on a common basis)
- Each group has its own structure
- One file - one structure

### `pkg/sqlite`
- sqlite connector

## Dependency Injection
In order to remove the dependence of business logic on external packages, dependency injection is used.

## Difficulties
- Gracefully connect and shutdown websocket connection and proccess that converts API responses in `internal/service/order_matcher.go`

## Likes
- Free usage of  Binance websocket API

## Time spent
- 6-7 hours

## What shoud be done better?
- Refactor matcher method in a more readable format.
- Test channel communication