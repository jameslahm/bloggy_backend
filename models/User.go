package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// User user model
type User struct {
	gorm.Model
	Nickname string `gorm:"size:255;not null;unique" json:"nickname"`
	Email    string `gorm:"size:100;not null;unique" json:"email"`
	Password string `gorm:"size:100;not null;" json:"-"`
	Token    string `gorm:"-" json:"token"`
}

// Hash hash password
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword verify password
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func CreateUser(db *gorm.DB, user *User) error {
	hashedPassword, err := Hash(user.Password)
	user.Password = string(hashedPassword)
	if err != nil {
		return err
	}
	err = db.Debug().Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func FindAllUsers(db *gorm.DB) ([]User, error) {
	users := []User{}
	err := db.Debug().Model(&User{}).Find(&users).Error
	if err != nil {
		return users, err
	}
	return users, nil
}

func FindUserById(db *gorm.DB, id int) (*User, error) {
	var user User
	err := db.Debug().Model(&User{}).Where("id=?", id).First(&user).Error
	if err != nil {
		return &User{}, err
	}
	return &user, nil
}

func UpdateUser(db *gorm.DB, id int, obj map[string]interface{}) error {
	if password, ok := obj["password"]; ok {
		hashedPassword, err := Hash(password.(string))
		if err != nil {
			return err
		}
		obj["password"] = string(hashedPassword)
	}
	err := db.Debug().Model(&User{}).Where("id=?", id).Updates(obj).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(db *gorm.DB, id int) error {
	err := db.Debug().Model(&User{}).Where("id=?", id).Delete(&User{}).Error
	if err != nil {
		return err
	}
	return nil
}
