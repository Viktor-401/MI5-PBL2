make db:
	docker network create app-net || true
	docker run -d --name mongodb --network app-net -p 27017:27017 mongo:5.0 || true
	docker build -t app-go .
	docker run --rm --network app-net -e MONGO_URI=mongodb://mongodb:27017 -e DB_NAME=$$(date +%s) -p 8080:8080 app-go

make rundb:
	docker build -t app-go .
	docker run --rm --network app-net -e MONGO_URI=mongodb://mongodb:27017 -e DB_NAME=$$(date +%s) -p 8081:8080 app-go
