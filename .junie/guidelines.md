# Hobigon API Server Guidelines

## Project Overview
Hobigon is a playground API server project running on Kubernetes, implemented in Go. The project serves as a practical implementation of Domain-Driven Design (DDD) principles and clean architecture patterns.

## Project Structure
The project follows a clean architecture pattern with the following key directories:

```
.
├── app/                    # Main application code
│   ├── domain/            # Domain layer (core business logic)
│   │   ├── model/         # Domain models
│   │   ├── repository/    # Repository interfaces
│   │   └── gateway/       # External service interfaces
│   ├── infra/            # Infrastructure layer
│   │   ├── dao/          # Data Access Objects
│   │   ├── db/           # Database configurations
│   │   └── dto/          # Data Transfer Objects
│   ├── presentation/     # Presentation layer
│   │   ├── rest/         # REST API handlers
│   │   ├── graphql/      # GraphQL API handlers
│   │   └── cli/          # Command-line interface
│   └── usecase/         # Application use cases
├── cmd/                  # Application entry points
├── docker/              # Docker configurations
└── test/                # Test utilities and helpers
```

## Key Components

### Domain Layer
- **Blog**: Blog-related domain models and logic
- **Task**: Task management functionality
- **Slack**: Slack integration capabilities
- **Pokemon**: Pokemon-related notification features

### Infrastructure Layer
- Database access implementations
- DTO (Data Transfer Objects) for external services
- DAO (Data Access Objects) for data persistence

### Presentation Layer
- REST API endpoints
- GraphQL API interface
- CLI commands for various operations
- Middleware components (CORS, Prometheus, etc.)

### Use Cases
Implements application-specific business rules and coordinates the flow of data between the domain and presentation layers.

## Development Guidelines

### Architecture
- Follows Clean Architecture principles
- Implements DDD tactical patterns
- Clear separation of concerns between layers

### Code Organization
- Domain logic should be kept in the domain layer
- Business rules should be implemented in use cases
- Infrastructure concerns should be isolated in the infra layer
- Presentation layer should only handle I/O and data transformation

### Current Focus Areas
- Implementing DDD tactical patterns
- Fixing anemic domain models, particularly in:
  - Slack integration
  - Task management

### Testing
- Test files should be placed alongside the code they test
- Integration tests are available in the test directory
- Use test utilities provided in the test package

### API Interfaces
- REST API with proper middleware support
- GraphQL API for flexible data querying
- CLI interface for administrative tasks

## Technology Stack
- Go 1.25+
- Kubernetes for deployment
- Docker for containerization
- GraphQL for flexible API queries
- Prometheus for monitoring
