FROM golang:1.17

# copy server into container folder /go/src/github.com/eputnam/health-check-server
WORKDIR /go/src/github.com/eputnam/health-check-server
COPY . .

# install dependencies and then install the module
RUN go get -d -v ./...
RUN go install -v ./...

# expose 8080 for http
EXPOSE 8080

# start the server
CMD ["health-check-server"]
