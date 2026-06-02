# Go Echo REST API

A RESTful API built with Go, Echo framework, MySQL database, Zap logger, and JWT authentication.

## Project Structure

Based on the provided STRUCTURE.md:

```
GO-ECHO-API/
│
├── config/          → Application configuration setup
├── database/        → Database connection and setup
├── handlers/        → Request handling / business logic
├── logs/            → Application log files
├── middlewares/     → Middleware functions (auth, logging, etc.)
├── models/          → Database models / structs
├── routes/          → API route definitions
├── uploads/         → Uploaded files storage
├── utils/           → Helper functions and JWT utilities
├── zplogger/        → Custom logging implementation
├── main.go          → Entry point of the application
├── Makefile         → Build and automation commands
└── schemasql        → Database schema file or SQL definitions
```

## Features

- **Framework**: Echo v4
- **Database**: MySQL (using schema.sql)
- **Logging**: Uber Zap logger
- **Authentication**: JWT-based authentication
- **Environment Configuration**: Using .env files
- **RESTful API**: Clean route organization

## Setup Instructions

### 1. Prerequisites

- Go 1.25+
- MySQL Server
- Git

### 2. Database Setup

1. Create a MySQL database:
   ```sql
   CREATE DATABASE echo_api;
   ```

2. Update the `.env` file with your MySQL credentials:
   ```env
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=your_mysql_username
   DB_PASSWORD=your_mysql_password
   DB_NAME=echo_api
   SERVER_PORT=8080
   JWT_SECRET=your-secret-key-change-this-in-production
   ```

3. Run the schema.sql to create tables:
   ```bash
   mysql -u your_mysql_username -p echo_api < schema.sql
   ```

### 3. Running the Application

```bash
# Install dependencies
go mod tidy

# Build the application
go build -o go-echo-api .

# Run the application
.\go-echo-api.exe
```

Or run directly:
```bash
go run main.go
```

### 4. API Endpoints

#### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration

*(Additional endpoints for users, accounts, pockets, transactions to be implemented)*

### 5. Environment Variables

Copy `.env.example` to `.env` and update the values:

| Variable | Description | Default |
|----------|-------------|---------|
| `ENV_MODE` | Application environment | `development` |
| `SERVICE_NAME` | Service name for logging | `go-echo-api` |
| `HOST` | Server listen address | `0.0.0.0` |
| `PORT` | HTTP server port | `8000` |
| `BASE_PATH` | URL path prefix | `api` |
| `TZ` | Timezone | `Asia/Bangkok` |
| `DB_URI` | MySQL connection URI | `mysql://demo:1234@localhost:3306/games_store` |
| `ENCRYPTION_KEY` | Key for symmetric encryption | — |
| `JWT_PRIVATE_KEY` | RSA private key for JWT signing | — |
| `JWT_PUBLIC_KEY` | RSA public key for JWT verification | — |
| `UPLOAD_LIMIT_SIZE` | Max upload size (MB) | `10` |
| `IMAGE_COMPRESS_LEVEL` | Image compression level (0–100) | `70` |
| `REDIS_URI` | Redis connection URI | `redis://localhost:6379` |
| `MQTT_HOST` | MQTT broker host | `localhost` |
| `MQTT_PORT` | MQTT broker port | `8084` |
| `MQTT_PROTOCOL` | MQTT protocol (wss/tcp) | `wss` |
| `MQTT_USER` | MQTT username | `demo` |
| `MQTT_PASSWORD` | MQTT password | `1234` |
| `MQTT_PATH` | MQTT path | `mqtt` |
| `MQTT_TOPIC` | MQTT topic prefix | `dev` |
| `LIMIT_MAX_BALANCE` | Enforce max balance cap | `false` |
| `LOG_LEVEL` | Logging level (debug/info/warn/error) | `info` |

### 6. Logging

Logs are written to:
- `logs/demo-rest-api.log` - General application logs
- `logs/demo-rest-api-error.log` - Error logs

### 7. Project Modules

- **config**: Handles loading configuration from environment variables
- **database**: Manages MySQL database connections
- **handlers**: Contains HTTP request handlers for different resources
- **middlewares**: Custom middleware including request logging
- **models**: Data structures representing database tables
- **routes**: API route definitions and registration
- **utils**: Utility functions including JWT token handling
- **zplogger**: Wrapper around Uber Zap logger for structured logging

## Implementation Details

### Authentication
- Uses JWT (JSON Web Tokens) for stateless authentication
- Passwords are hashed using bcrypt before storage
- Tokens expire after 24 hours

### Logging
- Structured logging with Uber Zap
- Request/response logging middleware
- Separate log files for general logs and errors

### Database
- Connection pooling through standard library database/sql
- Automatic connection testing on startup
- Parameterized queries to prevent SQL injection

## Future Enhancements

1. Implement remaining CRUD endpoints for:
   - Users
   - Accounts
   - Pockets
   - Transactions

2. Add role-based access control (RBAC)

3. Implement input validation using validator library

4. Add API documentation with Swagger/OpenAPI

5. Implement unit and integration tests

6. Add Docker support for containerized deployment

7. Implement rate limiting middleware

8. Add pagination to list endpoints

## License

MIT