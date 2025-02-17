package main

import (
	"context"
	"fmt"
	go_redis "github.com/redis/go-redis/v9"
	"log"
)

var redisClient *go_redis.Client

// Connecting via a redis url
func RedisClientWithURI() *go_redis.Client {
	//url := "redis://user:password@localhost:6379/0?protocol=3"
	url := "redis://localhost:6379"
	opts, err := go_redis.ParseURL(url)
	if err != nil {
		panic(err)
	}

	return go_redis.NewClient(opts)
}

func initRedis() {
	// Initialize the Redis client
	redisClient = go_redis.NewClient(&go_redis.Options{
		Addr:     "localhost:6380",   // Redis server address
		Password: "8fQl3LtXHlcEdiCF", // no password set
		DB:       0,                  // Default DB
	})

	// Test the connection
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis!")
}

func main() {
	//Initialize Redis client
	initRedis()

	CheckWithClient(redisClient)
	CheckWithClient(RedisClientWithURI())

	// Use KubeDB  db-client-go methods
	//CheckWithKubeDBClient()
}

func CheckWithClient(rdb *go_redis.Client) {
	// Set a key-value pair in Redis
	err := rdb.Set(context.Background(), "mykey", "Hello, Redis!", 0).Err()
	if err != nil {
		log.Fatalf("Could not set key: %v", err)
	}

	// Get the value of the key
	val, err := rdb.Get(context.Background(), "mykey").Result()
	if err != nil {
		log.Fatalf("Could not get key: %v", err)
	}
	fmt.Printf("Value for 'mykey': %s\n", val)

	val2, err := rdb.Get(context.Background(), "key2").Result()
	if err == go_redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
}

//func CheckWithKubeDBClient() {
//	// Create a Kubernetes client (using client-go or controller-runtime)
//	k8sClient, err := client.New(config.GetConfigOrDie(), client.Options{})
//	if err != nil {
//		log.Fatalf("Failed to create Kubernetes client: %v", err)
//	}
//
//	// Assume that you already have a Redis instance (for example, from KubeDB)
//	redisInstance := &v1.Redis{} // Get your Redis instance object from the Kubernetes API
//	// You should already have the Redis instance loaded from your Kubernetes cluster, or it can be created
//
//	// Use the KubeDBClientBuilder to build the Redis client
//	builder := NewKubeDBClientBuilder(k8sClient, redisInstance)
//
//	// Optionally, set the pod name, URL, or database index if needed
//	// builder.WithPod("my-redis-pod")
//	// builder.WithURL("my-redis-url")
//	// builder.WithDatabase(1) // Optional database index
//
//	// Get the Redis client
//	redisClient, err := builder.GetRedisClient(context.Background())
//	if err != nil {
//		log.Fatalf("Failed to get Redis client: %v", err)
//	}
//
//	// Now you can interact with Redis using the returned redisClient
//	// Set a key-value pair
//	err = redisClient.Set(context.Background(), "mykey", "Hello, Redis!", 0).Err()
//	if err != nil {
//		log.Fatalf("Could not set key: %v", err)
//	}
//
//	// Get the value of the key
//	val, err := redisClient.Get(context.Background(), "mykey").Result()
//	if err != nil {
//		log.Fatalf("Could not get key: %v", err)
//	}
//	fmt.Printf("Value for 'mykey': %s\n", val)
//}
