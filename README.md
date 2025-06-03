# BDU Student Dormitory Management System Backend

## Overview
The **BDU Student Dormitory Management System** is a RESTful API backend developed using Go and the Gin framework to manage student dormitory assignments at Bahir Dar University (BDU). It provides robust functionality for handling buildings, floors, rooms, departments, students, and audit logs, with a focus on database-driven operations. The system leverages a MySQL database with stored procedures, triggers, and role-based access control to enforce business rules such as gender-based building assignments, room capacity limits, and audit tracking.

This project fulfills the requirements of a university assignment to create a dormitory management system, emphasizing:
- **Database Integrity**: CRUD operations with MySQL, enforced by triggers and stored procedures.
- **Business Logic**: Ensures gender compatibility, room capacity constraints, and audit logging.
- **RESTful API**: Comprehensive endpoints for managing entities and generating reports.
- **Modularity**: Clean, maintainable Go code with separation of concerns.

## System Workflow
The system processes HTTP requests through a structured workflow, interacting with the MySQL database to manage dormitory data. Below is the workflow for key operations:

1. **Request Reception**:
   - HTTP requests (e.g., `GET /students`, `POST /rooms`) are received by the Gin router defined in `cmd/main.go`.
   - The router delegates requests to handlers in `routes/routes.go`.

2. **Input Validation**:
   - Request payloads are validated using Gin's binding (e.g., `binding:"required"` for mandatory fields).
   - Invalid inputs return a `400 Bad Request` response with error details.

3. **Database Connection**:
   - `config/db.go` establishes a MySQL connection using credentials from the `.env` file.
   - The `*sql.DB` connection is passed to `routes.go` and `models.go` for database operations.

4. **Model Operations**:
   - Handlers in `routes/routes.go` call functions in `models/models.go` (e.g., `GetAllStudents`, `CreateBuilding`).
   - Model functions execute SQL queries or stored procedures, interacting with the database.

5. **Database Logic**:
   - **Stored Procedures**: Operations like `assign_student_to_room` and `update_room_amenities` enforce complex logic (e.g., gender checks, capacity limits).
   - **Triggers**: The `restrict_gender` and `restrict_room_capacity` triggers validate student assignments, while `log_student_changes` logs modifications to the `audit_log` table.
   - **Roles**: Database users (`admin`, `manager`, `viewer`) restrict access based on permissions.

6. **Response Generation**:
   - Successful operations return JSON responses with data or confirmation messages (e.g., `200 OK`).
   - Errors (e.g., database failures, constraint violations) return `500 Internal Server Error` or `400 Bad Request` with descriptive messages.

**Example Workflow (POST /students)**:
- A client sends a POST request with student details (e.g., `{"first_name":"Abebe","last_name":"Kebede","gender":"male","room_id":3,"department_id":1,"building_id":1}`).
- The handler validates the JSON payload.
- `models.CreateStudent` calls the `assign_student_to_room` stored procedure.
- Triggers check gender compatibility and room capacity.
- If valid, the student is inserted, and an audit log entry is created.
- The client receives `{"message":"Student assigned successfully"}`.

## Database Schema
The MySQL database (`bdu_student_dormitory_management_system`) is defined in `sql/dormitory_management.sql` and includes:

### Tables
- **building**:
  - Columns: `building_id` (PK, int), `building_name` (varchar), `gender_type` (char: 'male' or 'female')
  - Purpose: Stores dormitory buildings with gender restrictions.
- **floor**:
  - Columns: `floor_id` (PK, int), `floor_number` (int), `building_id` (FK, int)
  - Purpose: Represents floors within buildings.
- **room**:
  - Columns: `room_id` (PK, int), `room_number` (int), `capacity` (int), `amenities` (varchar, nullable), `floor_id` (FK, int)
  - Purpose: Stores room details, including capacity and amenities.
- **department**:
  - Columns: `department_id` (PK, int), `department_name` (varchar)
  - Purpose: Stores academic departments.
