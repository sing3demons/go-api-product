package seeds

import (
	"math/rand"
	"strconv"

	"github.com/sing3demons/app/database"

	"github.com/sing3demons/app/models"

	"github.com/bxcodec/faker/v3"
	"github.com/labstack/gommon/log"
)

func Load() {
	db := database.GetDB()

	var productsDB []models.Product
	err := db.Find(&productsDB).Error

	if len(productsDB) == 0 || err != nil {
		// Clean Database
		// db.DropTableIfExists("users", "products", "categories", "migrations")

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

		numOfCategories := 10
		categories := make([]models.Category, 0, numOfCategories)

		for i := 1; i <= numOfCategories; i++ {
			category := models.Category{
				Name: faker.Word(),
			}

			db.Create(&category)
			categories = append(categories, category)
		}

		// Add products
		log.Info("Creating products...")

		numOfProducts := 100000
		products := make([]models.Product, 0, numOfProducts)

		for i := 1; i <= numOfProducts; i++ {
			product := models.Product{
				Name:       faker.Name(),
				Desc:       faker.Word(),
				Price:      rand.Intn(9999),
				Image:      "https://source.unsplash.com/random/300x200?" + strconv.Itoa(i),
				CategoryID: uint(rand.Intn(numOfCategories) + 1),
			}

			db.Create(&product)
			products = append(products, product)
		}

	}

}
