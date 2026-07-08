package main

import (
	"fmt"
	"log"
	"time"

	"just-go/stage-2-business/09-data-persistence/dbx"
	"just-go/stage-2-business/09-data-persistence/gormdemo"
	"just-go/stage-2-business/09-data-persistence/sqlcrud"
	"just-go/stage-2-business/09-data-persistence/txdemo"
)

func main() {
	db, err := dbx.OpenMemory()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	pool := dbx.PoolConfig{MaxOpenConns: 4, MaxIdleConns: 2, ConnMaxLifetime: 30 * time.Minute}
	dbx.ConfigurePool(db, pool)
	if err := dbx.ApplyMigration(db, "stage-2-business/09-data-persistence/migrations/001_create_articles.sql"); err != nil {
		log.Fatal(err)
	}

	repo := sqlcrud.NewRepository(db)
	article, err := repo.Create("database/sql CRUD", "Prepared statements bind user input safely.", "gopher")
	if err != nil {
		log.Fatal(err)
	}
	if err := txdemo.CommitTwoArticles(db); err != nil {
		log.Fatal(err)
	}

	gormDB, err := gormdemo.OpenMemory()
	if err != nil {
		log.Fatal(err)
	}
	if err := gormdemo.AutoMigrate(gormDB); err != nil {
		log.Fatal(err)
	}
	if err := gormdemo.SeedUserWithPosts(gormDB); err != nil {
		log.Fatal(err)
	}
	user, err := gormdemo.LoadUserWithPosts(gormDB, "Grace")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("第 09 章：数据持久化")
	fmt.Println("database/sql: CRUD + prepared statement 参数绑定")
	fmt.Printf("连接池: %s\n", pool.Summary())
	fmt.Printf("迁移: articles 表已创建，示例文章 ID=%d Title=%q\n", article.ID, article.Title)
	fmt.Println("事务: CommitTwoArticles 成功提交两条记录，RollbackOnError 演示失败回滚")
	fmt.Printf("GORM: AutoMigrate + Preload 加载用户 %q 的 %d 篇文章\n", user.Name, len(user.Posts))
	fmt.Println("N+1: 使用 Preload 预加载关联，避免循环中逐条查询")
}
