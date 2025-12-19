# User Service

## Overview
The **User Service** is a core microservice in the E-Commerce ecosystem responsible for user management and authentication. It handles user registration, secure login via JWT (JSON Web Tokens), and profile management.

Built with **Django** and **Django REST Framework (DRF)**, it is designed to be lightweight, secure, and easily scalable.

## Tech Stack
-   **Framework**: Django 4.2, Django REST Framework
-   **Authentication**: SimpleJWT (Stateless JWT Authentication)
-   **Database**: PostgreSQL
-   **Containerization**: Docker, Docker Compose
-   **Documentation**: OpenAPI 3.0 (via `drf-spectacular`)

## Features
-   **User Registration**: Public endpoint for creating new user accounts.
-   **Authentication**:
    -   **Login**: Obtain Access and Refresh tokens.
    -   **Token Refresh**: Rotate access tokens securely.
    -   **Token Verify**: Check token validity.
-   **Profile Management**: Retrieve and update user profile details (includes address management).
-   **Role-Based Access**: Distinguishes between `USER` and `ADMIN` roles.

## API Endpoints

| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| **POST** | `/api/auth/register/` | Register a new user | No |
| **POST** | `/api/token/` | Login (Obtain JWT pair) | No |
| **POST** | `/api/token/refresh/` | Refresh access token | No |
| **GET** | `/api/auth/profile/` | Get current user's profile | **Yes** |
| **PUT/PATCH** | `/api/auth/profile/` | Update profile information | **Yes** |

*> Note: When running through the Broker Service (Gateway), these endpoints are prefixed with `/user-service` (e.g., `http://localhost/user-service/api/auth/profile/`).*

## Local Development

### Prerequisites
-   Docker and Docker Compose

### Running with Docker (Recommended)
This service is designed to run as part of the larger microservice mesh. To run it:

1.  Navigate to the project root directory.
2.  Run the full stack:
    ```bash
    docker-compose up --build
    ```
3.  The service will start on port `8000` (internal) and `8001` (external access for debugging).

### Environment Variables
The service relies on the following environment variables (defined in `docker-compose.yml`):

-   `DEBUG`: `True`/`False`
-   `SECRET_KEY`: Django secret key
-   `DATABASE_URL`: Connection string for PostgreSQL (e.g., `postgres://postgres:password@postgres:5432/user_service_db`)
-   `ALLOWED_HOSTS`: Comma-separated list of allowed hosts.

## Database
The service uses a **PostgreSQL** database named `user_service_db`.
Migrations are automatically applied on container startup via the root `docker-compose.yml` command override.

## Testing
To run tests inside the container:
```bash
docker-compose exec user-service python manage.py test
```
