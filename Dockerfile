FROM golang:latest
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go mod verify

COPY . .

RUN go build -o main .
EXPOSE 3000
CMD ["./main"]

