package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"os"
	"strconv"
	"syscall"
)

func main() {

	redisHost := "localhost:6379"
	if os.Getenv("GO_REDIS_ADDR") != "" {
		redisHost = os.Getenv("GO_REDIS_ADDR")
	}
	redisDB := 0
	if os.Getenv("GO_REDIS_DB") != "" {
		redisDB, _ = strconv.Atoi(os.Getenv("GO_REDIS_DB"))
	}
	redisPassword := ""
	if os.Getenv("GO_REDIS_PASSWORD") != "" {
		redisHost = os.Getenv("GO_REDIS_PASSWORD")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword, // no password set
		DB:       redisDB,       // use default DB
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatal(err.Error())
	}

	keys, err := redisClient.Do("KEYS", "MICROSERVICE_*").Result()
	if err != nil {
		fmt.Printf("error get redis : %v", err)
		return
	}else{
		for _, data := range keys.([]interface{}) {

			redisPid := redisClient.Get(data.(string)).Val()

			fmt.Println(fmt.Sprintf("redis_value : %v", redisPid))

			pid, err := strconv.Atoi(redisPid)

			if err != nil {
				fmt.Println(fmt.Sprintf("key %v invalid value", data.(string)))
				continue
			}

			p, err := os.FindProcess(pid)

			if err != nil {
				fmt.Println(fmt.Sprintf("error find : %v", err))
				continue
			}

			err = p.Signal(syscall.SIGTERM)

			if err == nil {
				fmt.Println("success")
			}else{
				fmt.Println(fmt.Sprintf("error kill : %v", err))
			}
		}
	}
}