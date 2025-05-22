hostserver:
	docker build -t app-go .
	docker run -it --name $(NAME) --rm -e MONGO_URI=mongodb://$(DBIP):27017 -e DB_NAME=$(NAME) -e PORT=$(PORT) -p $(PORT):$(PORT) --network host app-go

server:
	docker network create app-net || true
	docker build -t app-go .
	docker run -it --name $(NAME) --rm --network app-net \
		-e MONGO_URI=mongodb://mongodb:27017 \
		-e DB_NAME=$(NAME) \
		-e PORT=$(PORT) \
		-p $(PORT):$(PORT) \
		app-go

mongo:
	docker network create app-net || true
	docker run -d --name mongodb --network app-net -p 27017:27017 mongo:5.0 || true

hostmongo:
	docker run -d --name mongodb -p 27017:27017 mongo:5.0 || true

killmongo:
	docker stop mongodb
	docker rm mongodb