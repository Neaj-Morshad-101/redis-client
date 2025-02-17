# Redis Client Integration in Go Applications

This guide demonstrates how to integrate Redis into your Go application using both the standard go-redis client and the KubeDB client for Kubernetes deployments.

### Prerequisites
Go Environment: Ensure Go is installed on your system.
Redis Server: Have access to a running Redis instance.
Kubernetes Cluster: For KubeDB client integration, a functional Kubernetes cluster is required.


Setting Up the Go Application
Initialize the Go Module: `go mod init your_module_name`
Install Dependencies:
```bash
go get github.com/redis/go-redis/v9
go get kubedb.dev/db-client-go/redis
```

Application Code:
Create a main.go file with the following content:

```go
package main

import (
    "context"
    "fmt"
    "log"

    go_redis "github.com/redis/go-redis/v9"
    "k8s.io/client-go/kubernetes/scheme"
    v1 "kubedb.dev/apimachinery/apis/kubedb/v1"
    "kubedb.dev/db-client-go/redis"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
    // Initialize Redis client
    rdb := go_redis.NewClient(&go_redis.Options{
        Addr:     "localhost:6379",
        Password: "", // No password set
        DB:       0,  // Use default DB
    })

    // Test the connection
    _, err := rdb.Ping(context.Background()).Result()
    if err != nil {
        log.Fatalf("Could not connect to Redis: %v", err)
    }
    fmt.Println("Connected to Redis!")

    // Set a key-value pair
    err = rdb.Set(context.Background(), "mykey", "Hello, Redis!", 0).Err()
    if err != nil {
        log.Fatalf("Could not set key: %v", err)
    }

    // Get the value
    val, err := rdb.Get(context.Background(), "mykey").Result()
    if err != nil {
        log.Fatalf("Could not get key: %v", err)
    }
    fmt.Printf("mykey: %s\n", val)

    // KubeDB Client Integration
    if err := v1.AddToScheme(scheme.Scheme); err != nil {
        log.Fatalf("Failed to add KubeDB API group to scheme: %v", err)
    }

    k8sClient, err := client.New(config.GetConfigOrDie(), client.Options{})
    if err != nil {
        log.Fatalf("Failed to create Kubernetes client: %v", err)
    }

    redisInstance := &v1.Redis{}
    err = k8sClient.Get(context.Background(), client.ObjectKey{
        Namespace: "demo",
        Name:      "rd-standalone",
    }, redisInstance)
    if err != nil {
        log.Fatalf("Failed to get Redis instance: %v", err)
    }

    builder := redis.NewKubeDBClientBuilder(k8sClient, redisInstance)
    kubedbRedisClient, err := builder.GetRedisClient(context.Background())
    if err != nil {
        log.Fatalf("Failed to get Redis client: %v", err)
    }

    // Set a key-value pair using KubeDB client
    err = kubedbRedisClient.Set(context.Background(), "KubeDBKey", "Hello, KubeDB!", 0).Err()
    if err != nil {
        log.Fatalf("Could not set key: %v", err)
    }

    // Get the value
    val, err = kubedbRedisClient.Get(context.Background(), "KubeDBKey").Result()
    if err != nil {
        log.Fatalf("Could not get key: %v", err)
    }
    fmt.Printf("KubeDBKey: %s\n", val)
}
```


Building and Deploying the Application
Build the Docker Image:

Create a Dockerfile:

```dockerfile

# Use the official Golang image
FROM golang:1.20-alpine

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Command to run the application
CMD ["./main"]
Build the Docker image:
```

`docker build -t your_dockerhub_username/your_app_name:latest .`
Push the Image to Docker Hub:
`docker push your_dockerhub_username/your_app_name:latest`
Kubernetes Deployment:

Create a deployment.yaml file:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app-deployment
  labels:
    app: my-app
spec:
  replicas: 1
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
        - name: my-app-container
          image: your_dockerhub_username/your_app_name:latest
          ports:
            - containerPort: 8080
```
Apply the deployment:
`kubectl apply -f deployment.yaml`
Notes
Environment Variables: If you're using the KubeDB client, ensure your application runs within the Kubernetes cluster where the Redis instance is deployed. The KubeDB client fetches connection details from the cluster, eliminating the need for manual environment variable configuration.

RBAC Configuration: Ensure your application's service account has the necessary permissions to access KubeDB resources and secrets. This might involve setting up appropriate ClusterRole and ClusterRoleBinding.

For kubedb client to work, we need deploy our app in k8s.
- build docker image like 'docker build -t neajmorshad/rd-client:0.0.1 .'
- `docker push neajmorshad/rd-client:0.0.1`
- `kubectl apply -f deployment.yaml`
- `kubectl logs -f my-app-deployment-7c64dc9b7b-qm9mw` 
CheckWithKubeDBClient().......
Value for 'mykey': Hello, KubeDB!
key2 does not exist




