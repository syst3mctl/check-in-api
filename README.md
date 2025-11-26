# Check-In/Out Service API

A robust REST API for managing employee attendance, organizations, shifts, and tasks. Built with Go, PostgreSQL, and Docker.

## Features

- **Authentication**: User registration and login with JWT-based authentication.
- **Organization Management**: Create organizations, invite employees, and manage roles (OWNER, MANAGER, EMPLOYEE). Owners and Managers can view, update, and remove employees.
- **Shift & Group Management**: Define shifts with specific working hours and assign users to groups.
- **Attendance Tracking**:
  - General Check-in/out.
  - Task-based Check-in (with optional geofencing).
  - Late arrival detection.
- **Reporting**: Generate performance reports for groups.
- **Swagger Documentation**: Interactive API documentation.

## Tech Stack

- **Language**: Go (Golang)
- **Framework**: [Chi](https://github.com/go-chi/chi)
- **Database**: PostgreSQL
- **Driver**: [pgx](https://github.com/jackc/pgx)
- **Authentication**: JWT (JSON Web Tokens)
- **Validation**: [go-playground/validator](https://github.com/go-playground/validator)
- **Documentation**: [Swagger](https://github.com/swaggo/swag)
- **Containerization**: Docker & Docker Compose

## Prerequisites

- [Go](https://go.dev/) 1.22+
- [Docker](https://www.docker.com/) & Docker Compose
- [Make](https://www.gnu.org/software/make/) (optional, for using the Makefile)

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/syst3mctl/check-in-api.git
cd check-in-api
```

### 2. Environment Setup

Create a `.env` file in the root directory (or use the default values in `docker-compose.yml` for local dev):

```env
PORT=8080
DATABASE_URL=postgres://user:password@localhost:5433/checkin_db?sslmode=disable
JWT_SECRET=your_super_secret_key
```

### 3. Start Infrastructure

Start the PostgreSQL database using Docker Compose:

```bash
make up
```

### 4. Run Migrations

Apply database migrations to set up the schema:

```bash
make migrate
```

### 5. Run the Application

Start the API server:

```bash
make run
```

The server will start on port `8080` (or the port specified in `.env`).

## API Documentation

The API is documented using Swagger. Once the server is running, access the interactive documentation at:

[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

You can use the "Authorize" button to enter your Bearer token for testing protected endpoints.

## Development

### Common Commands

- `make build`: Build the binary.
- `make run`: Run the application.
- `make up`: Start Docker containers.
- `make down`: Stop Docker containers.
- `make migrate`: Run database migrations.
- `make swag`: Regenerate Swagger documentation.
- `make test`: Run tests.

### Project Structure

```
.
├── cmd/
│   └── server/         # Application entry point
├── internal/
│   ├── adapter/        # Database adapters (PostgreSQL)
│   ├── api/            # HTTP handlers, middleware, and router
│   ├── config/         # Configuration loading
│   ├── core/           # Domain logic (services, domain models, ports)
│   └── pkg/            # Utility packages (logger, validator, response)
├── migrations/         # SQL migration files
├── docs/               # Swagger documentation files
├── docker-compose.yml  # Docker Compose configuration
└── Makefile            # Development commands
```

## License

[MIT](LICENSE)
