package modeltests

import (
	"log"
	"testing"

	"github.com/jameslahm/bloggy_backend/models"
	. "github.com/jameslahm/bloggy_backend/tests/utils"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllPosts(t *testing.T) {
	err := RefreshUserAndPostTable(&server)
	if err != nil {
		log.Fatalf("Error refreshing user and post table %v\n", err)
	}
	_, _, err = SeedUsersAndPosts(&server)
	if err != nil {
		log.Fatalf("Error seeding user and post table %v\n", err)
	}
	posts, err := models.FindAllPosts(server.DB)
	if err != nil {
		t.Errorf("Getting the posts error: %v\n", err)
	}
	assert.Equal(t, len(posts), 2)
}

func TestGetPostByID(t *testing.T) {
	err := RefreshUserAndPostTable(&server)
	if err != nil {
		log.Fatalf("Error refreshing user and post table %v\n", err)
	}
	post, err := SeedOneUserAndOnePost(&server)
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	foundPost, err := models.FindPostById(server.DB, int(post.ID))
	if err != nil {
		t.Errorf("Getting one post error: %v\n", err)
	}
	assert.Equal(t, foundPost.ID, post.ID)
	assert.Equal(t, foundPost.Content, post.Content)
	assert.Equal(t, foundPost.Title, post.Title)
}

func TestUpdatePost(t *testing.T) {
	err := RefreshUserAndPostTable(&server)
	if err != nil {
		log.Fatalf("Err refreshing user and post table %v\n", err)
	}
	post, err := SeedOneUserAndOnePost(&server)
	if err != nil {
		log.Fatalf("Error seed one user and post %v\n", err)
	}
	var obj = map[string]interface{}{
		"title":   "Hello",
		"content": "Hello",
	}
	err = models.UpdatePost(server.DB, int(post.ID), obj)
	if err != nil {
		t.Errorf("Error UpdatePost %v\n", err)
	}
	newPost, err := models.FindPostById(server.DB, int(post.ID))
	if err != nil {
		t.Errorf("Error FindPostById %v\n", err)
	}
	assert.Equal(t, obj["title"], newPost.Title)
	assert.Equal(t, obj["content"], newPost.Content)

}

func TestDeletePost(t *testing.T) {
	err := RefreshUserAndPostTable(&server)
	if err != nil {
		log.Fatalf("Error refreshUserAndPostTable %v\n", err)
	}
	post, err := SeedOneUserAndOnePost(&server)
	if err != nil {
		log.Fatalf("Error seedOneUserAndPost %v\n", err)
	}
	err = models.DeletePost(server.DB, int(post.ID))
	if err != nil {
		t.Errorf("Error DeletePost %v\n", err)
	}
	posts, err := models.FindAllPosts(server.DB)
	if err != nil {
		t.Errorf("Error FindAllPosts %v\n", err)
	}
	assert.Equal(t, len(posts), 0)

}
