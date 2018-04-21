FROM golang:1.10

COPY . .

RUN go get -d -v ./...
RUN go build src/pi-alarm-rf-serial.go

COPY pi-alarm-rf-serial .

CMD ["./pi-alarm-rf-serial"]
