package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/micro/go-log"
	"os"
	"strconv"
	"strings"
)

const (
	REDIS_HOST="localhost:6379"
	REDIS_PASSWORD=""
	REDIS_DB=0
	REDIS_HOST_ENV="GO_REDIS_ADDR"
	REDIS_PASSWORD_ENV="GO_REDIS_PASSWORD"
	REDIS_DB_ENV="GO_REDIS_DB"
	DEFAULT_PATH_FILE="/home/sysadmin/go/src/microservice/setRedis.txt"
)


func main(){

	pathFile := flag.String("path_file", "", "Path File Redis Key Value")

	flag.Parse()

	redisHost := REDIS_HOST
	if os.Getenv(REDIS_HOST_ENV) != "" {
		redisHost = os.Getenv(REDIS_HOST_ENV)
	}
	redisDB := REDIS_DB
	if os.Getenv(REDIS_DB_ENV) != "" {
		redisDB, _ = strconv.Atoi(os.Getenv(REDIS_DB_ENV))
	}
	redisPassword := REDIS_PASSWORD
	if os.Getenv(REDIS_PASSWORD_ENV) != "" {
		redisHost = os.Getenv(REDIS_PASSWORD_ENV)
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

	if *pathFile == ""{
		*pathFile = DEFAULT_PATH_FILE
	}

	file, err := os.Open(*pathFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()


	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
		keyVal := strings.SplitN(scanner.Text(), ":", 2)

		if len(keyVal) == 2 {
			redisClient.Set(strings.TrimSpace(keyVal[0]), strings.TrimSpace(keyVal[1]), 0)
		}
	}
}