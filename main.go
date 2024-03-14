package main

import (
   "github.com/gin-gonic/gin"
   "gorm.io/driver/sqlite"
   "gorm.io/gorm"
)

var db *gorm.DB

func main() {
   // Connect to the SQLite database
   var err error
   db, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
   if err != nil {
      panic("Failed to connect to the database")
   }

   // Auto Migrate the Database
   db.AutoMigrate(&Task{})

   // Set up Gogin router
   router := gin.Default()

   // Define API routes
   router.GET("/tasks", getTasks)
   router.GET("/tasks/:id", getTask)
   router.POST("/tasks", createTask)
   router.PUT("/tasks/:id", updateTask)
   router.DELETE("/tasks/:id", deleteTask)


   // Run the server
   router.Run(":8080")
}

// Task model definition
type Task struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Status      string `json:"status"`
}

// CRUD operations


func getTasks(c *gin.Context) {
	var tasks []Task
	db.Find(&tasks)
	c.JSON(200, tasks)
}

func getTask(c *gin.Context) {
	var task Task
	if err := db.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(200, task)
}

// func getTasks(c *gin.Context) {
//     var tasks []Task
//     if err := db.Find(&tasks).Error; err != nil {
//         c.JSON(500, gin.H{"error": "Internal Server Error"})
//         return
//     }
//     c.JSON(200, tasks)
// }


// CreateTask handles the creation of a new task
func createTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Handle the case where DueDate is empty (you can add more validation)
	if task.DueDate == "" {
		c.JSON(400, gin.H{"error": "DueDate cannot be empty"})
		return
	}

	// Generate a unique ID for the task (assuming you have a mechanism for this, e.g., UUID)
	// For simplicity, let's use a simple incrementing ID in this example.
	task.ID = uint(db.Find(&Task{}).RowsAffected) + 1

	// Store the task in the database
	if err := db.Create(&task).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create the task"})
		return
	}

	c.JSON(201, task)
}


// func updateTask(c *gin.Context) {
// 	var task Task
// 	if err := db.First(&task, c.Param("id")).Error; err != nil {
// 		c.JSON(404, gin.H{"error": "Task not found"})
// 		return
// 	}

// 	if err := c.ShouldBindJSON(&task); err != nil {
// 		c.JSON(400, gin.H{"error": err.Error()})
// 		return
// 	}

// 	db.Save(&task)
// 	c.JSON(200, task)
// }

func updateTask(c *gin.Context) {
    var task Task
    if err := db.First(&task, c.Param("id")).Error; err != nil {
        c.JSON(404, gin.H{"error": "Task not found"})
        return
    }

    if err := c.ShouldBindJSON(&task); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    db.Save(&task)
    c.JSON(200, task)
}

func deleteTask(c *gin.Context) {
	var task Task
	if err := db.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(404, gin.H{"error": "Task not found"})
		return
	}

	db.Delete(&task)
	c.JSON(204, gin.H{"message": "Task deleted successfully"})
}
