package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/jameslahm/bloggy_backend/api/controllers"
	"github.com/jameslahm/bloggy_backend/models"
	"github.com/jinzhu/gorm"
)

func Database(server *controllers.Server) {
	var err error

	TestDbDriver := os.Getenv("TEST_DB_DRIVER")
	if TestDbDriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_PORT"), os.Getenv("TEST_DB_NAME"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("Error: ", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
	if TestDbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_PORT"), os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_NAME"), os.Getenv("TEST_DB_PASSWORD"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("Error: ", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
}

func RefreshUserTable(server *controllers.Server) error {
	err := server.DB.DropTableIfExists(&models.Post{}).Error
	err = server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func SeedOneUser(server *controllers.Server) (models.User, error) {
	user := models.User{
		Nickname: "Pet",
		Email:    "fake@example.com",
		Password: "password",
	}
	hashedPassword, err := models.Hash(user.Password)
	if err != nil {
		log.Fatalf("Error Hash password %v\n", err)
	}
	user.Password = string(hashedPassword)
	err = server.DB.Debug().Model(&models.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("Cannot seed users table: %v\n", err)
	}
	return user, nil
}

func SeedUsers(server *controllers.Server) error {
	users := []models.User{
		{
			Nickname: "Steven",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		{
			Nickname: "Kenny",
			Email:    "kenny",
			Password: "password",
		},
	}
	for i := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func RefreshUserAndPostTable(server *controllers.Server) error {
	err := server.DB.DropTableIfExists(&models.User{}, &models.Post{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.Post{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")

	return nil
}

func SeedOneUserAndOnePost(server *controllers.Server) (models.Post, error) {

	err := RefreshUserAndPostTable(server)
	if err != nil {
		return models.Post{}, err
	}
	user := models.User{
		Nickname: "Sam Phil",
		Email:    "sam@gmail.com",
		Password: "password",
	}
	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.Post{}, err
	}
	post := models.Post{
		Title:    "This is the title sam",
		Content:  "This is the content sam",
		AuthorID: uint32(user.ID),
	}
	err = server.DB.Model(&models.Post{}).Create(&post).Error
	if err != nil {
		return models.Post{}, err
	}
	return post, nil
}

func SeedUsersAndPosts(server *controllers.Server) ([]models.User, []models.Post, error) {

	var err error

	if err != nil {
		return []models.User{}, []models.Post{}, err
	}
	var users = []models.User{
		{
			Nickname: "Steven victor",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		{
			Nickname: "Magu Frank",
			Email:    "magu@gmail.com",
			Password: "password",
		},
	}
	var posts = []models.Post{
		{
			Title:   "Title 1",
			Content: "Hello world 1",
		},
		{
			Title:   "Title 2",
			Content: "Hello world 2",
		},
	}

	for i := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		posts[i].AuthorID = uint32(users[i].ID)

		err = server.DB.Model(&models.Post{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
	return users, posts, nil
}
