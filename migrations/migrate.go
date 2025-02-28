package migrations

import (
	"fmt"
	"log"

	"github.com/Al-Khaimah/khaimah-golang-backend/config"
	users "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/models"
)

func Migrate() {
	db := config.GetDB()
	fmt.Println("Running Migrations...")
	err := db.AutoMigrate(
		&users.User{},
		&users.Auth{},
		// &podcasts.Podcast{},
		// &notifications.Notification{},
		// &categories.Category{},
		// &bookmarks.Bookmarks{},
	//call models
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	fmt.Println("Migrations completed successfully!")
}
