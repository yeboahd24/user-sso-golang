# User SSO Golang Project

## Endpoints

### Public Routes
- `GET /`: Serves the login page

### Auth Routes
All auth routes are prefixed with `/auth`

- `POST /auth/login`: User login
- `POST /auth/register`: User registration
- `GET /auth/google/login`: Initiates Google OAuth login
- `GET /auth/google/callback`: Handles Google OAuth callback
- `GET /auth/verify`: Verifies user session
- `POST /auth/logout`: User logout

## Setup and Configuration
- The project uses Gin framework for routing
- Configuration is managed using Viper, supporting both config files and environment variables
- Database connection is established using the provided DSN in the config

## Running the Project
The main entry point is in `cmd/main.go`. The server starts on port 8080.

For more details on implementation, refer to the respective handler and service files.
