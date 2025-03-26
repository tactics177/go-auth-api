# Examen : Développement d'une API d'Authentification en Go

## Group:

- Gregoire Lequippe
- Ali Loudagh
- Nahel Kini
- Roman Sabechkine

# Go Auth API

A simple and secure REST API for user authentication, built using **Golang**, **Gin**, and **MongoDB**.  
The project follows a clean layered architecture and includes essential authentication features like registration,
login, password reset, token-based authorization, and logout with JWT blacklisting.

---

## Live Deployment

**Base URL**: [`https://go-auth-api-e7bo.onrender.com/`](https://go-auth-api-e7bo.onrender.com/)

All endpoints are prefixed with our group code 73f2fc18-3053-4c38-943a-416d16432450

---

## API Endpoints

| Method | Endpoint           | Description                                    |
|--------|--------------------|------------------------------------------------|
| GET    | `/health`          | Health check                                   |
| POST   | `/register`        | Register a new user                            |
| POST   | `/login`           | Login                                          |
| POST   | `/forgot-password` | Send password reset token                      |
| POST   | `/reset-password`  | Reset password using token                     |
| GET    | `/me`              | Get current user profile (protected route)     |
| POST   | `/logout`          | Logout and blacklist token   (protected route) |

> Note: All routes are prefixed with our group ID.

---

---

## Run Locally with Docker

```bash
docker run `
  --name go-auth-api `
  -p 8080:8080 `
  -e MONGO_URI="replace_with_mongodb_uri" `
  -e JWT_SECRET="replace_with_jwt_secret" `
  -e DB_NAME="replace_with_db_name" `
  efrei/go-auth-api
