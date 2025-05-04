# Dockerfile
FROM golang:1.21

WORKDIR /app

# copia dependências
COPY go.mod go.sum ./
RUN go mod download

# copia o código
COPY . ./

# compila o binário
RUN go build -o app

# entrada
CMD ["./app"]
