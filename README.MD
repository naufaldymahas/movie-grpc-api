# Movie Search
Movie Search using API from OMDb API, Golang, MySQL, GRPC, and REST API

*this project using Golang 1.17

## Prerequisites
- [GO 1.17](https://golang.org/)
- Protocol Buffer Compiler, click [here](https://grpc.io/docs/protoc-installation/) to see installation guide
- [MySQL](https://www.mysql.com/)

## Usage
- Create table **movie_logs** using file ```sqlFile/create_table.sql```

Set environtment variable first before to start this project
- Create ```.env``` inside the project
- Then use the example below
```.env
MYSQL_DB_USER=root (your MySQL user)
MYSQL_DB_PASSWORD=password (your MySQL password)
MYSQL_DB_NAME=bibit-test (your MySQL database name)

OMDB_API=http://www.omdbapi.com/
OMDB_API_KEY=xxxxx (your key from OMDb API)
GRPC_PORT=8080 (port for GRPC server, default value 8080)
REST_PORT=8081 (port for REST API server, default value 8081)
```


For Windows user
- To recompile Protobuf use this ```.\gen-protoc.bat```
- To start the server use this ```.\run.bat```

For Mac/Linux user
- To recompile Protobuf use this
```
protoc --proto_path=proto --go_out=plugins=grpc:pb movie.proto
protoc --proto_path=proto --grpc-gateway_out=logtostderr=true:pb movie.proto
protoc --proto_path=proto --openapiv2_out=logtostderr=true:docs movie.proto
```
- To start the server use this
```go
go run main.go
```
or
```go
go build
./movie-grpc-api
```

## Author
Naufaldy Mahas (naufaldymahas@gmail.com)