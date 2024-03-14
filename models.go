// models.go
package main

import "gorm.io/gorm"

type Task struct {
   gorm.Model
   Title       string
   Description string
   DueDate     string
   Status      string
}
