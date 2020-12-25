#!/bin/bash
# IP_ADRR = 157.245.126.199
# APP_NAME = "createmusic"
# Change to the directory with our code that we plan to work from

# !!!!!!!!!!!!!!!!!!!!!!!!!!!!
#
# was written with gopath wont work with modules
#
# !!!!!!!!!!!!!!!!!!!!!!!!!!!!

cd $GOPATH/src/go-createmusic
echo "=== Releasing createmusic ==="
echo "  Deleting the local binary if it exists (so it isn't uploaded)..."
rm createmusic
echo "  Done!"

echo "  Deleting existing code..."
ssh root@157.245.126.199 "rm -rf /root/go/src/go-createmusic"
echo "  Code deleted successfully!"

echo "  Uploading code..."
# The \ at the end of the line tells bash that our
# command isn't done and wraps to the next line
rsync -avr --exclude '.git/*' --exclude 'tmp/*' \
  --exclude 'image/*' ./ \
  root@157.245.126.199:/root/go/src/go-createmusic/
echo "  Code uploaded successfully!"

# Note: Whenever you add new packages to your source code you will need
# to update this section of your release script, or alternatively you
# could look into using a package manager like dep13.

echo "  Go getting deps..."
ssh root@157.245.126.199 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get golang.org/x/crypto/bcrypt"
ssh root@157.245.126.199 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/gorilla/mux"
ssh root@157.245.126.199 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/gorilla/schema"
ssh root@157.245.126.199 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/lib/pq"
ssh root@157.245.126.199 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/jinzhu/gorm"
ssh root@157.245.126.199 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/gorilla/csrf"
ssh root@157.245.126.199 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get gopkg.in/mailgun/mailgun-go.v1"

echo "  Building the code on remote server..."
ssh root@157.245.126.199 'export GOPATH=/root/go; \
  cd /root/app; \
  /usr/local/go/bin/go build -o ./server \
    $GOPATH/src/go-createmusic/*.go'
echo "  Code built successfully!"

echo "  Moving assets..."
ssh root@157.245.126.199 "cd /root/app; \
  cp -R /root/go/src/go-createmusic/assets ."
echo "  Assets moved successfully!"

echo "  Moving views..."
ssh root@157.245.126.199 "cd /root/app; \
  cp -R /root/go/src/go-createmusic/views ."
echo "  Views moved successfully!"

echo "  Moving Caddyfile..."
ssh root@157.245.126.199 "cd /root/app; \
  cp /root/go/src/go-createmusic/Caddyfile ."
echo "  Views moved successfully!"


echo "  Restarting the server..."
ssh root@157.245.126.199 "sudo service createmusic restart"
echo "  Server restarted successfully!"

echo "  Restarting Caddy server..."
ssh root@157.245.126.199 "sudo service caddy restart"
echo "  Caddy restarted successfully!"

echo "==== Done releasing createmusic ===="