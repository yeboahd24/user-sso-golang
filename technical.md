Context and scope

This project is a User Single Sign-On (SSO) system implemented in Go. It provides authentication functionality, including local user registration and login, as well as Google OAuth integration. The system uses Gin as the web framework, PostgreSQL for data storage, and JWT for session management.

Goals and non-goals

Goals:

Provide a secure user authentication system

Support both local and Google OAuth authentication methods

Implement session management using JWT

Offer a configurable setup through YAML configuration and environment variables

Non-goals:

Implement other OAuth providers (e.g., Facebook, Twitter)

Provide a full-fledged user management system

Handle authorization (focus is on authentication only)

Design

System-context-diagram

sequenceDiagram
    participant User
    participant Frontend
    participant SSO_Service
    participant Database
    participant Google_OAuth

    User->>Frontend: Interacts with UI
    Frontend->>SSO_Service: Sends authentication request
    SSO_Service->>Database: Validates user credentials
    SSO_Service->>Google_OAuth: Initiates OAuth flow (if selected)
    Google_OAuth-->>SSO_Service: Returns OAuth token
    SSO_Service-->>Frontend: Returns JWT token
    Frontend-->>User: Authenticates user

APIs

The system exposes the following main API endpoints:

POST /auth/register: User registration

POST /auth/login: User login

GET /auth/google/login: Initiate Google OAuth flow

GET /auth/google/callback: Handle Google OAuth callback

GET /auth/verify: Verify user session

POST /auth/logout: User logout

Data storage

The system uses PostgreSQL for data storage. The main entity is the User, which likely includes fields such as ID, email, password hash, and potentially OAuth-related information.

Code and pseudo-code

The project is structured into several packages:

cmd: Contains the main application entry point

config: Handles configuration loading and management

database: Manages database connections

handler: Implements HTTP request handlers

repository: Handles data access and storage

route: Sets up API routes

service: Implements business logic

util: Provides utility functions (e.g., JWT generation)

Degree of constraint

This project has a moderate degree of constraint. It uses established libraries and patterns (Gin, GORM, JWT) but allows for flexibility in implementation details and future extensions.

Alternatives considered

Using session-based authentication instead of JWT

Pros: Easier to invalidate sessions, potentially more secure for certain use cases

Cons: Requires server-side storage, less scalable

Implementing custom OAuth flow instead of using a library

Pros: More control over the OAuth process

Cons: Increased complexity, higher risk of security vulnerabilities

The chosen design (JWT-based authentication with library-supported OAuth) provides a good balance between security, scalability, and ease of implementation.

Now, regarding the sample of the config.yaml file, based on the provided config/config.go file, here's an example of what the config.yaml might look like:
```json
server:
  port: "8080"
  mode: "development"
  timeout: 30

database:
  host: "localhost"
  port: "5432"
  user: "your_db_user"
  password: "your_db_password"
  dbname: "your_db_name"
  sslmode: "disable"
  timezone: "UTC"

oauth:
  google:
    clientid: "your_google_client_id"
    clientsecret: "your_google_client_secret"
    redirecturl: "http://localhost:8080/auth/google/callback"

jwt:
  secret: "your_jwt_secret"
  expiresinhrs: 24
  refreshsecret: "your_refresh_token_secret"

redis:
  host: "localhost"
  port: "6379"
  password: ""
  db: 0
```
This YAML configuration file corresponds to the structure defined in the Config struct in config/config.go. It includes settings for the server, database connection, Google OAuth, JWT, and Redis. Users of this project would need to replace the placeholder values (e.g., "your_db_user", "your_google_client_id") with their actual configuration details.
