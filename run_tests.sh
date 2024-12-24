#!/bin/bash

# Run tests with coverage
go test ./... -coverprofile=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

echo "Test coverage report generated at coverage.html"
