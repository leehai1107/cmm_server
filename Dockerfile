FROM golang:1.22.0

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy

EXPOSE 8081

CMD ["go", "run", "./main.go", "api"]

