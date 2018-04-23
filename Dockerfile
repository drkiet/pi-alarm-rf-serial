FROM golang:1.10

WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...
RUN cd src; \
    go build pi-alarm-rf-serial.go utils.go http.go udp.go;
RUN mv src/pi-alarm-rf-serial $GOPATH/bin/.
RUN echo $PATH
RUN ls -al
CMD ["pi-alarm-rf-serial"]
