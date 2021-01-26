package modeltests

import (
	"github.com/jameslahm/bloggy_backend/models"
	. "github.com/jameslahm/bloggy_backend/tests/utils"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"testing"
)

func TestFindAllUser(t *testing.T) {
	err := RefreshUserTable(&server)
	if err != nil {
		log.Fatal(err)
	}
	err = SeedUsers(&server)
	if err != nil {
		log.Fatal(err)
	}
	users, err := models.FindAllUsers(server.DB)
	if err != nil {
		t.Errorf("Getting users %v\n", err)
		return
	}
	assert.Equal(t, len(users), 2)
}

func TestGetUserById(t *testing.T) {
	err := RefreshUserTable(&server)
	if err != nil {
		log.Fatal(err)
	}
	user, err := SeedOneUser(&server)
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}
	foundUser, err := models.FindUserById(server.DB, int(user.ID))
	if err != nil {
		t.Errorf("Getting user error: %v\n", err)
		return
	}
	assert.Equal(t, foundUser.Email, user.Email)
	assert.Equal(t, foundUser.Password, user.Password)
}

func TestUpdateUser(t *testing.T) {
	err := RefreshUserTable(&server)
	if err != nil {
		log.Fatal(err)
	}
	user, err := SeedOneUser(&server)
	if err != nil {
		log.Fatalf("Cannot seed user: %v\n", err)
	}
	var obj = map[string]interface{}{
		"email": "updated@gmail.com",
	}
	err = models.UpdateUser(server.DB, int(user.ID), obj)
	if err != nil {
		t.Errorf("Update user error: %v\n", err)
		return
	}
	updatedUser, err := models.FindUserById(server.DB, int(user.ID))
	assert.Equal(t, updatedUser.Email, obj["email"])
}

func TestDeleteUser(t *testing.T) {
	err := RefreshUserTable(&server)
	if err != nil {
		log.Fatal(err)
	}
	user, err := SeedOneUser(&server)
	if err != nil {
		log.Fatalf("Cannot seed user: %v\n", err)
	}
	err = models.DeleteUser(server.DB, int(user.ID))
	if err != nil {
		t.Errorf("Delete user error: %v\n", err)
		return
	}
	users, err := models.FindAllUsers(server.DB)
	assert.Equal(t, len(users), 0)
}
