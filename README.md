# Redis Client Integration in Go Applications



For kubedb client to work, we need deploy our app in k8s.
- build docker image like 'docker build -t neajmorshad/rd-client:0.0.1 .'
- `docker push neajmorshad/rd-client:0.0.1`
- `kubectl apply -f deployment.yaml`
- `kubectl logs -f my-app-deployment-7c64dc9b7b-qm9mw` 
CheckWithKubeDBClient().......
Value for 'mykey': Hello, KubeDB!
key2 does not exist




