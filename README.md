# 🔐 Multi-Tenant Access Control API

A production-style **multi-tenant authentication and authorization backend** built with Go, Gin, PostgreSQL, and Redis.

The system provides secure authentication, organization-level access control, RBAC permissions, team invitations, session management, audit logging, Redis caching, and API rate limiting.

## 🚀 Live Demo

* **Live API:** https://multi-tenant-access-control.onrender.com/health?utm_source=chatgpt.com
* **Health Check:** https://multi-tenant-access-control.onrender.com/health
* **Swagger API Docs:** https://multi-tenant-access-control.onrender.com/swagger/index.html
* **GitHub:** https://github.com/Ritik-x/Multi-Tenant-Access-Control

---

## 🛠️ Tech Stack

| Technology        | Purpose                      |
| ----------------- | ---------------------------- |
| Go                | Backend programming language |
| Gin               | HTTP web framework           |
| PostgreSQL        | Primary relational database  |
| Redis / Valkey    | RBAC caching & rate limiting |
| JWT               | Access authentication        |
| bcrypt            | Password hashing             |
| Docker            | Containerization             |
| GitHub Actions    | CI                           |
| Swagger / OpenAPI | API documentation            |
| Render            | Deployment                   |

---

## ✨ Features

### 🔑 Authentication

* User registration
* User login
* JWT access tokens
* Refresh tokens
* Refresh token rotation
* Token expiration
* Secure password hashing with bcrypt

### 🏢 Multi-Tenancy

* Organization creation
* Organization-based data isolation
* User memberships
* Organization-specific roles
* Tenant boundary validation

### 🛡️ RBAC & Permissions

* Role-Based Access Control
* Permission-based authorization
* Organization-specific roles
* Permission middleware
* PostgreSQL-backed permission management
* Redis-based permission caching

### 👥 Team Invitations

* Invite users to an organization
* Assign a role during invitation
* Role belongs-to-organization validation
* Secure invitation token hashing
* Invitation expiration
* Invitation acceptance workflow

### 🔐 Session Management

* Database-backed refresh sessions
* Refresh token rotation
* Session revocation
* Revoke all sessions
* Session expiration handling
* Row-level locking during token rotation

### 📋 Audit Logging

Security-sensitive actions are recorded with:

* User ID
* Organization ID
* Action
* Resource
* Resource ID
* Metadata
* IP address
* Timestamp

Audit logs are organization-scoped and paginated.

### ⚡ Redis

Redis is used for:

* RBAC permission caching
* Distributed API rate limiting

The RBAC authorization flow follows a cache-aside approach:

```text
Request
   ↓
Redis Cache
   ↓
Cache Hit ──────→ Check Permission
   │
   └─ Cache Miss
          ↓
     PostgreSQL
          ↓
     Store in Redis
          ↓
     Check Permission
```

### 🚦 Rate Limiting

Redis-based rate limiting is applied to sensitive authentication endpoints:

* Login
* Registration
* Refresh token

Rate-limit keys are scoped by client IP and endpoint.

### 📖 API Documentation

The REST API is documented using Swagger/OpenAPI.

**Swagger UI:** https://multi-tenant-access-control.onrender.com/swagger/index.html

---

## 🏗️ Architecture

The backend follows a layered architecture:

```text
                    ┌──────────────┐
                    │    Client    │
                    └──────┬───────┘
                           │
                           ▼
                    ┌──────────────┐
                    │     Gin      │
                    │   Handlers   │
                    └──────┬───────┘
                           │
                           ▼
                    ┌──────────────┐
                    │   Services   │
                    │ Business     │
                    │   Logic      │
                    └──────┬───────┘
                           │
                 ┌─────────┴─────────┐
                 ▼                   ▼
          ┌──────────────┐    ┌──────────────┐
          │ Repositories │    │    Redis     │
          │  PostgreSQL  │    │ Cache/Rate   │
          └──────────────┘    │   Limiting   │
                              └──────────────┘
```

### Main layers

**Handlers**

* Receive HTTP requests
* Validate request input
* Return HTTP responses

**Services**

* Contain business logic
* Coordinate repositories
* Manage transactions and workflows

**Repositories**

* Handle database operations
* Keep SQL/database logic separated from business logic

**Middleware**

* Authentication
* Authorization
* Rate limiting

---

## 🔐 Authentication Flow

