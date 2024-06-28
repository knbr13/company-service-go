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

To run the application without Docker and Docker Compose:

1. Ensure you have Go 1.22 or later installed
2. Install dependencies:
```bash
go mod tidy
```
3. Set up your local MySQL database and Kafka instance
4. Update the `.env` file with your local configuration
5. Run the migrations, you have to have the golang-migrate tool, if not, download it:
    - `go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
6. Run the application:
```bash
go run ./cmd/api/
```

To run the application with Docker and Docker Compose:
1. Make sure you have Docker and Docker Compose installed.
2. Update the `.env` file with your local configuration.
3. Run docker compose up (the command is `docker compose` or `docker-compose`, depends on the installation):
```bash
docker compose up --build
```
4. Run the migrations, the golang-migrate tool is already installed for you in the container, check the `Dockerfile` for reference, so just run the migration command from the container:
```bash
docker compose exec -T app /usr/local/bin/migrate -path /migrations -database "mysql://root:$(MYSQL_ROOT_PASSWORD)@tcp(db-mysql:3306)/$(MYSQL_DATABASE)" up
```
