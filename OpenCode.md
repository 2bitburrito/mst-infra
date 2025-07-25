# MST Infrastructure - OpenCode Guide

## Build/Test Commands
- `go build` - Build the main application
- `go run main.go` - Run the application locally
- `go test ./...` - Run all tests
- `go test ./jwt` - Run tests for specific package (jwt example)
- `go test -v ./jwt/create-JWT_test.go` - Run single test file with verbose output
- `go fmt ./...` - Format all Go code
- `go vet ./...` - Run Go static analysis

## Code Style Guidelines
- **Imports**: Standard library first, then third-party, then local packages (see main.go)
- **Naming**: Use camelCase for variables/functions, PascalCase for exported types
- **Types**: Define custom types like `PlanType` as string constants in utils/shared-types.go
- **Error Handling**: Use explicit error returns, log errors before returning JSON responses
- **JSON**: Use struct tags for JSON serialization, follow existing JsonReturn pattern
- **Database**: Use sqlc for type-safe SQL queries, store in db/sqlc/
- **Testing**: Use testify/require for assertions, create mock functions for dependencies
- **Middleware**: Follow existing pattern with apiMiddleware for CORS and auth
- **Config**: Load from environment variables via config.LoadConfig(), use godotenv for local dev
- **Structs**: Group related fields, use pointers for optional fields (see License struct)