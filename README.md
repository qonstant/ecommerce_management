# Ecommerce Management Service

## Goal

Creating a RESTful API for a Ecommerce Management Service with integration of [EPAY](https://epayment.kz/) API from Halyk Bank. The service is deployed on Render, containerized using Docker, and run with Docker Compose and Makefile. 

## Features
- **Database was created by using dbdiagram.**

![alt text](https://github.com/qonstant/ecommerce_management/blob/main/dbdiagram.png)

- **CRUD operations for managing tasks.**
  - CRUD was created using SQLC. For regeneration of CRUD:
    ```bash
    make sqlc
    ```

- **Swagger UI for API documentation.**
  - For regeneration of Swagger documentation:
    ```bash
    make swagger
    ```

- **Docker support for containerization.**
  - For running up:
    ```bash
    make up
    ```
  - For shutting down:
    ```bash
    make down
    ```
  - For restart:
    ```bash
    make restart
    ```
    
## Prerequisites

- Go 1.22 or later
- Docker (for containerization)
- SQlC (for CRUD generation)
```bash
brew install sqlc
```
or
```bash
go get github.com/kyleconroy/sqlc/cmd/sqlc
```

## Getting Started

### Clone the Repository

```bash
https://github.com/qonstant/ecommerce_management.git
cd ecommerce_management
```
## Build and Run Locally

### Build and Run the application:

Before running the app, create app.env file and copy and paste all the data from .env.dist file.

After this, you can run this command:
```bash
make up
```

After running it, the server will start and be accessible at 

```bash
http://localhost:8080/swagger/index.html
```

Link to the deployment: 
```bash
https://ecommerce-management-kwsu.onrender.com/swagger/index.html
```

### Health Check

Health can by checked by [LINK](https://ecommerce-management-kwsu.onrender.com/status)

### Generate Swagger Documentation

```bash
make swagger
```
### Run Tests
```bash
make tests
```

## Docker
### For database container running
```bash
make postgres
```
### For database creation
```bash
make createdb
```
### For database drop
```bash
make dropdb
```
### For docker compose up
```bash
make up
```
### For docker compose down
```bash
make down
```
### For restarting container
```bash
make restart
```

## API Endpoints
### All API endpoints can be accessed through swagger, but here is data for post requests

### Create a New User
- URL: http://localhost:8080/users
- URL: https://ecommerce-management-kwsu.onrender.com//users
- Method: POST
- Description: Create a new user
- Request Body:
```json
{
  "address": "Tole Bi 59",
  "email": "admin@kbtu.kz",
  "full_name": "Astana Nazarbayev",
  "role": "Project Manager"
}
```

### Create a New Product
- URL: http://localhost:8080/products
- URL: https://ecommerce-management-kwsu.onrender.com/products
- Method: POST
- Description: Create a new product
- Request Body:
```json
{
  "category": "Shampoo",
  "description": "Shampoo Zhumaisynba, against dandruff",
  "name": "Shampoo Zhumaisynba",
  "price": "300",
  "stock_quantity": 100
}
```

### Create a New Order
- URL: http://localhost:8080/orders
- URL: https://ecommerce-management-kwsu.onrender.com/orders
- Method: POST
- Description: Create a new order. It may contain several types of products at the same time.
- Request Body:
```json
{
  "items": [
    {
      "product_id": 1,
      "quantity": 10
    }
  ],
  "user_id": 1
}
```

### Create a New Payment
- URL: http://localhost:8080/payments
- URL: https://ecommerce-management-kwsu.onrender.com/payments
- Method: POST
- Description: Create a new payment
- Request Body:
```json
{
  "cvc": "815",
  "expDate": "0125",
  "hpan": "4405639704015096",
  "order_id": 95,
}
```

### Test Cards

| PAN             | Expire Date | CVC  | Status  |
|-----------------|-------------|------|---------|
| 4405639704015096| 01/25       | 815  | unlock  |
| 5522042705066736| 01/25       | 525  | unlock  |
| 377514500004820 | 01/25       | 4169 | lock    |
| 4003032704547597| 09/20       | 170  | lock    |
| 5578342710750560| 09/20       | 254  | lock    |



- Here is the link to [SWAGGER](https://ecommerce-management-kwsu.onrender.com/swagger/index.htm)

# Swagger: HTTP tutorial for beginners

1. Add comments to your API source code, See [Declarative Comments Format](#declarative-comments-format).

2. Download swag by using:
```sh
go install github.com/swaggo/swag/cmd/swag@latest
```
To build from source you need [Go](https://golang.org/dl/) (1.17 or newer).

Or download a pre-compiled binary from the [release page](https://github.com/swaggo/swag/releases).

3. Run `swag init` in the project's root folder which contains the `main.go` file. This will parse your comments and generate the required files (`docs` folder and `docs/docs.go`).
```sh
swag init
```

  Make sure to import the generated `docs/docs.go` so that your specific configuration gets `init`'ed. If your General API annotations do not live in `main.go`, you can let swag know with `-g` flag.
  ```sh
  swag init -g internal/handler/handler.go
  ```

4. (optional) Use `swag fmt` format the SWAG comment. (Please upgrade to the latest version)

  ```sh
  swag fmt
  ```

## Project Structure

- Makefile: Makefile for building, running, testing, and Docker tasks.
- Dockerfile: Dockerfile for containerizing the application.
- internal/handlers: Contains the HTTP handlers for the API endpoints.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for any improvements.