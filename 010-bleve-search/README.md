# Bleve Search

## Purpose

This project provides a REST API, GRPC for managing search indexes using Bleve, a modern text indexing library for Go. The API allows users to create and delete indexes, index documents, search for documents, and perform various other operations related to index management. The primary goal is to offer a robust and efficient search solution that can be easily integrated into other applications.

---

## Structure
The project is structured as follows:

```bash
├── cmd
│   ├── grpc
│   └── restapi
├── go.mod
├── go.sum
├── integration
│   ├── integration_tests
│   ├── intergration_test.go
│   └── poc_test.go
├── internal
│   ├── blevefunc
│   ├── controllers
│   ├── grpc
│   ├── repositories
│   ├── routes
│   └── utils
├── pkg
│   └── log
├── proto
│   ├── index_grpc.pb.go
│   ├── index.pb.go
│   └── index.proto
└── README.md
```
- cmd/server/: Contains the main server code and the .env file for environment configuration.
  - cmd/server/docs/: Contains the generated Swagger documentation files.
- integration/: Contains integration tests and related resources.
- internal/controllers/: Contains the controllers that handle the API requests.
- internal/repositories/: Contains the code for data access and manipulation.
- internal/routes/: Contains the routing logic for the API endpoints.
- internal/utils/: Contains utility functions and helpers used throughout the project.


---
## Environment
Create a file name ```.env``` in *cmd/server*:
```bash
cmd/server/.env

INDEX_DIRECTORY = "./INDEXES"
PORT = "8080"
DOMAIN = ""
PREFIX_PATH = ""
```

Init Swagger file:
```bash
swag init -d ./cmd/server/,./internal/controllers/ -o ./cmd/server/docs/ -md ./cmd/server/docs/
```

Generate protoc file:
```bash
protoc --go_out=proto --go-grpc_out=proto proto/index.proto
```
