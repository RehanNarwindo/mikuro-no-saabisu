cat > README.md << 'EOF'
# mikuro-no-saabisu

A microservices project built with Go.

## Services

- **auth-service** - Handles authentication and authorization
- **user-service** - Manages user data and profiles

## Prerequisites

- [Go](https://golang.org/dl/) 1.21+
- [Docker](https://www.docker.com/) (optional)

## Getting Started

### Clone the repository
```bash
git clone https://github.com/RehanNarwindo/mikuro-no-saabisu.git
cd mikuro-no-saabisu
```

### Run a service
```bash
cd auth-service
go run main.go
```
