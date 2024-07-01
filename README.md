# Company Service

This is a microservice for handling company data, providing CRUD operations for company records.

## Features

- Register, login as a user
- Create, Read, Update, and Delete company records
- JWT Authentication
- Event production for mutating operations
- Dockerized application and services
- MySQL database
- Kafka for event streaming
- Graceful shutdown
- Automatic linting check as you push to remote or create a PR to the `main` branch

## Prerequisites

- Docker
- Docker Compose

## Getting Started

1. Clone the repository
2. Copy the `.env.example` file to `.env` and adjust the values as needed:
```bash
cp .env.example .env
```
3. Build and run the services using Docker Compose:
```bash
docker-compose up --build
```
This will start the following services:
- MySQL database
- Kafka
- Zookeeper
- The company service application

4. The application will be available at `http://localhost:8080`

## Environment Variables

The following environment variables are used in the project:

- `JWT_KEY`: Secret key for JWT token generation and validation
- `KAFKA_BROKER`: Kafka broker address
- `MYSQL_DATABASE`: MySQL database name
- `MYSQL_ROOT_PASSWORD`: MySQL root password
- `MYSQL_ALLOW_EMPTY_PASSWORD`: Whether to allow empty MySQL password
- `DB_DSN`: Database connection string (not needed when using Docker Compose)

## API Endpoints

### Public Endpoints (No Authentication Required)

- `POST /register`: Register a new account to get an Auth Token (JWT).
- `POST /login`: Login if you already have an account, you will get back an Auth Token.

- `GET /companies/{id}`: Get a single company

### Protected Endpoints (JWT Authentication Required)

- `POST /companies`: Create a new company
- `PATCH /companies/{id}`: Update an existing company
- `DELETE /companies/{id}`: Delete a company

## Authentication

The service uses JWT for authentication. To access protected endpoints, include a valid JWT token in the Authorization header of your requests.

## Development

1. Ensure you have Docker and Docker Compose installed
2. Run:
```bash
docker compose up --build
```
If this is the first time you run `docker compose up`, you need to append the `--build` flag in order to build the Go project using the `Dockerfile`.
Otherwise you can just run it without the `--build` flag.

Note: depending on how you have installed Docker Compose, you might need to write the command in this way:
```bash
docker-compose up --build
```
(it depends whether you have installed it as a plugin or a standalone tool).
