package gormdemo

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User demonstrates one-to-many associations.
type User struct {
	ID    uint
	Name  string
	Posts []Post
}

// Post belongs to a user and can have many tags.
type Post struct {
	ID     uint
	Title  string
	Body   string
	UserID uint
	Tags   []Tag `gorm:"many2many:post_tags;"`
}

// Tag demonstrates a many-to-many relation with posts.
type Tag struct {
	ID   uint
	Name string `gorm:"uniqueIndex"`
}

// OpenMemory opens a shared GORM SQLite memory database.
func OpenMemory() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("file:just_go_stage09_gorm?mode=memory&cache=shared"), &gorm.Config{})
}

// AutoMigrate creates the tables used in the examples.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Post{}, &Tag{})
}

// SeedUserWithPosts creates a user with two posts.
func SeedUserWithPosts(db *gorm.DB) error {
	user := User{
		Name: "Grace",
		Posts: []Post{
			{Title: "Preload avoids N+1", Body: "Load associations with the parent query."},
			{Title: "Transactions protect invariants", Body: "Commit all changes or roll them back."},
		},
	}
	return db.Create(&user).Error
}

// LoadUserWithPosts loads posts with the user.
func LoadUserWithPosts(db *gorm.DB, name string) (User, error) {
	var user User
	err := db.Preload("Posts").Where("name = ?", name).First(&user).Error
	return user, err
}

// SeedPostWithTags creates one post with two tags.
func SeedPostWithTags(db *gorm.DB) (uint, error) {
	post := Post{
		Title: "GORM associations",
		Body:  "Many-to-many associations use a join table.",
		Tags:  []Tag{{Name: "gorm"}, {Name: "association"}},
	}
	if err := db.Create(&post).Error; err != nil {
		return 0, err
	}
	return post.ID, nil
}

// LoadPostWithTags loads tags with the post.
func LoadPostWithTags(db *gorm.DB, id uint) (Post, error) {
	var post Post
	err := db.Preload("Tags").First(&post, id).Error
	return post, err
}
