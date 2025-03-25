# Movie API

## Overview

This API allows users to manage movies, including creating, updating, retrieving, and deleting movie records.

### Default Admin Credentials:

- **Username**: `admin`
- **Password**: `password123`

---

## Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) (1.22 or later)
- [Docker](https://www.docker.com/)
- [PostgreSQL](https://www.postgresql.org/)

### Running Locally

#### 1. Clone the Repository

```sh
$ git clone https://github.com/boriyevmahmud/itv-task
$ cd itv-task
```

#### 2. Set Up Environment Variables

Create a `.env` file and configure the database:

```env
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=movies_db
JWT_SECRET=your_jwt_secret
```

#### 3. Run the Application

```sh
$ go run main.go
```

The API will start on `http://localhost:8080`.

---

## Running with Docker

### Build and Run the Container

```sh
$ docker compose up --build
```

---

## API Documentation (Swagger)

Once the server is running, open Swagger UI:

```
http://localhost:8080/swagger/index.html
```

---

## API Endpoints

### Authentication

#### Login (Get JWT Token)

**POST** `/auth/login`

##### Request Body:

```json
{
  "username": "admin",
  "password": "password123"
}
```

##### Response:

```json
{
  "token": "eyJhbGciOiJIUzI1..."
}
```

### Movies

#### Create Movies in Bulk

**POST** `/movies`

##### Request Body:

```json
{
  "movies": [
    {
      "title": "Inception",
      "director": "Christopher Nolan",
      "year": 2010,
      "plot": "A thief who enters people's dreams."
    }
  ]
}
```

##### Response:

```json
{
  "message": "Movies created"
}
```

#### Get Movie by ID

**GET** `/movies/{id}`

##### Response:

```json
{
  "id": 1,
  "title": "Inception",
  "director": "Christopher Nolan",
  "year": 2010,
  "plot": "A thief who enters people's dreams.",
  "created_at": "2025-03-22T15:04:05Z",
  "updated_at": "2025-03-22T15:04:05Z"
}
```

---

## Additional Notes

- Ensure that the database is running before starting the application.
- Use the JWT token received from login in the `Authorization` header as `Bearer <token>` for protected routes.
