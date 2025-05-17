# Overview
A Project Starter for https://github.com/irdaislakhuafa/go-sdk.git.

This project includes setup for:
- SQL Query Builder using [sqlc](https://sqlc.dev/)
- SMTP Client using [gomail](https://github.com/go-gomail/gomail)
- Rest API Server with [fiber](https://github.com/gofiber/fiber)
- SDK Library from [go-sdk](https://github.com/irdaislakhuafa/go-sdk)

## Project Structure
```bash
.
├── deploy -- configurations for deploy the apps.
├── docker-compose.yaml -- docker compose configurations.
├── docs -- documents assets, can be email template/sql code/etc.
│   └── sql -- sql queries and schemas
├── etc -- other configurations for app.
├── flake.nix -- nix flake configurations.
├── sqlc.yaml -- sqlc query generator configurations.
└── src -- source code of app.
    ├── business -- for business layers of app.
    │   ├── domain -- for low level layer to access resource like database/third party/etc.
    │   │   ├── domain.go -- instantiate domain layer.
    │   └── usecase -- layer for craft business logic aligns with business needs here.
    │       └── usecase.go -- instantiate usecase layer.
    ├── cmd -- command line app layer.
    │   └── main.go -- entrypoint code of app.
    ├── entity -- entity layer, craft your entity for db or etc here.
    │   ├── gen -- generated from sqlc.
    │   ├── rest.go -- craft your entity for rest api implementation here.
    ├── handler -- craft your handler collections for rest/graphql/grpc.
    │   └── rest -- rest api setup
    │       ├── helper.go -- code helper for rest api implementation
    │       ├── rest.go -- rest api entrypoint.
    │       ├── route.go -- write your rest api routes here.
    └── utils -- code utilities for app.
        ├── config -- contains code for app config configuration.
        ├── connection -- contains code for connection to db.
        ├── ctxkey -- context key collections.
        ├── pagination -- code collections to implement paginations.
        └── validation -- code collections for validation purpose.
```

## Usages
Below is usage docs for this project.

- [Installation](#installation)
- [Development](#development)
- [Run and Build](#run-and-build)

### Installation

A simple way to using this project, you just need to clone this project and custom setup based on your purpose.
```bash
$ git clone https://github.com/irdaislakhuafa/amartha-billing-engine.git
```

### Development
Big Thanks for [nix](https://github.com/NixOS/nix) and [nixpkgs](https://github.com/NixOS/nixpkgs) the awesome project for developer.

This project includes bundle with development tools for Go project:
- [helix](https://github.com/helix-editor/helix): A simple code editor but powerfull like vim/neovim.
- [gopls](https://github.com/golang/tools/tree/master/gopls): Golang language server protocol that provide auto completion and etc.
- [gotools](https://go.googlesource.com/tools): Golang tools collection for development that provide formatter, auto import and etc. 
- [go](https://github.com/golang/go): Golang compiler with version 1.24.
- [simple-completion-language-server](https://github.com/estin/simple-completion-language-server): A language server for text buffer.

#### Full Setup
Full setup contains all tools above for development. For usage you need to install [nix](https://github.com/NixOS/nix) package manager first. Then just type below:

```bash
$ nix develop .#ide
```

And [nix](https://github.com/NixOS/nix) will setup development environment for Go and you ready to develop your app without any effort to setup compiler or etc.

### Minimal Setup
May you need minimal setup if you:
- already have own `go` compiler installed.
- want to use your own `ide` or `code editor`.
- don't need auto complete for plain text from buffer.

And just need a minimal setup to develop your Go app you can type.

```bash
$ nix develop .#ide-minimal
```

This setup only contains:
- `gopls`
- `gotools`


After that you can open your project using `helix`.
```bash
$ hx .
```

Or using any IDE that you like.

### Run and Build
Build your app.
```bash
$ go build -o main ./src/cmd/main.go
```

Run the app.
```bash
$ ./main
```

Build with docker.
```bash
$ docker compose up app-dev -d
```
or
```bash
$ docker compose up app-prod -d
```
Depends about your [docker-compose](./docker-compose.yaml) configurations.

## Todo

- [x] Setup nix flake
- [x] Setup SQL Query Builder using sqlc.
- [x] Setup for Rest API
- [ ] Setup for GraphQL
- [ ] Setup for gRPC 
- [x] Setup for Minio Object Storage
