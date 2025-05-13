hostdb:
	docker run -d --name mongodb --network host -p 27017:27017 mongo:5.0 || true
	docker build -t app-go .
	docker run -it --rm --name server1 -e MONGO_URI=mongodb://172.25.236.13:27017 -e DB_NAME=server1 -e PORT=$(PORT) -p $(PORT):$(PORT) app-go

hostrundb:
	docker build -t app-go . || true
	docker run --name server2 --rm -e MONGO_URI=mongodb://172.25.236.13:27017 -e DB_NAME=server2 -e PORT=$(PORT) -p $(PORT):$(PORT) app-go


db:
	make killmongo
	docker network create app-net || true
	docker run -d --name mongodb --network app-net -p 27017:27017 mongo:5.0 || true
	docker build -t app-go .
	docker run -it --name $(NAME) --rm --network app-net \
		-e MONGO_URI=mongodb://mongodb:27017 \
		-e DB_NAME=$(NAME) \
		-e PORT=$(PORT) \
		-p $(PORT):$(PORT) \
		app-go

rundb:
	docker build -t app-go .
	docker run -it --name $(NAME) --rm --network app-net \
		-e MONGO_URI=mongodb://mongodb:27017 \
		-e DB_NAME=$(NAME) \
		-e PORT=$(PORT) \
		-p $(PORT):$(PORT) \
		app-go
killmongo:
	docker stop mongodb
	docker rm mongodb