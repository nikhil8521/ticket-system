# Ticket Management System Backend

This project implements a production-ready backend for a Ticket Management System using Golang, Gin, GORM, and SQLite. It includes user authentication with JWT, ticket management functionalities, and is designed with a clean architecture.

## Table of Contents
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Setup and Running Locally](#setup-and-running-locally)
  - [Go Installation](#go-installation)
  - [Environment Variables](#environment-variables)
  - [Running with Go](#running-with-go)
  - [Running with Docker](#running-with-docker)
- [API Documentation](#api-documentation)
  - [Authentication](#authentication)
  - [Ticket APIs](#ticket-apis)
- [Deployment](#deployment)
- [Testing](#testing)

## Features
- User Registration and Login
- JWT-based Authentication
- Secure Password Hashing with bcrypt
- Ticket Creation, Listing, and Viewing (owner-restricted)
- Ticket Status Updates with Validated Transitions
- RESTful API Design
- Clean Architecture
- Dockerized Application

## Tech Stack
- **Language**: Go (latest stable version)
- **Web Framework**: [Gin Gonic](https://gin-gonic.com/)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: SQLite
- **Authentication**: JWT (JSON Web Tokens)
- **Password Hashing**: bcrypt
- **Dependency Management**: Go Modules
- **Containerization**: Docker, Docker Compose
- **Environment Variables**: godotenv

## Project Structure
```
ticket-system/
├── cmd/
│   └── server/             # Main application entry point
├── config/                 # Configuration files (if any, currently not used explicitly)
├── controllers/            # Handles HTTP requests and responses
├── middleware/             # Custom middleware, e.g., authentication
├── models/                 # Defines database models (User, Ticket)
├── repository/             # Data access layer (interfaces and implementations)
├── services/               # Business logic layer (interfaces and implementations)
├── routes/                 # Defines API routes
├── utils/                  # Utility functions, e.g., JWT, password hashing
├── database/               # Database connection and migration setup
├── docs/                   # API documentation (Postman collection will be here)
├── Dockerfile              # Docker build instructions
├── docker-compose.yml      # Docker Compose configuration
├── README.md               # Project README
├── .env.example            # Example environment variables
└── go.mod                  # Go module definition
```

## Setup and Running Locally

### Go Installation
Ensure you have Go installed on your system. You can download it from the [official Go website](https://golang.org/doc/install).

### Environment Variables
Create a `.env` file in the root directory of the project based on `.env.example`:

```bash
cp .env.example .env
```

Edit `.env` with your desired values:

```
PORT=8080
JWT_SECRET=your_super_secret_key_change_this_in_production
DB_PATH=ticket.db
```

- `PORT`: The port on which the server will run.
- `JWT_SECRET`: A secret key used for signing JWT tokens. **Change this to a strong, unique value in production.**
- `DB_PATH`: The path to the SQLite database file.

### Running with Go
1. Navigate to the project root:
   ```bash
   cd ticket-system
   ```
2. Download dependencies:
   ```bash
   go mod tidy
   ```
3. Run the application:
   ```bash
   go run ./cmd/server/main.go
   ```
   The server will start on the port specified in your `.env` file (default: `8080`).

### Running with Docker
1. Navigate to the project root:
   ```bash
   cd ticket-system
   ```
2. Build the Docker image:
   ```bash
   docker build -t ticket-system .
   ```
3. Run the Docker container:
   ```bash
   docker run -p 8080:8080 ticket-system
   ```
   Alternatively, you can use `docker-compose`:
   ```bash
   docker-compose up --build
   ```
   The application will be accessible at `http://localhost:8080`.

## API Documentation
The API endpoints are designed to be RESTful. All protected routes require a `Authorization: Bearer <token>` header.

### Authentication

#### Register a new user
- **Endpoint**: `POST /auth/register`
- **Body**:
  ```json
  {
    "name": "John Doe",
    "email": "john.doe@example.com",
    "password": "password123"
  }
  ```
- **Validation**: `name` required, `email` valid and unique, `password` minimum 6 characters.
- **Response**: `201 Created` on success, `400 Bad Request` on validation error or email already exists.

#### Login user
- **Endpoint**: `POST /auth/login`
- **Body**:
  ```json
  {
    "email": "john.doe@example.com",
    "password": "password123"
  }
  ```
- **Validation**: `email` required, `password` required.
- **Response**: `200 OK` with JWT token on success, `401 Unauthorized` on invalid credentials.
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```

### Ticket APIs

#### Health Check
- **Endpoint**: `GET /health`
- **Response**: `200 OK`
  ```json
  {
    "status": "ok"
  }
  ```

#### Create a new ticket
- **Endpoint**: `POST /tickets`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "title": "Bug in login page",
    "description": "Users are unable to log in after recent deployment."
  }
  ```
- **Validation**: `title` required, `description` required.
- **Response**: `201 Created` with the new ticket object.

#### Get all tickets for the logged-in user
- **Endpoint**: `GET /tickets`
- **Headers**: `Authorization: Bearer <token>`
- **Response**: `200 OK` with an array of tickets owned by the user.

#### Get a specific ticket by ID
- **Endpoint**: `GET /tickets/:id`
- **Headers**: `Authorization: Bearer <token>`
- **Response**: `200 OK` with the ticket object if owned by the user. `403 Forbidden` if not owned by the user. `404 Not Found` if ticket does not exist.

#### Update ticket status
- **Endpoint**: `PATCH /tickets/:id/status`
- **Headers**: `Authorization: Bearer <token>`
- **Body**:
  ```json
  {
    "status": "in_progress"
  }
  ```
- **Supported Statuses**: `open`, `in_progress`, `closed`.
- **Status Transition Rules**:
  - `open` → `in_progress`
  - `in_progress` → `closed`
  - `closed` cannot be reopened.
  - `open` cannot directly become `closed`.
- **Response**: `200 OK` on successful update. `400 Bad Request` on invalid status or transition. `403 Forbidden` if not owned by the user. `404 Not Found` if ticket does not exist.

## Deployment
This application is designed for deployment on platforms like Render. The `Dockerfile` and `docker-compose.yml` facilitate containerized deployment. The health endpoint (`/health`) should be publicly accessible.

## Testing
A Postman Collection will be provided in the `docs/` directory for easy testing of all API endpoints.
