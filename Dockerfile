FROM golang:1.10

WORKDIR /go/src/app

COPY . .

RUN ls -al

RUN rm .private

RUN go get -d -v ./...
RUN cd src; \
    go alarm_main.go email.go ethernet.go event.go http.go log_util.go monitor.go pi_alarm.go pi_alarm_cfg.go pi_network.go rf_base.go sensor.go;
RUN mv src/pi_alarm_rf_serial $GOPATH/bin/.

CMD ["pi_alarm_rf_serial"]
