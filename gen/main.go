package main

import (
	"streambot/db"
	"streambot/models"

	"gorm.io/gen"
)

// This name follows Go conventions for interfaces but I still don't like it
// Feel free to rename
type Subscriber interface {
	// GetByHealth queries for instances that meet a health check threshold
	//
	// SELECT * from @@table WHERE times_failed >= @min AND times_failed < @max AND deleted_at IS NULL
	GetByHealth(min int, max int) ([]*gen.T, error)
}

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	// gormdb, _ := gorm.Open(mysql.Open("root:@(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"))
	g.UseDB(db.Connection) // reuse your gorm db

	// Generate basic type-safe DAO API for struct `model.User` following conventions
	g.ApplyBasic(models.All...)
	g.ApplyInterface(func(Subscriber) {}, models.Subscription{})

	// Generate the code
	g.Execute()
}
