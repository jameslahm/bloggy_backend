package controllertests

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jameslahm/bloggy_backend/tests/utils"
	"gopkg.in/go-playground/assert.v1"
)

func TestLogin(t *testing.T) {
	err := utils.RefreshUserTable(&server)
	if err != nil {
		log.Fatalf("Error RefreshUserTable %v\n", err)
	}
	_, err = utils.SeedOneUser(&server)
	if err != nil {
		log.Fatalf("Error SeedOneUser %v\n", err)
	}
	samples := []struct {
		inputJSON  string
		statusCode int
	}{
		{
			inputJSON:  `{"email":"123","password":"123"}`,
			statusCode: 400,
		},
		{
			inputJSON:  `{"email":"fake@example.com","password":"password"}`,
			statusCode: 200,
		},
	}
	for _, v := range samples {
		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(v.inputJSON)))
		if err != nil {
			t.Errorf("Error NewRequest %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.Login)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)
	}
}
