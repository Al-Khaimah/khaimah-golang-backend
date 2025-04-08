package migrations

import (
	"fmt"
	categories "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/models"
	notifications "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/notifications/models"
	podcasts "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/models"
	"log"

	"github.com/Al-Khaimah/khaimah-golang-backend/config"
	users "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/models"
)

func Migrate() {
	db := config.GetDB()
	fmt.Println("Running Migrations...")
	err := db.AutoMigrate(
		&users.User{},
		&users.IamAuth{},
		&categories.Category{},
		&notifications.Notification{},
		&podcasts.Podcast{},
		&podcasts.UserPodcast{},
		&podcasts.BookmarkPodcast{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	fmt.Println("Migrations completed successfully!")
}
