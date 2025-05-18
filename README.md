# Overview
This is Amartha Billing Engine that provide billing services for loaning. This project is an assignment for me as test from Amartha for Senior/Principal Engineer position.

This project is built with my SDK for Go Development [go-sdk](https://github.com/irdaislakhuafa/go-sdk).

## Usages
This app use docker as container so it's easy to deploy, you don't need to configure anything if you already have docker installed in your machine.

### Installation
Just run the following command in your terminal
```bash
$ docker compose up -d
```
Then wait until app is ready.

### API Docs
- Get Outstanding: `GET /api/v1/loan/transaction/calculate/:user_id`
- Is Delinquent: `GET /api/v1/users/:id`. Users have flag `deliquent_level` to identify they are delinquent or not.
- Make Payment: `POST /api/v1/loan/transaction/pay`

You can see full API docs in `docs/rest/Amartha Billing` directory. You can use [bruno](https://www.usebruno.com/) to open it. [Bruno](https://www.usebruno.com/) is a Postman alternative that fully free and open source.


### Demo for Loan and Pay loan bill
[Click Demo as Video Here](https://drive.google.com/file/d/1Dz3v0cPRlJah58pkyGK4OzHickxUfzCB/view?usp=sharing)
Below is a flow that used on demo video.
- Flow Admin
1. Create Loan plan
2. Setup Setting for `eod_date` and `limit_billing_for_delinquent`

- Flow User
1. Register
2. Login and use `token` on `Authorization` header
3. Create loan transaction
4. Check outstanding loan before billing
5. Then admin should update `eod_date` to next bill date
6. Check outstanding loan after billing
7. Pay loan transaction



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