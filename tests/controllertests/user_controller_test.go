package controllertests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jameslahm/bloggy_backend/models"
	"github.com/jameslahm/bloggy_backend/tests/utils"
	"gopkg.in/go-playground/assert.v1"
)

func TestRegister(t *testing.T) {
	err := utils.RefreshUserTable(&server)
	if err != nil {
		log.Fatalf("Error RefreshUserTable %v\n", err)
	}
	samples := []struct {
		inputJSON  string
		statusCode int
	}{
		{
			inputJSON:  `{"nickname":"Pet","email":"fake@fake.com","password":"123"}`,
			statusCode: 200,
		},
	}
	for _, v := range samples {
		req, err := http.NewRequest("POST", "/users", bytes.NewBuffer([]byte(v.inputJSON)))
		if err != nil {
			t.Errorf("Error NewRequest %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.Register)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, v.statusCode)
	}
}

func TestGetUsers(t *testing.T) {
	err := utils.RefreshUserTable(&server)
	if err != nil {
		log.Fatalf("Error RefreshUserTable %v\n", err)
	}
	err = utils.SeedUsers(&server)
	if err != nil {
		log.Fatalf("Error SeedUsers %v\n", err)
	}
	req, err := http.NewRequest("GET", "/users", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetUsers)
	handler.ServeHTTP(rr, req)
	var users []models.User
	err = json.Unmarshal(rr.Body.Bytes(), &users)
	if err != nil {
		t.Errorf("Error Unmarshal %v\n", err)
	}
	assert.Equal(t, len(users), 2)
	assert.Equal(t, rr.Code, http.StatusOK)
}

func TestGetUserById(t *testing.T) {
	err := utils.RefreshUserTable(&server)
	if err != nil {
		log.Fatalf("Error RefreshUserTable %v\n", err)
	}
	user, err := utils.SeedOneUser(&server)
	if err != nil {
		log.Fatalf("Error SeedOneUser %v\n", err)
	}
	req, err := http.NewRequest("GET", "/users", nil)
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(user.ID))})
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetUser)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
	var foundUser models.User
	err = json.Unmarshal(rr.Body.Bytes(), &foundUser)
	if err != nil {
		t.Errorf("Error Unmarshal %v\n", err)
	}
	assert.Equal(t, foundUser.ID, user.ID)
	assert.Equal(t, foundUser.Email, user.Email)
	assert.Equal(t, foundUser.Nickname, user.Nickname)
}
