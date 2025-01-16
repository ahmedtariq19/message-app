# Event-Driven Real-Time Messaging Application

This project is an event-driven, real-time messaging application that follows a microservices architecture. Follow the steps below to set up and run the application.

## Prerequisites

1. **Go**: Install Go by following the instructions [here](https://go.dev/doc/install).
2. **Docker**: Install Docker by following the instructions [here](https://docs.docker.com/engine/install/).
3. **Docker Compose**: Install Docker Compose by following the official documentation.
4. **Make**: Install Make utility.

## Setup Instructions

### 1. Run Project Dependencies in Docker

Start the project dependencies using Docker Compose:

```bash
docker-compose -f docker-compose-deps.yaml up -d
```

### 2. Configure the Application

- Copy the `example_conf.json` to `conf.json`:

```bash
cp example_conf.json conf.json 
```

### 3. Run the Project

Run the application using the following command:

```bash
go run main.go

```
### 4. Run Tests

To run the tests, follow these steps:

1. Install mock generation:

```bash
make install-mockgen
```

2. If the path is not set, run:

```bash
make path
```

3. Generate mocks:

```bash
make mocks
```

4. Run unit and integration tests:

```bash
make test
```
This command will run both unit tests and integration tests. A comprehensive test report will be generated as cover.html.


## 5. Deployment with Minikube and Helm

1. Start Minikube

Start Minikube with the following command:

```bash
minikube start
```

2. Verify Minikube Status

Check the status of Minikube to ensure it's running:

```bash
minikube status
```

3. Install the Application with Helm

To install the application using Helm, run:

```bash
helm install message-app message-app
```

4. Verify Deployment

To verify the deployment of your application, run:

```bash
kubectl get pods
```

## Clean Up Deployment

1. Stop Minikube

To stop the Minikube cluster, run:

```bash
minikube stop
```

2. Uninstall the Helm Application

To uninstall the Helm app, run:

```bash
helm uninstall message-app
```

3. To stop the Docker dependencies, run:

```bash
docker-compose -f docker-compose-deps.yaml down
```