```text
Register / Login
       ↓
 Validate Credentials
       ↓
 Generate Access Token
       ↓
 Generate Refresh Token
       ↓
 Store Refresh Session
       ↓
 Return Tokens
```

Access tokens are short-lived, while refresh tokens are stored as hashes in PostgreSQL.

Refresh token rotation follows:

```text
Old Refresh Token
       ↓
Find Session
       ↓
FOR UPDATE
       ↓
Validate Session
       ↓
Revoke Old Session
       ↓
Generate New Refresh Token
       ↓
Create New Session
       ↓
Commit Transaction
```

This prevents concurrent refresh requests from safely reusing the same refresh session.

---

## 🛡️ Authorization Flow

```text
HTTP Request
     ↓
JWT Middleware
     ↓
Extract User + Organization
     ↓
Permission Middleware
     ↓
Redis Permission Cache
     ↓
   Cache Hit?
    /     \
  Yes      No
   ↓        ↓
Check    PostgreSQL
Permission   ↓
   │      Cache Result
   └───────┘
      ↓
 Allow / Deny
```

Permissions are evaluated within the user's organization to maintain tenant isolation.

---

## 📋 API Endpoints

### Authentication

```text
POST   /register
POST   /login
POST   /refresh
POST   /logout
```

### Sessions

```text
GET    /sessions
DELETE /sessions/:id
POST   /sessions/revoke-all
```

### Invitations

```text
POST   /invitations
POST   /invitations/accept
```

### Audit Logs

```text
GET    /audit-logs
```

### Health

```text
GET    /health
```

For request/response schemas and authorization requirements, see the Swagger documentation.

---

## 🐳 Docker

The application is containerized using a multi-stage Docker build.

The project can run with:

```text
Go API
   │
   ├── PostgreSQL
   │
   └── Redis
```

Docker Compose is used for local development with PostgreSQL and Redis services.

---

## 🔄 CI/CD

GitHub Actions is configured to run on pushes and pull requests to the `main` branch.

The CI pipeline:

```text
Git Push / Pull Request
          ↓
     Checkout Code
          ↓
      Setup Go
          ↓
   Download Dependencies
          ↓
       Run Tests
          ↓
     Build Application
```

The project also builds and publishes its Docker image through GitHub Actions.

---

## 📦 Project Structure

```text
.
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── config/
│   ├── database/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── redis/
│   ├── repository/
│   └── services/
│
├── migrations/
│
├── docs/
│
├── .github/
│   └── workflows/
│
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

---

## ⚙️ Local Setup

### 1. Clone the repository

```bash
git clone https://github.com/Ritik-x/Multi-Tenant-Access-Control.git

cd Multi-Tenant-Access-Control
```

### 2. Configure environment variables

Create a `.env` file:

```env
APP_ENV=development
PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=team_access_control
DB_SSLMODE=disable

JWT_SECRET=your_jwt_secret

REDIS_ADDR=localhost:6379
```

### 3. Start PostgreSQL and Redis

Make sure PostgreSQL and Redis are running locally.

Alternatively, use Docker Compose for the supporting services.

### 4. Run database migrations

Apply the migrations using your migration tool.

### 5. Start the API

```bash
go run ./cmd/server
```

The API will be available at:

```text
http://localhost:8080
```

Swagger:

```text
http://localhost:8080/swagger/index.html
```

---

## 🌐 Deployment

The application is deployed as a Docker-based web service on Render.

Deployment architecture:

```text
GitHub
   ↓
GitHub Actions
   ↓
Docker Image
   ↓
Render
   ↓
Go API
   ├── PostgreSQL
   └── Redis / Valkey
```

### Production API

[https://multi-tenant-access-control.onrender.com](https://multi-tenant-access-control.onrender.com/health?utm_source=chatgpt.com)

### Swagger

https://multi-tenant-access-control.onrender.com/swagger/index.html

---

## 🎯 Engineering Concepts Demonstrated

This project was designed to demonstrate practical backend engineering concepts including:

* REST API design
* Layered architecture
* Authentication & authorization
* Multi-tenancy
* RBAC
* Database transactions
* PostgreSQL row-level locking
* Refresh token rotation
* Session management
* Secure password hashing
* Redis caching
* Distributed rate limiting
* Auditability
* API pagination
* Docker
* CI/CD
* API documentation
* Cloud deployment

---

## 👨‍💻 Author

**Ritik Garg**



GitHub: https://github.com/Ritik-x
