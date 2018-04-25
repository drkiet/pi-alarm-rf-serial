#
# Required fields:
# examples: 
# RUNNING_MODE: 
# 	RF_RECEIVER_TO_UDP
# 	RF_RECEIVER_TO_HTTP
# 	EVENT_UDP_SERVER
# 	EVENT_HTTP_SERVER
# 	UDP_REPEATER
# 
# SERVER_ENDPOINT: 
#	192.168.1.63:9999 (udp)
#	192.168.1.63:9090 (http)
#
# FILE_NAME: ./alarm.log
#
export PI_ALARM_RUNNING_MODE=$1
export PI_ALARM_SERVER_ENDPOINT=$2
export PI_ALARM_LOG_FILE_NAME=$3

#
# Optional fields:
# example: 172.17.0.2:9999 (udp)
#
export PI_ALARM_REPEATER_ENDPOINT=$4

cd src
go install \
	pi_alarm_rf_serial.go \
	utils.go \
	http.go \
	udp.go \
	rf_tx_rx.go \
	sensor.go \
	event.go
cd ..
$GOBIN/pi_alarm_rf_serial

