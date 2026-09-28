# Multi-Tenant Access Control API

A production-style backend API for managing authentication, organizations, team memberships, roles, permissions, sessions, invitations, and security auditing in a multi-tenant environment.

The system is built with **Go and Gin**, uses **PostgreSQL** for persistent data, and **Redis** for RBAC permission caching and distributed rate limiting. Authentication is implemented using short-lived JWT access tokens and rotating refresh tokens backed by database sessions.

The project focuses on real-world backend engineering concepts such as tenant isolation, role-based access control, transactional workflows, session security, auditability, caching, rate limiting, API documentation, containerization, and CI/CD.

### Key Features

* JWT-based authentication
* Access and refresh token flow
* Refresh token rotation and session revocation
* Multi-tenant organizations
* Organization memberships
* Role-Based Access Control (RBAC)
* Permission-based authorization
* Team invitation and acceptance workflow
* Audit logging for sensitive actions
* Redis-based RBAC permission caching
* Redis-based API rate limiting
* PostgreSQL transactions for security-sensitive workflows
* Paginated audit log API
* Swagger/OpenAPI documentation
* Docker containerization
* GitHub Actions CI
* Deployment on Render

### Tech Stack

**Backend:** Go, Gin
**Database:** PostgreSQL
**Cache & Rate Limiting:** Redis / Valkey
**Authentication:** JWT, bcrypt
**API Documentation:** Swagger / OpenAPI
**Containerization:** Docker
**CI/CD:** GitHub Actions
**Deployment:** Render
