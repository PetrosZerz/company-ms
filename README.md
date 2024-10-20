# Company Microservice

## Overview
This project is a Company Microservice built using Go, designed to manage company data with support for persistent storage and event publishing.

## Technology Stack
- **Go**: v1.22.0
- **MongoDB**: (Docker Image: `mongo:6.0`)
- **Kafka**: (Docker Image: `wurstmeister/kafka:2.13-2.8.1`)
- **Zookeeper**: (Docker Image: `wurstmeister/zookeeper:3.4.6`)

## Project Structure 
```bash
├── Dockerfile
├── README.md
├── cmd
│   ├── __debug_bin2578255471
│   └── main.go
├── config.toml
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal
│   ├── adapters
│   │   ├── http_handlers.go
│   │   ├── kafka_producer.go
│   │   └── mongo_repository.go
│   ├── application
│   │   ├── company_repository.go
│   │   ├── company_service.go
│   │   ├── errors.go
│   │   └── message_producer.go
│   ├── config
│   │   └── config.go
│   └── model.go
└── tests
    └── company_service_integration_test.go
```

## API Endpoints

### 1. Retrieve All Companies
- **GET** `/v1/companies`
- Description: Retrieve all companies persisted in the database.

### 2. Persist a Company
- **POST** `/v1/companies`
- Description: Create a new company entry.
- **Example Payload**:
    ```json
    {
        "name": "company1",
        "description": "",
        "amount_of_employees": 121,
        "registered": true,
        "type": "NonProfit"
    }
    ```

### 3. Get a Company by ID
- **GET** `/v1/companies/{id}`
- Description: Retrieve a company by providing its ID.

### 4. Update a Company
- **PATCH** `/v1/companies/{id}`
- Description: Update an existing company identified by its ID.
- **Example Payload**:
    ```json
    {
        "amount_of_employees": 121,
        "registered": true
    }
    ```

### 5. Delete a Company
- **DELETE** `/v1/companies/{id}`
- Description: Delete a company entry by its ID.

## Setup and Run
1. Use Docker Compose to set up MongoDB, Kafka, and Zookeeper:
   ```bash 
   docker-compose up -d
   ```
2. Run the Application using go run ./cmd/main.go on the top level  of the project 
    ```bash
    go run ./cmd/main.go 
    ```
3. The application is running at the port `8080`.

 **Note**: A Dockerfile for the application has been created; however, some issues were encountered in building stage that have not yet been resolved.  


## Project Notes

- **Event Publishing**: Kafka is utilized for publishing events during mutate operations.
- **Docker**: The project uses Docker images for MongoDB, Kafka, and Zookeeper, orchestrated through a Docker Compose file.
- **API Design**: The service follows a RESTful architecture throughout the project.
- **Testing**: An integration test is implemented for the create company operation, ensuring functionality with both MongoDB and Kafka.
- **Configuration**: The configuration file `conf.toml` is located at the top level of the project.
- **Authentication**: Note that authentication using JWT has not been implemented in this version of the application.




