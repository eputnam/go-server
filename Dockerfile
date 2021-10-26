FROM golang:1.17

# copy server into container folder /go/src/go-server
WORKDIR /go/src/go-server
COPY . .

# install dependencies and then install the module
RUN go get -d -v ./...
RUN go install -v ./...

# expose 443 for https
EXPOSE 443

# start the server
CMD ["go-server"]
