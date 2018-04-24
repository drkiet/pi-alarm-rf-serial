FROM golang:1.10

WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...
RUN cd src; \
    go build pi_alarm_rf_serial.go utils.go http.go udp.go rf_tx_rx.go sensor.go;
RUN mv src/pi-alarm-rf-serial $GOPATH/bin/.
RUN echo $PATH
RUN ls -al
CMD ["pi_alarm_rf_serial"]
