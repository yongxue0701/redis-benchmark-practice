package redis

import (
	"context"
	"fmt"
	"net"

	"github.com/go-redis/redis/v8"
	"github.com/hhxsv5/go-redis-memory-analysis"
)

type RedisService struct {
	host   string
	port   uint16
	client *redis.Client
}

func NewRedisService(host string, port uint16) *RedisService {
	return &RedisService{
		host: host,
		port: port,
	}
}

func (s *RedisService) Run() {
	address := net.JoinHostPort(s.host, string(s.port))
	client := redis.NewClient(&redis.Options{
		Addr:         address,
		Password:     "",
		DB:           0,
		PoolSize:     128,
		MinIdleConns: 100,
		MaxRetries:   5,
	})
	s.client = client

	s.writeToRedis(10000, "size_10_num_10k", s.getValue(10))
	s.writeToRedis(50000, "size_10_num_50k", s.getValue(10))
	s.writeToRedis(200000, "size_10_num_200k", s.getValue(10))
	s.writeToRedis(500000, "size_10_num_500k", s.getValue(10))

	s.writeToRedis(10000, "size_1000_num_10k", s.getValue(1000))
	s.writeToRedis(50000, "size_1000_num_50k", s.getValue(1000))
	s.writeToRedis(200000, "size_1000_num_200k", s.getValue(1000))
	s.writeToRedis(500000, "size_1000_num_500k", s.getValue(1000))

	s.writeToRedis(10000, "size_5000_num_10k", s.getValue(5000))
	s.writeToRedis(50000, "size_5000_num_50k", s.getValue(5000))
	s.writeToRedis(200000, "size_5000_num_200k", s.getValue(5000))
	s.writeToRedis(500000, "size_5000_num_500k", s.getValue(5000))

	s.generateReport()
}

func (s *RedisService) writeToRedis(num int, key string, value string) {
	for i := 0; i < num; i++ {
		key := s.getKey(key, i)

		cmd := s.client.Set(context.Background(), key, value, -1)
		if cmd.Err() != nil {
			fmt.Println(cmd.String())
		}
	}
}

func (s *RedisService) getKey(key string, index int) string {
	return fmt.Sprintf("%s_%d", key, index)
}

func (s *RedisService) getValue(size int) string {
	arr := make([]byte, size)
	for idx := 0; idx < len(arr); idx++ {
		arr[idx] = 'k'
	}
	return string(arr)
}

func (s *RedisService) generateReport() {
	conn, err := gorma.NewAnalysisConnection(s.host, s.port, "")
	if err != nil {
		fmt.Println("failed to connect server: ", err)
		return
	}
	defer conn.Close()

	conn.Start([]string{":"})

	err = conn.SaveReports("./reports")
	if err == nil {
		fmt.Println("done")
	} else {
		fmt.Println("failed to save report: ", err)
	}
}
