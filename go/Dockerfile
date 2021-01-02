# Official docker container image
FROM golang:1.15.2

# environment variables
#ENV REPO_URL=github.com/steven7/go-createmusic

#ENV GOPATH=/app

#ENV APPPATH=$GOPATH/src/$REPO_URL

#ENV WORKPATH=$GOPATH/src/$REPO_URL
#/createmusic

# Copy the local package files to the container's workspace.
#ADD . /go/src/go-createmusic

#WORKDIR /go/src/go-createmusic
#WORKDIR '~/Documents/Go/go-createmusic'

#WORKDIR $WORKPATH
WORKDIR /app/go-createmusic

#COPY . .
#COPY . $WORKPATH

COPY go.mod .
COPY go.sum .

#RUN go get -d -v ./...
#RUN go install -v ./...


#RUN go build -o go-createmusic .
RUN go mod download

COPY .. .

RUN go build -o ./out/go-createmusic

# Run the go-createmusic command by default when the container starts.
#ENTRYPOINT /go/bin/go-createmusic

EXPOSE 5000

CMD ["./go-createmusic"]


#ENTRYPOINT go run main.go