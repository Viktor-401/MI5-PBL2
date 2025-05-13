make db:
	docker run -d --name mongodb --network host -p 27017:27017 mongo:5.0 || true
	docker build -t app-go .
	docker run -it --rm --name server1 -e MONGO_URI=mongodb://172.16.103.10:27017 -e DB_NAME=server1 -e PORT=$(PORT) -p $(PORT):$(PORT) app-go

make rundb:
	docker build -t app-go . || true
	docker run --name server2 --rm -e MONGO_URI=mongodb://172.16.103.10:27017 -e DB_NAME=server2 -e PORT=$(PORT) -p $(PORT):$(PORT) app-go

killmongo:
	docker stop mongodb
	docker rm mongodb