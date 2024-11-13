package redis

import (
	"fmt"
	"jdy/config"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
)

func Init() {

	var (
		conf = config.Config.Redis
	)

	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.Db,
	})
}

func Close() {
	Client.Close()
}
