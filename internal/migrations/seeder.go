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
	podcasts "github.com/Al-Khaimah/khaimah-golang-backend/internal/modules/podcasts/models"
)

func ClearDatabase(db *gorm.DB) {
	fmt.Println("⚠️ Clearing existing database...")
	db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
}

func SeedDatabase(db *gorm.DB) {
	var count int64
	db.Model(&podcasts.Podcast{}).Count(&count)
	if count > 5 {
		fmt.Println("✅ Podcasts already seeded! Skipping...")
		return
	}

	fmt.Println("🌱 Seeding only podcasts table with test data...")
	seedPodcasts(db)
	fmt.Println("🎉 Podcasts seeding completed successfully!")
}

func seedPodcasts(db *gorm.DB) {
	var categoriesList []categories.Category
	if err := db.Find(&categoriesList).Error; err != nil {
		fmt.Println("❌ Failed to load categories:", err)
		return
	}

	if len(categoriesList) == 0 {
		fmt.Println("❌ No categories found. Cannot seed podcasts without categories.")
		return
	}

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
			AudioURL:              fullAudioURL,
			CoverImageURL:         fmt.Sprintf("https://source.unsplash.com/random/400x400?sig=%d", i),
			CoverImageDescription: faker.Sentence(),
			CategoryID:            categoriesList[rand.Intn(len(categoriesList))].ID,
		}
		podcastsList = append(podcastsList, podcast)
	}

	db.Create(&podcastsList)
	fmt.Println("✅ Podcasts seeded using existing categories!")
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
