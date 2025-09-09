# ScribeHost Authentication Module

## Overview

This module implements a secure authentication system for the ScribeHost application, featuring email verification with OTP (One-Time Password), JWT-based authentication, and refresh tokens. The system follows a clean MVC (Model-View-Controller) architecture for maintainability and scalability.

## Authentication Flow

1. **Registration**:
   - User submits registration request with email and password
   - System creates user account with `email_verified = false`
   - User receives verification email with OTP

2. **Email Verification**:
   - User receives OTP via email
   - User submits OTP to `/api/v1/auth/verify-otp` endpoint
   - System verifies OTP and sets `email_verified = true`
   - User is now verified and can log in

3. **Login**:
   - User submits email and password
   - System authenticates user and issues JWT access token
   - Access token is sent as `accessToken` cookie (HttpOnly, Secure)
   - Refresh token is sent as `refreshToken` cookie (HttpOnly, Secure)

4. **Session Management**:
   - User can refresh access token using `/api/v1/auth/refresh`
   - User can log out using `/api/v1/auth/logout`
   - Session is ended by revoking refresh tokens

## API Endpoints

| Method | Endpoint | Description | Authentication Required |
|--------|----------|-------------|-------------------------|
| POST | `/api/v1/auth/register` | Register a new user | No |
| POST | `/api/v1/auth/login` | Log in a user | No |
| POST | `/api/v1/auth/refresh` | Refresh access token | Yes (Access Token) |
| POST | `/api/v1/auth/logout` | Log out user | Yes (Access Token) |
| POST | `/api/v1/auth/verify-otp` | Verify email OTP | No |
| POST | `/api/v1/auth/resend-otp` | Resend OTP for email verification | No |

## Rate Limiting

The authentication system implements rate limiting using the Token Bucket algorithm to prevent abuse and protect against brute force attacks. Rate limits are applied to the following endpoints:

| Endpoint | Rate Limit | Identifier |
|----------|------------|------------|
| `/api/v1/auth/register` | 1 request per hour | IP address |
| `/api/v1/auth/login` | 5 requests per hour | IP address |
| `/api/v1/auth/verify-otp` | 3 requests per hour | Email address |
| `/api/v1/auth/resend-otp` | 3 requests per hour | IP address |

This rate limiting is implemented using Redis for storage, ensuring it works correctly even if multiple instances of the application are running. The rate limits are hard-coded in the application and do not require any additional configuration.

## Architecture

The authentication module follows a clean MVC (Model-View-Controller) architecture:

### Models (Entity Layer)
- `entity/user_entity.go`: User model with email, password hash, verification status
- `entity/refresh_token_entity.go`: Refresh token model for session management

### Controllers (HTTP Layer)
- `controller/auth_controller.go`: Handles HTTP requests and responses
- `controller/user_controller.go`: Handles user-related requests

### Services (Business Logic)
- `service/auth_service.go`: Contains authentication business logic
- `service/user_service.go`: Contains user business logic

### Repositories (Data Access)
- `repository/auth_repo.go`: Handles authentication data operations
- `repository/user_repo.go`: Handles user data operations
- `repository/refresh_token_repo.go`: Handles refresh token operations
- `repository/otp_redis_repo.go`: Handles OTP storage in Redis

### DTOs (Data Transfer Objects)
- `dto/auth_dto.go`: DTOs for authentication requests/responses
- `dto/user_dto.go`: DTOs for user-related requests/responses

### Helper Functions
- `helper/otp_helper.go`: OTP generation and verification
- `helper/email_helper.go`: Email sending functionality
- `helper/jwt_helper.go`: JWT token generation and validation
- `helper/password_helper.go`: Password hashing and comparison
- `helper/redis_client_helper.go`: Redis connection helper

### Middleware
- `middleware/jwt_auth_middleware.go`: JWT authentication middleware
- `middleware/logger_middleware.go`: Request logging middleware
- `middleware/rate_limiter.go`: Rate limiting middleware

## Setup

1. **Environment Configuration**:
   - Create a `.env` file in the project root with the following configuration (use dummy data for testing, replace with your actual credentials for production):

```
# Application Configuration
PORT=4000

# Database Configuration
DB_USERNAME=admin
DB_PASSWORD=admin
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=scribe_host

# Redis Configuration (for OTP storage)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Goose Configuration (for database migrations)
GOOSE_DRIVER=mysql
GOOSE_DBSTRING="${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?parseTime=true"

# Email Configuration (for OTP verification)
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USER=your_email@gmail.com
EMAIL_SMTP_PASSWORD=your_app_password
EMAIL_FROM=your_email@gmail.com
```

**Explanation of each variable**:

- `PORT=4000`: The port the application listens on (default for development)
- `DB_USERNAME=admin`: Database username (e.g., "root" or your custom username)
- `DB_PASSWORD=admin`: Database password (use a secure password in production)
- `DB_HOST=127.0.0.1`: Database host (localhost for local development)
- `DB_PORT=3306`: Database port (default for MySQL)
- `DB_NAME=scribe_host`: Name of your database (will be created during migrations)
- `REDIS_HOST=localhost`: Redis server host
- `REDIS_PORT=6379`: Redis server port
- `REDIS_PASSWORD=`: Redis server password (empty for default configuration)
- `REDIS_DB=0`: Redis database number to use
- `GOOSE_DRIVER=mysql`: Database driver for migration tool (should match your database)
- `GOOSE_DBSTRING`: Connection string for the migration tool (automatically generated from other variables)
- `EMAIL_SMTP_HOST=smtp.gmail.com`: SMTP server host (e.g., "smtp.gmail.com" for Gmail)
- `EMAIL_SMTP_PORT=587`: SMTP server port (587 for TLS, 465 for SSL)
- `EMAIL_SMTP_USER=your_email@gmail.com`: Email address used to send OTPs
- `EMAIL_SMTP_PASSWORD=your_app_password`: Password for the email account (use an app password for Gmail)
- `EMAIL_FROM=your_email@gmail.com`: Email address that appears as the sender

> **Important Note**: For Gmail, you'll need to generate an "App Password" instead of using your regular account password. Go to your Google Account > Security > App Passwords to generate one. This is required because Gmail blocks sign-in attempts from apps that don't use 2-Step Verification.

> **Redis Note**: Redis is used to store OTPs temporarily with automatic expiration. The Redis server must be running for OTP functionality to work. You can install Redis using your package manager or Docker.

2. **Database**:
   - Run migrations to set up database schema
   - Use `make migrate-up` to apply migrations

3. **Redis**:
   - Install and start Redis server (see [Redis installation guide](https://redis.io/docs/getting-started/installation/))
   - Verify Redis is running with `redis-cli ping` (should return "PONG")

4. **Running the Application**:
   - Use `make dev` for development with live reload
   - Use `make up` to run the application with Docker

## Future Features

- OAuth 2.0 integration for social logins
- Account recovery flow (password reset)
- Two-factor authentication (2FA) with TOTP

> **Note**: This documentation is part of the `feature/auth` branch and will be updated as new features are added to the authentication module.
