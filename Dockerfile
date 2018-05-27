FROM golang:1.10

WORKDIR /go/src/app

COPY src src
COPY web web
COPY pi_alarm.cfg ./.private

RUN ls -al

RUN go get -d -v ./...

RUN cd src; \
    go build alarm_main.go email.go ethernet.go event.go http.go log_util.go monitor.go pi_alarm.go pi_alarm_cfg.go pi_network.go rf_base.go sensor.go udp.go;

RUN mv src/alarm_main $GOPATH/bin/.

RUN rm -r src/*
RUN rmdir src

CMD ["alarm_main"]
