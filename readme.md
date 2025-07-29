# Meta Sound Tools Backend

A Go-based REST API backend for an audio metadata management application. This service handles user authentication, license management, and application distribution for [Meta Sound Tools](https://metasoundtools.com).

## Features

- User registration and authentication with AWS Cognito integration
- License key generation and validation
- JWT-based authentication for desktop applications
- Email notifications for beta users
- Application release management and binary distribution
- PostgreSQL database with migration support

## Tech Stack

- **Language**: Go 1.24.1
- **Database**: PostgreSQL with SQLC for type-safe queries
- **Authentication**: JWT tokens with ECDSA signing, AWS Cognito
- **Email**: AWS SES for transactional emails
- **Deployment**: Docker, Fly.io
- **Database Migrations**: Goose

## Setup

### Prerequisites

- Go 1.24.1 or later
- PostgreSQL
- AWS account (for SES and Cognito)

### Environment Variables

Create a `.env` file with:

```
PORT=8080
API_KEY=your_api_key
DB_URL=your_postgres_connection_string
DEV_DB_URL=your_dev_postgres_connection_string
COGNITO_POOL_ID=your_cognito_pool_id
AWS_ACCESS_KEY_ID=your_aws_access_key
AWS_SECRET_ACCESS_KEY=your_aws_secret_key
ENV=dev
```

### Database Setup

1. Run migrations:

```bash
cd db
./migrateup.sh
```

### Running Locally

```bash
go mod download
go run main.go
```

### Docker

```bash
docker-compose up
```

## Project Structure

- `/config` - Configuration management
- `/db` - Database migrations and SQLC generated code
- `/email` - Email service implementation
- `/jwt` - JWT token creation and validation
- `/lambda` - AWS Lambda functions for Cognito integration
- `/licence` - License validation logic
- `/store` - In-memory verification store
- `/utils` - Shared utilities
