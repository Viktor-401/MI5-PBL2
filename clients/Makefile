hostcar:
	docker build -t car-client -f car/Dockerfile .
	docker run -it --rm -e PORT=$(PORT) -p $(PORT):8080 car-client

hoststation:
	docker build -t station-client -f station/Dockerfile .
	docker run -it --rm -e PORT=$(PORT) -p $(PORT):8080 station-client

hostbroker:
	docker build -t mosquitto-broker -f broker/Dockerfile broker/
	docker run  --rm --name mosquitto -p 1885:1885 -p 9001:9001 mosquitto-broker

car:
	docker build -t car-client -f car/Dockerfile .
	docker run -it --rm --network app-net -e PORT=$(PORT) -p $(PORT):8080 car-client

station:
	docker build -t station-client -f station/Dockerfile .
	docker run -it --rm --network app-net -e PORT=$(PORT) -p $(PORT):8080 station-client

broker:
	docker build -t mosquitto-broker -f broker/Dockerfile broker/
	docker run  --rm --name mosquitto -p 1885:1885 -p 9001:9001 --network app-net  mosquitto-broker