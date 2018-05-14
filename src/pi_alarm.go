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
	Zones                 map[string]*Zone
	CurState, WantedState string	
	Updated               time.Time
}

var piAlarm PiAlarm

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
	zone := new (Zone)
	zone.SensorId = strings.Trim(zoneTokens[0], " ")
	zone.ZoneName = strings.Trim(zoneTokens[1], " ")
	piAlarm.Zones[zone.SensorId] = zone
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
func piAlarmInit() {
	piAlarm.MacId = getMacAddr()
	piAlarm.Zones = make(map[string]*Zone)
	
	piAlarm.Updated = time.Now()
	loadPiAlarmCfg()
	rfInit(RFBaseStationSerial, SerialPortSpeed)
	emailInit()

}

// Managing Pi alarms with RF base station and sensors.
func managePiAlarm() {
	log.Println("**** Alarm Manager ****")
	piAlarmInit()	

	sensorCh := make(chan Sensor)
	go rfReceiver(sensorCh)
	go healthMonitor()

	for {
       sensor := <- sensorCh
       log.Println("**** managePiAlarm: ", sensor)
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
			return
		}
	}
	zoneName = fmt.Sprintf("*** Unknown zone %s ***", sensorId)
	return
}
