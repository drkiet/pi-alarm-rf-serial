FROM golang:latest 

# RUN mkdir /go/src/github.com/drkiet/pi-alarm-rf-serial
# RUN mkdir /go/src/github.com/mikepb/go-serial

ADD . /go/src/github.com/drkiet/pi-alarm-rf-serial

WORKDIR /go/src/github.com/drkiet/pi-alarm-rf-serial

RUN go get github.com/mikepb/go-serial

RUN go build -o main . 
CMD ["/go/src/github.com/drkiet/pi-alarm-rf-serial/main"]
