# Introduction
The purpose of this application is to run as a UDP server that expects incoming messages from a PI that process incoming messages from RF sensors. In other words, the code that runs on the PI listens to the RF message and then forward exactly the messages to a UDP listener acting as a UDP Server. The application runs on the PI acts a UDP Client.

Microservice architecture - in other words, it is small, independently executed and stateless. A single program can run as:

- An RF service: 
	+ listens to an incoming RF data from the sensors (zones). 
	+ posts the event to another service via HTTP or UDP protocol.
	+ should runs on a Raspberry PI with a Wireless Base Station for Raspberry Pi from private eye pi located here `http://privateeypi.com`
	+ run example:
		Posting the sensor events to a UDP server
		```
		./run.sh RF_RECEIVER_TO_UDP 192.168.1.63:9999 rf-alarm.log
		```
		Posting the sensor events to a HTTP server
		```
		./run.sh RF_RECEIVER_TO_HTTP http://192.168.1.63:9090 rf-alarm.log
		```


- a service that listens (HTTP or UDP) to an incoming event from an alarm application running on a Raspberry PI, receives it and processes it.

- a service that listens (UDP) incoming event from an alarm application running on a Raspberry PI, receives it and forward to another endpoint via UDP or HTTP

- to run use this shell:
`
run.sh
`

## Docker Image:
- Built & pushed to `hub.docker.com`. You can get a copy there.
```
sudo docker build drtran/pi-alarm-rf-serial .
sudo docker tag pi-alarm-rf-serial:latest drtran/pi-alarm-rf-serial:latest
sudo docker login
sudo docker push drtran/pi-alarm-rf-serial:latest

```

- Imported an image into Openshift & run an application in Openshift using these commands:
```
oc import-image drtran/pi-alarm-rf-serial --from=drtran/pi-alarm-rf-serial:latest  --confirm=true --insecure=true
oc new-app -e PI_ALARM_SERVER_ENDPOINT=:8080 -e PI_ALARM_RUNNING_MODE=EVENT_HTTP_SERVER -e PI_ALARM_LOG_FILE_NAME=./alarm.log pi-alarm-rf-serial:latest
oc expose deploymentconfigs/pi-alarm-rf-serial --port=8080
oc expose services/pi-alarm-rf-serial

```
http://pi-alarm-rf-serial-rest-example.1d35.starter-us-east-1.openshiftapps.com/
http://pi-alarm-rf-serial-rest-example.1d35.starter-us-east-1.openshiftapps.com

## configuring serial port on Raspberry Pi 3
Use this:

sudo systemctl stop serial-getty@ttyAMA0.service
sudo systemctl disable serial-getty@ttyAMA0.service

sudo systemctl stop serial-getty@ttyS0.service
sudo systemctl disable serial-getty@ttyS0.service
sudo vim /boot/config.txt
add these:

```
dtoverlay=pi3-miniuart-bt
core_freq=250
enable_uart=1
```

make sure /boot/cmdline.txt contains this (A SINGLE LINE!!!):
```
dwc_otg.lpm_enable=0 console=tty1 root=PARTUUID=201422d0-02 rootfstype=ext4 elevator=deadline fsck.repair=yes rootwait
```