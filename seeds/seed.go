package seed

import (
	"app/config"
	"app/migrations"
	"app/models"
	"math/rand"
	"strconv"

	"github.com/bxcodec/faker/v3"
	"github.com/labstack/gommon/log"
)

func Load() {
	db := config.GetDB()

	// Clean Database
	db.DropTableIfExists("users", "articles", "categories", "migrations")
	migrations.Migrate()

	// Add Admin
	log.Info("Creating admin...")

	admin := models.User{
		Email:    "admin@sing3demons.com",
		Password: "passw0rd",
		Name:     "Admin",
		Role:     "Admin",
		Avatar:   "https://i.pravatar.cc/100",
	}

	admin.Password = admin.GenerateEncryptedPassword()
	db.Create(&admin)

	// Add normal users
	log.Info("Creating users...")

	numOfUsers := 50
	users := make([]models.User, 0, numOfUsers)
	userRoles := [2]string{"Editor", "Member"}

	for i := 1; i <= numOfUsers; i++ {
		user := models.User{
			Name:     faker.Name(),
			Email:    faker.Email(),
			Password: "passw0rd",
			Avatar:   "https://i.pravatar.cc/100?" + strconv.Itoa(i),
			Role:     userRoles[rand.Intn(2)],
		}

		user.Password = user.GenerateEncryptedPassword()
		db.Create(&user)
		users = append(users, user)
	}

	// Add categories
	log.Info("Creating categories...")

	numOfCategories := 20
	categories := make([]models.Category, 0, numOfCategories)

	for i := 1; i <= numOfCategories; i++ {
		category := models.Category{
			Name: faker.Word(),
		}

		db.Create(&category)
		categories = append(categories, category)
	}

	// Add articles
	// log.Info("Creating articles...")

	// numOfArticles := 50
	// articles := make([]models.Product, 0, numOfArticles)

	// for i := 1; i <= numOfArticles; i++ {
	// 	article := models.Product{
	// 		Title:      faker.Sentence(),
	// 		Excerpt:    faker.Sentence(),
	// 		Body:       faker.Paragraph(),
	// 		Image:      "https://source.unsplash.com/random/300x200?" + strconv.Itoa(i),
	// 		CategoryID: uint(rand.Intn(numOfCategories) + 1),
	// 		UserID:     uint(rand.Intn(numOfUsers) + 1),
	// 	}

	// 	db.Create(&article)
	// 	articles = append(articles, article)
	// }
}
