# gin-ent-example  

Example application of Gin (Go web framework) and Ent (ORM).

## Setup A Go Environment
```shell
$ go mod tidy
```

## Install entc
```shell
$ go install entgo.io/ent/cmd/entc
```

## Create Your Schema
```shell
# Run `entc init` from the root directory of the project as follows:
$ entc init User
```

# Generate Code
```shell
# Run `go generate` from the root directory of the project as follows:
$ go generate ./ent
```

## Schema describe
```shell
$ entc describe ./ent/schema
```

## Install swag
```shell
$ go get -u github.com/swaggo/swag/cmd/swag
```

## Build and Run
```shell
$ go mod tidy
$ go run main.go
```  
and now you can visit this url for the test of API: <http://localhost:8088/swagger/index.html>  
