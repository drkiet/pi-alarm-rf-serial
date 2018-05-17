#
# Required parameters: There are three (3) required parameters are passed
# 	into the program via OS environment variables. They are provided to 
# 	this script in sequence.
#
# RUN_MODE: This pi_alarm_rf_serial program can run as one of these modes:
#
# 	1) RF_RECEIVER_TO_UDP
# 	2) RF_RECEIVER_TO_HTTP
# 	3) EVENT_UDP_SERVER
# 	4) EVENT_HTTP_SERVER
# 	5) UDP_REPEATER
#   6) UDP_HTTP_REPEATER
# 
# SERVER_ENDPOINT:
#	192.168.1.63:9999 (udp)
#	http://192.168.1.63:9090 (http)
#
# FILE_NAME: alarm.log
#
export RUN_MODE=$1
export PI_ALARM_SERVER_ENDPOINT=$2
export PI_ALARM_LOG_FILE_NAME=$3

#
# Optional parameter: the repeater endpoint is required when the program runs 
# as a UDP_REPEATER or UDP_HTTP_REPEATER mode.
#
# PI_ALARM_REPEATER_ENDPOINT: 
# 	192.168.1.63:9999 (udp)
# 	http://192.168.1.63:9090
#
export PI_ALARM_REPEATER_ENDPOINT=$4

export PI_ALARM_MONGODB_HOST=172.17.0.6:27017
export PI_ALARM_MONGODB_NAME=pialarmdb
export PI_ALARM_MONGODB_USERNAME=pialarmuser
export PI_ALARM_MONGODB_PASSWORD=Password12341234

cd src
go install alarm_main.go \
	email.go \
	ethernet.go \
	event.go \
	http.go \
	log_util.go \
	monitor.go \
	pi_alarm.go \
	pi_alarm_cfg.go \
	pi_network.go \
	rf_base.go \
	sensor.go 

cd ..
$GOBIN/alarm_main

