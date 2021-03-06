package modeltests

import (
	"log"
	"os"
	"testing"

	"github.com/jameslahm/bloggy_backend/api/controllers"
	"github.com/jameslahm/bloggy_backend/tests/utils"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	utils.Database(&server)
	os.Exit(m.Run())
}
