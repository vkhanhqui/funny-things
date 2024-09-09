# Bleve Search

## Purpose

This project provides a REST API and gRPC for managing search indexes using Bleve, a modern text indexing library for Go. The API allows users to create and delete indexes, index documents, search for documents, and perform various other operations related to index management. The primary goal is to offer a robust and efficient search solution that can be easily integrated into other applications.

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
│   ├── integration_test.go
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
- `cmd/server/`: Contains the main server code and the `.env` file for environment configuration.
  - `cmd/server/docs/`: Contains the generated Swagger documentation files.
- `integration/`: Contains integration tests and related resources.
- `internal/controllers/`: Contains the controllers that handle the API requests.
- `internal/repositories/`: Contains the code for data access and manipulation.
- `internal/routes/`: Contains the routing logic for the API endpoints.
- `internal/utils/`: Contains utility functions and helpers used throughout the project.

---

## Environment
Create a file named `.env` in `cmd/server`:

```bash
cmd/server/.env

INDEX_DIRECTORY="./INDEXES"
PORT="8080"
DOMAIN=""
PREFIX_PATH=""
```

Initialize Swagger files:

```bash
swag init -d ./cmd/server/,./internal/controllers/ -o ./cmd/server/docs/ -md ./cmd/server/docs/
```

Generate protobuf files:

```bash
protoc --go_out=proto --go-grpc_out=proto proto/index.proto
```

---

## CVE Description

**Important Notice:**

Bleve is a text indexing library for Go, including HTTP utilities under the `bleve/http` package. These utilities are for demonstration purposes only and lack robust security measures. Using handlers like `CreateIndexHandler` and `DeleteIndexHandler` can allow attackers to manipulate the filesystem where Bleve indexes reside, creating or deleting directories where the server has write permissions.

**Warning:** Do not use `bleve/http` in production. It lacks proper Role-Based Access Controls (RBAC), authentication, and authorization. Secure your application by implementing custom handlers with appropriate security measures.

For more details, refer to the [official CVE-2022-31022 description](https://nvd.nist.gov/vuln/detail/CVE-2022-31022/change-record?changeRecordedOn=06/02/2022T10:15:56.647-0400).