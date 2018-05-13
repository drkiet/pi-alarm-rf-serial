package main 

import (
	"log"
	"time"
	"strings"
	"fmt"
)

const RFBaseStationSerial = "/dev/ttyAMA0" // best to get it from environment var.
const SerialPortSpeed int = 9600
const MaxZones = 20

const (
	ViaEmail = "Via Email"
	ViaText = "Via Text"
	ViaPhone = "Via Phone"
)

type Zone Sensor

type PiAlarm struct {
	MacId, Owner, Email, 
	Cell, NotifyVia       string
	Zones                 []Zone
	CurState, WantedState string	
	Updated               time.Time
}

var piAlarm PiAlarm
var zoneIndex int = 0

func setOwner(owner string) {
	piAlarm.Owner = owner
}

func setEmail(email string) {
	piAlarm.Email = email
}

func setCell(cell string) {
	piAlarm.Cell = cell
}

func setNotifyVia(notifyVia string) {
	piAlarm.NotifyVia = notifyVia
}

func addZone (zoneCfg string) {
	zoneTokens := strings.Split(zoneCfg, "=")
	piAlarm.Zones = piAlarm.Zones[:zoneIndex+1]
	piAlarm.Zones [zoneIndex].SensorId = strings.Trim(zoneTokens[0], " ")
	piAlarm.Zones [zoneIndex].ZoneName = strings.Trim(zoneTokens[1], " ")
	zoneIndex++
}

func setCurState(curState string) {
	piAlarm.CurState = curState
}

func setWantedState(wantedState string) {
	piAlarm.WantedState = wantedState
}

func setUpdated(updated time.Time) {
	piAlarm.Updated = updated
}

// Initializing the Pi alarm before operation.
func piAlarmInitialize() {
	piAlarm.MacId = getMacAddr()
	piAlarm.Zones = make([]Zone, 0, MaxZones)
	
	piAlarm.Updated = time.Now()
	loadPiAlarmCfg()
	rfInitialize(RFBaseStationSerial, SerialPortSpeed)
}

// Managing Pi alarms with RF base station and sensors.
func managePiAlarm() {
	log.Println("Managing PI Alarm System")
	piAlarmInitialize()	

	for {
        data := rfReceive()
        log.Println("managePiAlarm: ", data)
        sensor := makeSensorEvent(data)
        log.Println(sensor)
	}
}

// Print the content of PiAlarm object
func printPiAlarm() {
	fmt.Println(piAlarm)
}

func lookupZoneName(sensorId string) (zoneName string) {
	for _, zone := range piAlarm.Zones {
		if sensorId == zone.SensorId {
			zoneName = zone.ZoneName
		}
	}
	zoneName = fmt.Sprintf("*** Unknown zone %s ***", sensorId)
	return
}