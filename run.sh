export PI_ALARM_SERVER_ENDPOINT=$1
export PI_ALARM_RUNNING_MODE=$2

cd src
go install pi-alarm-rf-serial.go
cd ..
$GOBIN/pi-alarm-rf-serial

