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

## Environment Configuration

The application uses two environment files for different environments:

1. **`.env.local`** - For development (default)
2. **`.env.production`** - For production

### Example `.env.local`:
```
# .env.local
PORT=4000

DB_USERNAME=admin
DB_PASSWORD=admin
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=scribe_host

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

GOOSE_DRIVER=mysql
GOOSE_DBSTRING="${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?parseTime=true"

EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USER=your_email@gmail.com
EMAIL_SMTP_PASSWORD=your_app_password
EMAIL_FROM=your_email@gmail.com
```

### Example `.env.production`:
```
# .env.production
PORT=8080

DB_USERNAME=scribehost_user
DB_PASSWORD=scribepassword123
DB_HOST=db
DB_PORT=3306
DB_NAME=scribehost

REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=redisPassword
REDIS_DB=0

GOOSE_DRIVER=mysql
GOOSE_DBSTRING="${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?parseTime=true"

EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USER=production@scribehost.com
EMAIL_SMTP_PASSWORD=production_app_password
EMAIL_FROM=production@scribehost.com
```

**Key Differences**:
- Development uses `localhost` for database and Redis
- Production uses Docker service names (`db` and `redis`)
- Production uses different database credentials

## Setup

1. **Environment Configuration**:
   - Create `.env.local` for development (default)
   - Create `.env.production` for production
   - The application automatically loads the appropriate file based on `APP_ENV` (default: `development`)

2. **Database**:
   - Run migrations to set up database schema
   - Use `make migrate-up` to apply migrations

3. **Redis**:
   - For development: Run Redis locally (port 6379)
   - For production: The Docker Compose setup includes the Redis service

4. **Docker Compose for Production**:
   - Use `make up` to run the application with Docker Compose (production environment)
   - The Docker Compose file is configured for production deployment
   - It includes the `redis` service and sets `APP_ENV=production`

5. **Running the Application**:
   - For development: `make dev` (uses `.env.local`)
   - For production: `make up` (uses `.env.production` and Docker Compose)

## Future Features

- OAuth 2.0 integration for social logins
- Account recovery flow (password reset)
- Two-factor authentication (2FA) with TOTP

> **Note**: This documentation is part of the `feature/auth` branch and will be updated as new features are added to the authentication module.
