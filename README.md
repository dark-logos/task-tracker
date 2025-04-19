Task Tracker 
Task Tracker is a RESTful API for task management, built with Go using the Gin framework. The project allows users to register, log in, manage tasks (create, view, update, delete), and provides Prometheus metrics support.

logging + task creating (PLEASE, WAIT 2-3SEC FOR GIFS LOADING!)
![demo](docs/gifs/demo_task.gif)

update tasks + upload the page
![demo](docs/gifs/demo_update.gif)



________________________________________
Table of Contents:

Requirements
Installation
Configuration
Running
API Endpoints
Project Structure
Contributing

Requirements

Go 1.24 or higher
PostgreSQL (for data storage)
Git (for cloning the repository)
(Optional) Docker for containerization

Installation

Clone the repository:
git clone https://github.com/yourusername/task-tracker.git
cd task-tracker


Install dependencies:Ensure Go is installed, then run:
go mod tidy

Set up PostgreSQL:

Install PostgreSQL if not already installed.
Create a database:CREATE DATABASE task_tracker;

Configuration

Create a .env file in the root directory:
touch .env

Add the following environment variables:
DATABASE_URL=postgres://user:password@localhost:5432/task_tracker?sslmode=disable
JWT_SECRET= *jwt_token*
PORT=8080


Replace user, password, and other values with your own.

Running

Run the server:
go run cmd/api/main.go
Run the front:
node frontend/server.js 

The server will start on the port specified in .env (default 8080).

________________________________________________

API Endpoints
Authentication

POST /register — Register a new user.
Request body: {"username": "user", "password": "pass"}
Response: 200 OK or 400 Bad Request

POST /login — Log in.
Request body: {"username": "user", "password": "pass"}
Response: 200 OK (with JWT token) or 401 Unauthorized

POST /refresh — Refresh token.
Header: Authorization: Bearer <refresh_token>
Response: 200 OK (with new token)

Tasks (requires authentication)

GET /tasks — Get list of tasks.
Response: 200 OK with task array

POST /tasks — Create a task.
Request body: {"title": "Task", "description": "Desc", "status": "pending", "priority": 1}
Response: 200 OK

GET /tasks/:id — Get task by ID.
Response: 200 OK or 404 Not Found

PUT /tasks/:id — Update a task.
Request body: {"title": "Updated", "description": "Updated Desc", "status": "in_progress", "priority": 2}
Response: 200 OK or 404 Not Found

DELETE /tasks/:id — Delete a task.
Response: 200 OK or 404 Not Found

______________________________________________

Contributing

Fork the repository.
Create a new branch: git checkout -b feature/your-feature.
Make changes and commit: git commit -m "Add your feature".
Push changes: git push origin feature/your-feature.
Create a Pull Request.

