# filepath: /home/linuxuser/PBL/MI5-PBL2/clients/Makefile
PORT ?= 8080

car:
	docker build -t car-client -f car/Dockerfile .
	docker run -it --rm --name car-client -h 172.16.103.9 -e PORT=$(PORT) -p $(PORT):8080 car-client

station:
	docker build -t station-client -f station/Dockerfile .
	docker run -it --rm --name station-client -h 172.16.103.9 -e PORT=$(PORT) -p $(PORT):8080 station-client

broker:
	docker build -t mosquitto-broker -f broker/Dockerfile broker/
	docker run  --rm --name mosquitto -p 1883:1883 -p 9001:9001 mosquitto-broker

clean:
	docker rm -f car-client station-client mosquitto 2>/dev/null || true