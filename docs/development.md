
## Usage

This project is using Go version 1.21 ideally use the GVM to be able to run it 

    gvm install go1.12.1
    gvm use go1.12.1

If you need to install the gvm simply run 

    bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)



## Development

### Getting Started

    # install golang
    brew install golang

    # install the golangci linter 
    # more details: https://golangci-lint.run/
    brew install golangci-lint
    
    # install pre-commit
    pip install pre-commit
    pre-commit install

    # Download all dependencies
    go mod download


### Testing

All test files are named *_test.go. Github workflow automatically run the tests when code is pushed and will return a report with results when finished.

You can also run the tests locally:

    go test ./...

To run the tests with coverage:

    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out

