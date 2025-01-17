FROM golang:1.23-alpine

WORKDIR /app

# Copy go.mod and go.sum first
COPY go.mod go.sum ./

# Now download dependencies
RUN go mod download

# Then copy the rest of the code
COPY . .

RUN go build -o main main.go

EXPOSE 8080

CMD ["./main"]