package routes

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/abik1221/bdu-dormitory-backend/models"

	"github.com/gin-gonic/gin"
)

type StudentInput struct {
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
	Gender       string `json:"gender" binding:"required"`
	RoomID       int    `json:"room_id" binding:"required"`
	DepartmentID int    `json:"department_id" binding:"required"`
	BuildingID   int    `json:"building_id" binding:"required"`
}

type BuildingInput struct {
	Name       string `json:"building_name" binding:"required"`
	GenderType string `json:"gender_type" binding:"required"`
}

type FloorInput struct {
	Number     int `json:"floor_number" binding:"required"`
	BuildingID int `json:"building_id" binding:"required"`
}

type RoomInput struct {
	Number    int    `json:"room_number" binding:"required"`
	Capacity  int    `json:"capacity" binding:"required"`
	Amenities string `json:"amenities"`
	FloorID   int    `json:"floor_id" binding:"required"`
}

type DepartmentInput struct {
	Name string `json:"department_name" binding:"required"`
}

type RoomAmenitiesInput struct {
	Amenities string `json:"amenities" binding:"required"`
}

func SetupRoutes(r *gin.Engine, db *sql.DB) {

	r.GET("/buildings", func(c *gin.Context) {
		buildings, err := models.GetAllBuildings(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch buildings: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, buildings)
	})

	r.POST("/buildings", func(c *gin.Context) {
		var input BuildingInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}
		err := models.CreateBuilding(db, input.Name, input.GenderType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create building: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Building created successfully"})
	})

	r.GET("/floors", func(c *gin.Context) {
		floors, err := models.GetAllFloors(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch floors: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, floors)
	})

	r.POST("/floors", func(c *gin.Context) {
		var input FloorInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}
		err := models.CreateFloor(db, input.Number, input.BuildingID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create floor: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Floor created successfully"})
	})

	r.GET("/rooms", func(c *gin.Context) {
		rooms, err := models.GetAllRooms(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rooms: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, rooms)
	})

	r.POST("/rooms", func(c *gin.Context) {
		var input RoomInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}
		err := models.CreateRoom(db, input.Number, input.Capacity, input.FloorID, input.Amenities)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Room created successfully"})
	})

	r.PUT("/rooms/:id/amenities", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}
		var input RoomAmenitiesInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}
		err = models.UpdateRoomAmenities(db, id, input.Amenities)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update amenities: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Room amenities updated successfully"})
	})

	r.GET("/departments", func(c *gin.Context) {
		departments, err := models.GetAllDepartments(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch departments: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, departments)
	})

	r.POST("/departments", func(c *gin.Context) {
		var input DepartmentInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}
		err := models.CreateDepartment(db, input.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create department: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Department created successfully"})
	})

	r.GET("/students", func(c *gin.Context) {
		students, err := models.GetAllStudents(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch students: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, students)
	})

	r.POST("/students", func(c *gin.Context) {
		var input StudentInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}
		err := models.CreateStudent(db, input.FirstName, input.LastName, input.Gender, input.RoomID, input.DepartmentID, input.BuildingID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign student: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Student assigned successfully"})
	})

	r.PUT("/students/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
			return
		}
		var input StudentInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}
		err = models.UpdateStudent(db, id, input.FirstName, input.LastName, input.Gender, input.RoomID, input.DepartmentID, input.BuildingID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Student updated successfully"})
	})

	r.DELETE("/students/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
			return
		}
		err = models.DeleteStudent(db, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete student: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
	})

	r.GET("/audit_logs", func(c *gin.Context) {
		logs, err := models.GetAllAuditLogs(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch audit logs: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, logs)
	})

	r.GET("/reports/occupancy", func(c *gin.Context) {
		rows, err := db.Query("CALL generate_occupancy_report()")
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
