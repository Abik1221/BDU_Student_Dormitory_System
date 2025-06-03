package routes

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Student represents the student model for JSON binding
type Student struct {
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	Gender       string `json:"gender" binding:"required"`
	RoomID       int    `json:"room_id" binding:"required"`
	DepartmentID int    `json:"department_id" binding:"required"`
	BuildingID   int    `json:"building_id" binding:"required"`
}

// SetupRoutes configures the API endpoints
func SetupRoutes(r *gin.Engine, db *sql.DB) {
	// GET /students: Retrieve all students
	r.GET("/students", func(c *gin.Context) {
		rows, err := db.Query(`
            SELECT student_id, first_name, last_name, gender, room_id, department_id, building_id 
            FROM student
        `)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch students: " + err.Error()})
			return
		}
		defer rows.Close()

		var students []Student
		for rows.Next() {
			var s Student
			var studentID int
			if err := rows.Scan(&studentID, &s.FirstName, &s.LastName, &s.Gender, &s.RoomID, &s.DepartmentID, &s.BuildingID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan student: " + err.Error()})
				return
			}
			students = append(students, s)
		}

		c.JSON(http.StatusOK, students)
	})

	// POST /students: Assign a new student using stored procedure
	r.POST("/students", func(c *gin.Context) {
		var student Student
		if err := c.ShouldBindJSON(&student); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		// Call the stored procedure
		result, err := db.Exec(`
            CALL assign_student_to_room(?, ?, ?, ?, ?, ?)
        `, student.FirstName, student.LastName, student.Gender, student.RoomID, student.DepartmentID, student.BuildingID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign student: " + err.Error()})
			return
		}

		// Check if the procedure was successful
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No rows affected, possible database error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Student assigned successfully"})
	})

	// GET /reports/occupancy: Generate occupancy report
	r.GET("/reports/occupancy", func(c *gin.Context) {
		rows, err := db.Query(`CALL generate_occupancy_report()`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report: " + err.Error()})
			return
		}
		defer rows.Close()

		var reports []map[string]interface{}
		columns, _ := rows.Columns()
		for rows.Next() {
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))
			for i := range values {
				valuePtrs[i] = &values[i]
			}
			if err := rows.Scan(valuePtrs...); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan report: " + err.Error()})
				return
			}

			report := make(map[string]interface{})
			for i, col := range columns {
				val := values[i]
				if b, ok := val.([]byte); ok {
					report[col] = string(b)
				} else {
					report[col] = val
				}
			}
			reports = append(reports, report)
		}

		c.JSON(http.StatusOK, reports)
	})
}
