package redis

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/ory/dockertest"
)

var (
	client    *redis.Client
	redisPort string
)

func TestMain(m *testing.M) {
	var closeRedis func()
	var err error

	client, redisPort, closeRedis, err = startRedis()
	defer closeRedis()
	if err != nil {
		log.Println(err)
	}

	m.Run()
}

// startRedis はDockerを使用してRedisコンテナを起動し、redisへの接続を確立する関数です。
func startRedis() (*redis.Client, string, func(), error) {
	// Dockerのデフォルト接続方法を使用（Windowsではtcp/http、Linux/OSXではsocket）
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Printf("Could not construct pool: %s\n", err)
		return nil, "", nil, err
	}

	// Dockerに接続を試みる
	err = pool.Client.Ping()
	if err != nil {
		log.Printf("Could not connect to Docker: %s", err)
		return nil, "", nil, err
	}

	// Dockerコンテナを起動する際に指定する設定定義
	redisOptions := &dockertest.RunOptions{
		Repository: "redis",
		Tag:        "5.0",
		Env: []string{
			"REDIS_PASSWORD=",
		},
	}

	redisResource, err := pool.RunWithOptions(redisOptions)
	if err != nil {
		log.Printf("Could not start Redis resource: %s", err)
		return nil, "", nil, err
	}

	// Redisのポートを取得
	redisPort = redisResource.GetPort("6379/tcp")

	// Redisへの接続確認
	err = pool.Retry(func() error {
		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("localhost:%s", redisPort),
			Password: "",
			DB:       0,
		})
		cmd := client.Ping(context.Background())
		_, err = cmd.Result()
		return err
	})
	if err != nil {
		log.Printf("Could not connect to Redis container: %s", err)
		return nil, "", nil, err
	}

	log.Println("start Redis container🐳")

	// redisへの接続とクリーンアップ関数を返却
	return client, redisPort, func() { closeRedis(client, pool, redisResource) }, nil
}

// closeMySQL はMySQLデータベースの接続を閉じ、Dockerコンテナを停止・削除する関数
func closeRedis(client *redis.Client, pool *dockertest.Pool, resource *dockertest.Resource) {
	// redisへの接続を切断
	if err := client.Close(); err != nil {
		log.Fatalf("Failed to close redis: %s", err)
	}

	// Dockerコンテナを停止して削除
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Failed to purge resource: %s", err)
	}

	log.Println("close Redis container🐳")
}

func ValidateErr(t *testing.T, err error, wantErr error) {
	if (err != nil) != (wantErr != nil) {
		t.Errorf("error = %v, wantErr %v", err, wantErr)
	} else if err != nil && wantErr != nil && err.Error() != wantErr.Error() {
		t.Errorf("error = %v, wantErr %v", err, wantErr)
	}
}
