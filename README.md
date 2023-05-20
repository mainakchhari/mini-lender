# mini-lender

## Description

This is a lending app (backend) that lets customers apply for loans, with a (weekly) repayment
period of their choice. After approval, customers can repay their loans weekly with any `amount >
due amount`.

## Local Setup

Within the repo rootDir run `make build` to build a binary in /bin folder. Run the binary
using `make run`. In another open terminal, run `curl -v http://localhost:8080` to get `HTTP/1.1 200 OK` as response.

Run all unit tests using `make test`, which will also generate a coverage profile file `coverage.out`. Now you can check out your coverage results in a browser window by running `go tool cover -html=coverage.out`

## Package Structure

```shell
.
└── internal
    └── app
        ├── adapter
        │   ├── controllers                   # Controllers
        │   │   ├── loan.go
        │   │   └── user.go
        │   ├── middleware
        │   │   └── auth.go                   # Middleware
        │   └── repository                    # Repository Implementation
        │       ├── loan.go
        │       ├── payment.go
        │       └── user.go
        ├── application
        │   └── errors                        # Errors
        │   └── usecase                       # Usecases (Business Logic/Flow)
        │       └── approverLifecycle.go
        │       └── customerLifecycle.go
        │       └── userLifecycle.go
        └── domain
            ├── constants.go                  # Domain constants
            ├── loan.go                       # Entity
            ├── payment.go                    # Entity
            ├── user.go                       # Entity
            │── repository                    # Repository Interfaces
            │   ├── loan.go
            │   ├── payment.go
            │   └── user.go
            │── mock_repository               # Repository Mocks
            │   ├── loan.go
            │   ├── payment.go
            │   └── user.go
            │── factory                       # Domain Factory
            │   ├── loan.go
            │   ├── payment.go
            │   └── user.go
            └── mock_factory                  # Factory Mocks
                ├── loan.go
                ├── payment.go
                └── user.go

```
