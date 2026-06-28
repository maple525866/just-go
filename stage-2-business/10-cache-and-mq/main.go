package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"just-go/stage-2-business/10-cache-and-mq/cachex"
	"just-go/stage-2-business/10-cache-and-mq/mqdemo"
)

func main() {
	ctx := context.Background()
	store := cachex.NewStore(time.Now())
	source := map[string]string{"article:1": "Cache-Aside keeps hot data close"}
	loader := func(context.Context, string) (string, bool, error) {
		value, ok := source["article:1"]
		return value, ok, nil
	}
	value, ok, err := cachex.CacheAside(ctx, store, "article:1", time.Minute, loader)
	if err != nil {
		log.Fatal(err)
	}

	readThrough := cachex.NewReadThrough(cachex.NewStore(time.Now()), time.Minute, loader)
	_, _, _ = readThrough.Get(ctx, "article:1")
	writeThrough := cachex.NewWriteThrough(cachex.NewStore(time.Now()), time.Minute, func(ctx context.Context, key, value string) error {
		source[key] = value
		return nil
	})
	if err := writeThrough.Set(ctx, "article:2", "Write-Through updates source and cache"); err != nil {
		log.Fatal(err)
	}

	lockStore := cachex.NewStore(time.Now())
	lock := cachex.NewLockManager(lockStore)
	locked := lock.Acquire("article:1", "token-a", time.Minute)

	broker := mqdemo.NewBroker(time.Now(), time.Minute)
	broker.Publish("article.created", "article-1")
	delivery, delivered := broker.Fetch()
	if delivered {
		broker.Ack(delivery.DeliveryID)
	}

	fmt.Println("第 10 章：缓存与消息")
	fmt.Printf("Cache-Aside: value=%q ok=%v\n", value, ok)
	fmt.Println("Read-Through: cache owns the loader and fills itself on miss")
	fmt.Println("Write-Through: source and cache are updated together")
	fmt.Printf("缓存问题: negative cache 防穿透，TTL jitter 防雪崩，singleflight 防击穿（示例 TTL=%s）\n", cachex.JitterTTL(time.Minute, 10, "article:1"))
	fmt.Printf("分布式锁: token + TTL 获取结果=%v\n", locked)
	fmt.Printf("消息 ack: topic=%q delivered=%v ack 后不再投递\n", delivery.Topic, delivered)
}
