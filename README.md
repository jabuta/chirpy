
# Chirpy Web Application

## Overview

Chirpy is a web application that allows users to post, view, and manage short messages called "chirps." The application supports user authentication, chirp management, and admin metrics. It is built in Go, using the Chi router for handling HTTP requests and a JSON file for persistent storage.

## Environment Setup

- **Environment Variables:** The application requires setting `JWT_SECRET` for token generation and `POLKA_API` for external API interactions. These should be defined in a `.env` file at the root of the project.

## Initialization

The application is initialized in the `main` function, where the database is set up, routes are defined, and the HTTP server is started. The Chi router is used to define and mount routes for the main application, API, and admin metrics.

## Database

### Schema

The database schema is defined in `database.go` and consists of three main entities:
- **Chirps**: Messages posted by users.
- **Users**: User accounts in the system.
- **Tokens**: Active tokens for user sessions.

### Operations

The `database` package provides functionality for creating, reading, updating, and deleting data related to chirps and users, as well as managing tokens.

## Endpoints

---

### Chirps (`/api/chirps`)

- **GET `/chirps`**: Fetches a list of chirps. This endpoint now accepts an optional `author_id` query parameter to filter chirps by the author. If `author_id` is not provided, all chirps are returned. The endpoint has been updated to directly filter the chirps in the database layer, improving efficiency and simplifying the logic in the API layer.

#### Updated Functionality

- If an `author_id` is provided in the query parameters, the function will only return chirps authored by the specified user ID. If no `author_id` is given, it returns all chirps.
- The database function `GetChirpsList` has been updated to support this filtering by accepting a user ID (`uid`) as a parameter. It locks the database for reading, loads the current database state, and iterates over the chirps to compile a list filtered by the author ID if provided.

This update improves the flexibility of the `getChirps` endpoint, making it more useful for clients that need to display chirps for specific users or all chirps without filtering.

### Example Usage

- To fetch all chirps: `GET /api/chirps`
- To fetch chirps by a specific author: `GET /api/chirps?author_id=123`

This enhancement to the `getChirps` endpoint enables clients to tailor chirp retrieval to their specific needs, whether displaying a user's chirps or the entire chirp feed.

---

### Users (`/api/users`)

- **POST `/users`**: Register a new user with email and password.
- **PUT `/users`**: Update the authenticated user's information.
- **POST `/login`**: Authenticate a user and return an access and a refresh token.
- **POST `/refresh`**: Issue a new access token using the refresh token.
- **POST `/revoke`**: Revoke the user's current access token.

### Admin (`/admin/metrics`)

- **GET `/metrics`**: Provides the number of times the Chirpy application has been accessed.

### Health Check (`/api/healthz`)

- **GET `/healthz`**: Returns a 200 OK status, indicating the application is running.

### Polka Webhooks (`/api/polka/webhooks`)

- **POST `/polka/webhooks`**: Handles external webhook events for user upgrades.

## Middleware

### CORS

- Configured to allow all origins and methods. It is applied globally to the application to enable cross-origin requests.

### Metrics Middleware

- Increments a counter every time the file server is accessed. This is part of the admin metrics functionality.

## Authentication

Authentication is handled using JWT tokens. The application supports token generation, refresh, and revocation. Passwords are hashed before storage, and token-based authentication is required for accessing protected routes.

## Error Handling

Standardized error responses are provided for various failure scenarios, such as unauthorized access, invalid parameters, and database errors.