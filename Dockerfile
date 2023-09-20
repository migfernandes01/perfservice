FROM golang:1.19.0

WORKDIR /app

# install air
RUN go install github.com/cosmtrek/air@latest

# copy all files to the container
COPY . .

# tidy go.mod
RUN go mod tidy

