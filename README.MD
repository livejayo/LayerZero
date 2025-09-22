Part 1 - Monitring System

Use the docker-compose.yml file which builds the monitor image and starts both containers

Part 2 - Monitor deploy

Place the main.go, Dockerfile (from the previous answer), and docker-compose.yml in a directory.

Initialize the Go module with go mod init
goat-monitor and then run go mod tidy.

run docker-compose up --build.

visit http://localhost:8080 in your browser to see the status.


Part 3 - Instructions to Deploy

Build and load your monitor image into your local Kubernetes cluster (e.g., minikube image load goat-monitor:latest).

Save the YAML above as goatmon.yaml.

Run kubectl apply -f goatmon.yaml.

Find your cluster's IP (e.g., minikube ip) and then access the monitor in your browser at http://<your-cluster-ip>:30001.



