docker network create my-network
docker build -t my-mongo2 -f .\phase1\database\Dockerfile.db .\phase1\database
docker build -t my-go -f .\phase1\Dockerfile .\phase1
docker build -t my-angular -f .\clinic-appointment-app\Dockerfile .\clinic-appointment-app
docker run -d -p 27017:27017 --name mongo-container --net my-network my-mongo2
docker run -d -p 8080:8080 --name go-container --net my-network my-go
docker run -d -p 4200:4200 --name angular-container --net my-network my-angular
