# Official docker container image
FROM golang:1.15.6

ENV GO111MODULE=on

WORKDIR /app

# copy go mod files first
COPY go/go.mod go/go.sum ./

# download go mod dependencies
RUN go mod download

# copy the rest of the code
COPY go .

# go build and name output file
RUN go build -o go-createmusic .

#expore port to outside the container
EXPOSE 5000

# run go output file
CMD ["./go-createmusic"]
