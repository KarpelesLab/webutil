# Go Webutil Guidelines

## Build & Test Commands
- Build: `make all` or `go build -v ./...`
- Run all tests: `make test` or `go test -v ./...`
- Run single test: `go test -v -run TestName` (e.g., `go test -v -run TestParseIPPort`)
- Format/Lint: `goimports -w -l .`
- Install dependencies: `make deps` or `go get -v -t ./...`

## Code Style Guidelines
- **Package Structure**: Main package is `webutil`, tests use `webutil_test`
- **Imports**: Group and organize imports using `goimports`
- **Naming**: CamelCase for exported functions/types, camelCase for unexported
- **Error Handling**: Use custom error types and `errors.As()` pattern; wrap errors with context using `fmt.Errorf` with `%w`
- **Comments**: Document functions in godoc format; mark deprecated functions
- **Types**: Use interfaces when appropriate; define custom struct types for specialized purposes
- **Tests**: Implement table-driven tests with clear expected outputs
- **Function Design**: Follow single responsibility principle; use clear parameter and return types