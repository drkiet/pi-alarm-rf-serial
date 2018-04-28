FROM golang:1.10

WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...
RUN cd src; \
    go build pi_alarm_rf_serial.go utils.go http.go udp.go rf_tx_rx.go sensor.go alarm_unit.go event.go;
RUN mv src/pi_alarm_rf_serial $GOPATH/bin/.

CMD ["pi_alarm_rf_serial"]
