# Chirpy API Specification
This document details the API endpoints for the Chirpy web application, focusing on authentication, user management, and chirps interaction. It specifies required request formats, response structures, and the need for authentication tokens for certain operations.

## Authentication

### User Login
- **POST** `/api/login`
- **Auth Required**: No
- **Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```
- **Response**: Contains access and refresh tokens, along with user details.

### Refresh Token
- **POST** `/api/refresh`
- **Auth Required**: Yes (via refresh token)
- **Purpose**: Obtain a new access token using a valid refresh token.
- **Response**:
  ```json
  {
    "token": "new_access_token"
  }
  ```

### Revoke Token
- **POST** `/api/revoke`
- **Auth Required**: Yes
- **Purpose**: Invalidate the current access token.
- **Response**: Confirmation message.

## Endpoints

### Chirps

#### Get All Chirps
- **GET** `/api/chirps`
- **Auth Required**: Optional
- **Query Params**: `author_id` (integer, optional)
- **Response**: Array of chirp objects.

#### Create Chirp
- **POST** `/api/chirps`
- **Auth Required**: Yes
- **Body**:
  ```json
  {
    "body": "Chirp content"
  }
  ```
- **Response**: Chirp object.

### Users

#### Register User
- **POST** `/api/users`
- **Auth Required**: No
- **Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "password123",
    "is_chirpy_red": false
  }
  ```
- **Response**: User object.

## Response Structures

### Chirp Object
```json
{
  "id": 1,
  "body": "Chirp content",
  "authorID": 123
}
```

### User Object
```json
{
  "id": 123,
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

### Login Response
```json
{
  "id": 123,
  "email": "user@example.com",
  "is_chirpy_red": false,
  "token": "access_token",
  "refresh_token": "refresh_token"
}
```

## Notes

- Protected endpoints require an `Authorization` header with a valid JWT token presented as `Bearer <token>`.
- The API uses HTTP status codes to indicate the success or failure of requests.
- Error responses are returned in a standard JSON format with an error message.
