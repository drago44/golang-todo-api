// @title Golang Todo API
// @version 1.0
// @description Simple Todo API built with Gin and GORM.
// @BasePath /api/v1
// @schemes http https
package main

import "github.com/drago44/golang-todo-api/internal/app"

func main() {
	app.Run()
}
