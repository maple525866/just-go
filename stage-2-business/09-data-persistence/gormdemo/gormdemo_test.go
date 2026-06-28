package gormdemo

import (
	"testing"

	"gorm.io/gorm"
)

func TestAutoMigrateAndCRUD(t *testing.T) {
	db := openGormDB(t)
	if err := AutoMigrate(db); err != nil {
		t.Fatalf("AutoMigrate returned error: %v", err)
	}

	user := User{Name: "Ada"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Create user returned error: %v", err)
	}

	var got User
	if err := db.First(&got, user.ID).Error; err != nil {
		t.Fatalf("First user returned error: %v", err)
	}
	if got.Name != "Ada" {
		t.Fatalf("user name = %q, want Ada", got.Name)
	}
}

func TestPreloadPosts(t *testing.T) {
	db := openGormDB(t)
	if err := AutoMigrate(db); err != nil {
		t.Fatalf("AutoMigrate returned error: %v", err)
	}
	if err := SeedUserWithPosts(db); err != nil {
		t.Fatalf("SeedUserWithPosts returned error: %v", err)
	}

	user, err := LoadUserWithPosts(db, "Grace")
	if err != nil {
		t.Fatalf("LoadUserWithPosts returned error: %v", err)
	}
	if len(user.Posts) != 2 {
		t.Fatalf("posts length = %d, want 2", len(user.Posts))
	}
}

func TestPreloadTags(t *testing.T) {
	db := openGormDB(t)
	if err := AutoMigrate(db); err != nil {
		t.Fatalf("AutoMigrate returned error: %v", err)
	}
	postID, err := SeedPostWithTags(db)
	if err != nil {
		t.Fatalf("SeedPostWithTags returned error: %v", err)
	}

	post, err := LoadPostWithTags(db, postID)
	if err != nil {
		t.Fatalf("LoadPostWithTags returned error: %v", err)
	}
	if len(post.Tags) != 2 {
		t.Fatalf("tags length = %d, want 2", len(post.Tags))
	}
}

func openGormDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := OpenMemory()
	if err != nil {
		t.Fatalf("OpenMemory returned error: %v", err)
	}
	return db
}
