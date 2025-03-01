package migrations

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/bxcodec/faker/v3"
	"gorm.io/gorm"

	categories "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/categories/models"
	notifications "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/notifications/models"
	podcasts "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/models"
	users "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/users/models"
)

func ClearDatabase(db *gorm.DB) {
	fmt.Println("⚠️ Clearing existing database...")
	db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
}

func SeedDatabase(db *gorm.DB) {
	var count int64
	db.Model(&podcasts.Podcast{}).Count(&count)
	if count > 5 {
		fmt.Println("✅ Database already seeded! Skipping...")
		return
	}

	fmt.Println("🌱 Seeding database with test data...")

	categoriesList := seedCategories(db)
	usersList := seedUsers(db)
	podcastsList := seedPodcasts(db, categoriesList)
	seedUserPodcasts(db, usersList, podcastsList, categoriesList)
	seedUserBookmarks(db, usersList, podcastsList)
	seedUserCategories(db, usersList, categoriesList)
	seedAuth(db, usersList)
	seedNotifications(db, usersList)

	fmt.Println("🎉 Seeding completed successfully!")
}

func seedCategories(db *gorm.DB) []categories.Category {
	categoryNames := []string{"Technology", "Sports", "Finance", "Health", "Music", "Science", "Entertainment", "Education"}
	var categoriesList []categories.Category
	for _, name := range categoryNames {
		category := categories.Category{
			Name:        name,
			Description: faker.Sentence(),
		}
		categoriesList = append(categoriesList, category)
	}
	db.Create(&categoriesList)
	fmt.Println("✅ Categories seeded!")
	return categoriesList
}

func seedUsers(db *gorm.DB) []users.User {
	var usersList []users.User
	for i := 0; i < 50; i++ {
		user := users.User{
			FirstName: faker.FirstName(),
			LastName:  faker.LastName(),
			Email:     faker.Email(),
		}
		usersList = append(usersList, user)
	}
	db.Create(&usersList)
	fmt.Println("✅ Users seeded!")
	return usersList
}

func seedPodcasts(db *gorm.DB, categoriesList []categories.Category) []podcasts.Podcast {
	var podcastsList []podcasts.Podcast

	audioDir := "audio"
	if _, err := os.Stat(audioDir); os.IsNotExist(err) {
		os.Mkdir(audioDir, os.ModePerm)
	}

	for i := 0; i < 200; i++ {
		filename := fmt.Sprintf("audio/podcast_%d.wav", i+1)
		generateRandomWAV(filename)
		fullAudioURL := fmt.Sprintf("http://localhost:8080/%s", filename)

		podcast := podcasts.Podcast{
			Title:                 faker.Word() + " Podcast",
			Description:           faker.Paragraph(),
			AudioURL:              fullAudioURL, // Store full URL
			CoverImageURL:         fmt.Sprintf("https://source.unsplash.com/random/400x400?sig=%d", i),
			CoverImageDescription: faker.Sentence(),
			CategoryID:            categoriesList[rand.Intn(len(categoriesList))].ID,
		}
		podcastsList = append(podcastsList, podcast)
	}
	db.Create(&podcastsList)
	fmt.Println("✅ Podcasts seeded with WAV files!")
	return podcastsList
}

func seedUserPodcasts(db *gorm.DB, usersList []users.User, podcastsList []podcasts.Podcast, categoriesList []categories.Category) {
	var userPodcasts []podcasts.UserPodcast
	for i := 0; i < 500; i++ {
		userPodcast := podcasts.UserPodcast{
			UserID:         usersList[rand.Intn(len(usersList))].ID,
			PodcastID:      podcastsList[rand.Intn(len(podcastsList))].ID,
			CategoryID:     categoriesList[rand.Intn(len(categoriesList))].ID,
			ResumePosition: rand.Intn(1000),
			IsCompleted:    rand.Intn(2) == 1,
		}
		userPodcasts = append(userPodcasts, userPodcast)
	}
	db.Create(&userPodcasts)
	fmt.Println("✅ User-Podcast relationships seeded!")
}

func seedUserBookmarks(db *gorm.DB, usersList []users.User, podcastsList []podcasts.Podcast) {
	for i := 0; i < 300; i++ {
		user := &usersList[rand.Intn(len(usersList))]
		podcast := &podcastsList[rand.Intn(len(podcastsList))]

		db.Model(user).Association("Bookmarks").Append(podcast)
	}
	fmt.Println("✅ User-Bookmarks seeded!")
}

func seedUserCategories(db *gorm.DB, usersList []users.User, categoriesList []categories.Category) {
	for i := 0; i < 200; i++ {
		user := &usersList[rand.Intn(len(usersList))]
		category := &categoriesList[rand.Intn(len(categoriesList))]

		db.Model(user).Association("Categories").Append(category)
	}
	fmt.Println("✅ User-Categories seeded!")
}

func seedAuth(db *gorm.DB, usersList []users.User) {
	var authRecords []users.IamAuth
	for _, user := range usersList {
		auth := users.IamAuth{
			UserID:   user.ID,
			Password: faker.Password(),
			IsActive: rand.Intn(2) == 1,
		}
		authRecords = append(authRecords, auth)
	}
	db.Create(&authRecords)
	fmt.Println("✅ IAM Auth records seeded!")
}

func seedNotifications(db *gorm.DB, usersList []users.User) {
	var notificationList []notifications.Notification
	for i := 0; i < 150; i++ {
		notification := notifications.Notification{
			UserID:      usersList[rand.Intn(len(usersList))].ID,
			Title:       faker.Word(),
			Description: faker.Sentence(),
			IsRead:      rand.Intn(2) == 1,
			Type:        "general",
		}
		notificationList = append(notificationList, notification)
	}
	db.Create(&notificationList)
	fmt.Println("✅ Notifications seeded!")
}

func generateRandomWAV(filename string) {
	duration := time.Duration(rand.Intn(14)+1) * time.Minute

	file, err := os.Create(filename)
	if err != nil {
		log.Printf("❌ Failed to create WAV file: %s\n", err)
		return
	}
	defer file.Close()

	data := make([]byte, 44100*2*int(duration.Seconds()))
	_, err = file.Write(data)
	if err != nil {
		log.Printf("❌ Failed to write WAV data: %s\n", err)
	}
}
