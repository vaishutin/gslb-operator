package storage

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vaishutin/gslb-operator/watcher/checkers"
)

// Init connect
func GetClient(addr string, password string, db int, poolName string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "qwerty",
		DB:       0,
		PoolSize: 1000,
	})

	ctx := context.Background()

	cmd := client.Ping(ctx)
	if cmd.Err() != nil {
		log.Println("Redis connect error")
		return nil, cmd.Err()
	}
	log.Println(cmd)

	cn := client.Conn()

	ttlReg := time.Duration(10) * time.Second
	resp := cn.Set(ctx, "REG_"+poolName, "Registered", ttlReg)

	log.Println(resp)
	return client, nil
}

// Write Global Name struct to redis
func InitPool(conn *redis.Conn, conf checkers.WatcherConfig) error {
	ctx := context.Background()

	data, _ := json.Marshal(conf)

	cmd := conn.Set(ctx, conf.GlobalName, data, 0)

	if cmd.Err() != nil {
		log.Println(cmd.Err())
		return cmd.Err()
	}

	return nil
}

// Write member status
func WriteStat(conn *redis.Conn, config checkers.WatcherConfig, memberName string, memberHealth bool) error {
	ctx := context.Background()

	status := checkers.HealthData{
		Health:    memberHealth,
		LastCheck: time.Now().Unix(),
		IP:        config.Members[memberName].Ip,
	}

	keyVal := config.GlobalName + "/" + memberName + "/health"
	value, _ := json.Marshal(status)

	cmd := conn.Set(ctx, keyVal, value, 0)

	return cmd.Err()

}