# Cosmos CLI

Code scaffolding tool `cosm` for Cosmos SDK applications.

## Installation

```
go get github.com/allinbits/cosm
```

Note: Make sure to run this command outside of an existing Go project, otherwise it will be added as a dependency to `go.mod`.

## Creating an application

```
cosm app [modulePath]
```

This command creates an empty template for a Cosmos SDK application. By default it also includes a module with the same name as the package.

To create a new application called `blog`, run:

```
cosm app github.com/example/blog
```

## Running an application

```
cosm serve
```

To start the server, go into you application's directory and run `cosm serve`. This commands installs dependencies, builds and initializes the app and runs both Tendermint RPC server (by defalut on `localhost:26657`) as well as LCD (by defalut on `localhost:1317`).

## Creating types

```
cosm type [typeName] [field1] [field2:bool] ...
```

This command generates messages, handlers, keepers, CLI and REST clients and type definition for `typeName` type. A type can have any number of `field` arguments. By default fields are strings, but `bool` and `int` are supported.

For example,

```
cosm type post title body
```

This command generates a type `Post` with two fields: `title` and `body`.

To add a post run `blogcli tx blog create-post "My title" "This is a blog" --from=me`.
