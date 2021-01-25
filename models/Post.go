package models

import (
	"errors"

	"github.com/jameslahm/bloggy_backend/utils"
	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model
	Title    string `gorm:"size:255;not null;unique" json:"title"`
	Content  string `gorm:"size:255;not null;unique" json:"content"`
	Author   User   `json:"author"`
	AuthorID uint32 `gorm:"not null" json:"author_id"`
}

func (p *Post) Validate() error {
	if p.Title == "" {
		return errors.New(utils.EMPTY_TITLE)
	}
	if p.Content == "" {
		return errors.New(utils.EMPTY_CONTENT)
	}
	if p.AuthorID < 1 {
		return errors.New(utils.EMPTY_AUTHOR)
	}
	return nil
}

func CreatePost(db *gorm.DB, post *Post) error {
	err := db.Debug().Model(&Post{}).Create(post).Error
	if err != nil {
		return err
	}
	return nil
}

func FindAllPosts(db *gorm.DB) ([]Post, error) {
	var posts []Post
	err := db.Debug().Model(&Post{}).Find(&posts).Error
	if err != nil {
		return []Post{}, err
	}
	for _, post := range posts {
		err := db.Debug().Model(&User{}).Where("id=?", post.AuthorID).Take(&post.Author).Error
		if err != nil {
			return []Post{}, err
		}
	}
	return posts, nil
}

func FindPostById(db *gorm.DB, id int) (*Post, error) {
	var post Post
	err := db.Debug().Model(&Post{}).Where("id=?", id).First(&post).Error
	if err != nil {
		return &Post{}, err
	}
	err = db.Debug().Model(&User{}).Where("id=?", post.AuthorID).Take(&post.Author).Error
	if err != nil {
		return &Post{}, err
	}
	return &post, nil
}

func UpdatePost(db *gorm.DB, id int, obj map[string]interface{}) error {
	err := db.Debug().Model(&Post{}).Where("id=?", id).Updates(obj).Error
	if err != nil {
		return err
	}
	return nil
}

func DeletePost(db *gorm.DB, id int) error {
	var post Post
	err := db.Debug().Model(&Post{}).Where("id=?", id).Delete(&post).Error
	if err != nil {
		return err
	}
	return nil
}
