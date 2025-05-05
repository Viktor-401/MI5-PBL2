# Use the official Eclipse Mosquitto image as the base image
FROM eclipse-mosquitto:latest

# Copy custom configuration file if needed
# Uncomment the following lines if you have a custom mosquitto.conf
# COPY ./config/config.conf /mosquitto/config/mosquitto.conf

# Expose the default MQTT ports
EXPOSE 1883 9001

RUN apt update && \
    apt install -y mosquitto-clients

# Start Mosquitto broker
CMD ["mosquitto", "-c", "/mosquitto/config/mosquitto.conf"]