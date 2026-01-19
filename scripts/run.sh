#!/bin/bash

# Build the application
go build -o supermarket-app ./cmd/web/main.go

./supermarket-app