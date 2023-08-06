package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"test-pelindo/handlers"
	"test-pelindo/models"
)

func ApiRegistration(db *gorm.DB) {
	migrator := db.Migrator()
	if migrator.HasTable("users") {
		fmt.Println("The users table exists in the database.")
	} else {
		fmt.Println("The users table does not exist in the database.")
		db.AutoMigrate(&models.User{})
	}

	r := mux.NewRouter()
	userService := handlers.UserServices(db)

	r.HandleFunc("/user", userService.UserHandler)
	r.HandleFunc("/user/{id}", userService.GetUserDetail)

	http.Handle("/", r)

	http.ListenAndServe(":8090", nil)
}