- **student**:
  - Columns: `student_id` (PK, int), `first_name` (varchar), `last_name` (varchar), `gender` (char), `room_id` (FK, int), `department_id` (FK, int), `building_id` (FK, int)
  - Purpose: Stores student details and their assignments.
- **audit_log**:
  - Columns: `log_id` (PK, int), `table_name` (varchar), `operation` (varchar), `record_id` (int), `change_time` (datetime), `changed_by` (varchar), `details` (varchar)
  - Purpose: Logs changes to the `student` table for auditing.

### Stored Procedures
- **assign_student_to_room(first_name, last_name, gender, room_id, department_id, building_id)**:
  - Inserts a student, ensuring gender and capacity constraints are met (via triggers).
  - Used by `POST /students`.
- **update_room_amenities(room_id, amenities)**:
  - Updates a room’s amenities field.
  - Used by `PUT /rooms/:id/amenities`.
- **generate_occupancy_report()**:
  - Returns a report of room occupancy (building name, floor, room number, capacity, occupants, available beds).
  - Used by `GET /reports/occupancy`.

### Triggers
- **restrict_gender**:
  - Ensures a student’s gender matches the building’s `gender_type` before insertion.
  - Raises an error for mismatches (e.g., female in a male building).
- **restrict_room_capacity**:
  - Checks if a room’s capacity is exceeded before assigning a student.
  - Raises an error if the room is full.
- **log_student_changes**:
  - Logs insert, update, and delete operations on the `student` table to `audit_log`.
  - Captures operation type, record ID, timestamp, user, and details.

### Roles and Users
- **Roles**:
  - `admin`: Full permissions (all operations).
  - `manager`: Read and write permissions (no schema modifications).
  - `viewer`: Read-only permissions.
- **Users**:
  - `admin_user` (password: `admin_pass`, role: `admin`): Used by the application for full access.
  - `manager_user` (password: `manager_pass`, role: `manager`): For restricted write access.
  - `viewer_user` (password: `viewer_pass`, role: `viewer`): For read-only access.

## Requirements
### Software
- **Go**: Version 1.21 or later (download from [golang.org](https://golang.org/)).
- **MySQL**: Version 8.0 or later (download from [mysql.com](https://www.mysql.com/)).
- **VS Code**: Recommended IDE with the **Go** extension for development and debugging.
- **REST Client Extension**: Optional for testing API endpoints in VS Code.
- **MySQL Client**: For executing SQL scripts (e.g., MySQL Workbench or command-line client).

### Dependencies
- Go modules (defined in `go.mod`):
  - `github.com/gin-gonic/gin`: Web framework for routing.
  - `github.com/go-sql-driver/mysql`: MySQL driver.
  - `github.com/joho/godotenv`: Environment variable management.
- Install dependencies using `go mod tidy` (see setup instructions).

### Environment
- A `.env` file in the project root with MySQL credentials (see setup instructions).
- MySQL server running locally on port `3306` (default).

## Setup Instructions
Follow these steps to set up and run the backend locally:

1. **Clone the Repository**:
   - Clone or copy the project to your local machine:
     ```bash
     git clone github.com/abik1221/BDU_Student_Dormitory_System
     cd BDU_Student_Dormitory_System
     ```
2. **Create the .env file and set all the approprita data by mapping with the config file.**
```bash 
 touch .env
 ```

3.**Install all the 3rd part dependencies using the simple golang command.**
```bash
go mod tidy
```
4.Run the backend program using 
```bash
go run cmd/main.go
```

pls use postman or tender client to test the Api end points properly


## contributers 

 Meklit Anteneh (1507008), 
 Tamer Getaw (1507859), 
 Feruza Mohamed (1506085), 
 Naom Keneni (1507336), 
 Eyerus Tekto (1506000)

 Thankyou!

 Pls fork the repository and contribute to inheance our campus dormitory managemny ststem it is an open source project!

