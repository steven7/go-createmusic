# Official docker container image
FROM golang:1.14

# Copy the local package files to the container's workspace.
#ADD . /go/src/go-createmusic


#WORKDIR /go/src/go-createmusic
WORKDIR '~/Documents/Go/go-createmusic'

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

# Run the go-createmusic command by default when the container starts.
#ENTRYPOINT /go/bin/go-createmusic
ENTRYPOINT ~/Documents/Go/go-createmusic

EXPOSE 5000

CMD ["go-createmusic"]