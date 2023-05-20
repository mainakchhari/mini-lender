# mini-lender

## Description

This is a lending app (backend) that lets customers apply for loans, with a (weekly) repayment
period of their choice. After approval, customers can repay their loans weekly with any `amount >
due amount`.

## Local Setup

Within the repo rootDir run `make build` to build a binary in /bin folder. Run the binary
using `make run`. In another open terminal, run `curl -v http://localhost:8080` to get `HTTP/1.1 200 OK` as response
