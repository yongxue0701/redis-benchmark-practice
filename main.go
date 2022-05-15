package main

import (
	r "redis-benchmark-practice/redis"
)

func main() {
	redisService := r.NewRedisService("127.0.0.1", 6379)
	redisService.Run()
}
