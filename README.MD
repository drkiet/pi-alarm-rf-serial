# Introduction
The purpose of this application is to run as a UDP server that expects incoming messages from a PI that process incoming messages from RF sensors. In other words, the code that runs on the PI listens to the RF message and then forward exactly the messages to a UDP listener acting as a UDP Server. The application runs on the PI acts a UDP Client.

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
