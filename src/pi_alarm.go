package main 

import (
	"log"
	"time"
	"strings"
	"fmt"
)

const RFBaseStationSerial = "/dev/ttyAMA0" // best to get it from environment var.
const SerialPortSpeed int = 9600
const MaxZones = 100

const (
	ViaEmail = "Via Email"
	ViaText = "Via Text"
	ViaPhone = "Via Phone"
)

type Zone Sensor

type PiAlarm struct {
	Id          string    `json:"id"`
	MacId       string    `json:"macid"`
	Owner       string    `json:"owner"`
	Email       string    `json:"email"`
	Cell        string    `json:"cell"`
	NotifyVia   string    `json:"notifyvia"`
	Zones       map[string]*Zone `json:"zones"`
	CurState    string    `json:"curstate"`
	WantedState string	  `json:"wantedstate"`
	Updated     time.Time `json:"updatedby"`
}

var piAlarm PiAlarm

func getSystemInfo() (sysInfo *PiAlarm) {
	sysInfo = new(PiAlarm)
	sysInfo.Id          = "1234"
	sysInfo.MacId       = piAlarm.MacId
	sysInfo.Owner       = piAlarm.Owner
	sysInfo.Email       = piAlarm.Email
	sysInfo.Cell        = piAlarm.Cell
	sysInfo.CurState    = piAlarm.CurState
	sysInfo.WantedState = piAlarm.WantedState
	sysInfo.Updated     = piAlarm.Updated
	sysInfo.Zones       = piAlarm.Zones
	return sysInfo
}

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
	zone.Id = strings.Trim(zoneTokens[0], " ")
	zone.ZoneName = strings.Trim(zoneTokens[1], " ")
	zone.State = "UNK"
	piAlarm.Zones[zone.Id] = zone
}

func updateZoneName(id string, zoneName string) {
	piAlarm.Zones[id].ZoneName = zoneName
	piAlarm.Updated = time.Now()
}

func removeZone(id string) {
	delete(piAlarm.Zones, id)
	piAlarm.Updated = time.Now()
}

func updateZone(sensor *Sensor) {
	var zone Zone = Zone(*sensor)
	piAlarm.Zones[zone.Id] = &zone
	piAlarm.Updated = time.Now()
}

func getZoneById(id string) (zone *Zone) {
	if piAlarm.Zones[id] == nil {
		return nil
	}

	zone = new (Zone)
	copyZone(zone, piAlarm.Zones[id])
	return
}

func copyZone(zone1 *Zone, zone2 *Zone) {
	zone1.Id = zone2.Id
	zone1.Type = zone2.Type
	zone1.ZoneName = zone2.ZoneName
	zone1.State = zone2.State
	zone1.Subunit = zone2.Subunit
	zone1.Battery = zone2.Battery
	zone1.Data = zone2.Data
	zone1.Updated = zone2.Updated
}

func getCurState() (curState string) {
	return piAlarm.CurState
}

func setCurState(curState string) {
	piAlarm.CurState = curState
	piAlarm.Updated = time.Now()
}

func getWantedState() (wantedState string) {
	return piAlarm.WantedState
}

func setWantedState(wantedState string) {
	piAlarm.WantedState = wantedState
	piAlarm.Updated = time.Now()
}

func setUpdated(updated time.Time) {
	piAlarm.Updated = updated
}

func getLastPiAlarmUpdated() (lastUpdated time.Time) {
	return piAlarm.Updated
}

func getZones() (zones map[string]*Zone) {
	return piAlarm.Zones
}

func getFormattedZoneStates() (zonesState []string) {
	zonesState = make ([]string, len(piAlarm.Zones), len(piAlarm.Zones))
	var i int = 0
	for _, zone := range piAlarm.Zones {
		zonesState[i] = fmt.Sprintf("%s.%s : %s\n", zone.Id, zone.ZoneName, 
					  				zone.State)
		i++
	}

	return
}

const (
	ARMED = "ARMED"
	DISARMED = "DISARMED"
	PERIMETERED = "PERIMETERED"
	ALARMED = "ALARMED"
	FAULT = "FAULT"
	NOFAULT = "NOFAULT"
)

// Initializing the Pi alarm before operation.
func piAlarmInit() {
	piAlarm.MacId = getMacAddr()
	piAlarm.Zones = make(map[string]*Zone)
	
	piAlarm.Updated = time.Now()
	loadPiAlarmCfg()

	if runsOnPi() {
		rfInit(RFBaseStationSerial, SerialPortSpeed)
	} else {
		udpInit(getUdpEndpoint())
	}

	emailInit()
	setWantedState("DISARMED")
	setCurState("NOFAULT")

}

// Managing Pi alarms with RF base station and sensors.
func managePiAlarm() {
	log.Println("**** Alarm Manager ****")
	piAlarmInit()
	makeEventLog()

	sensorCh := make(chan Sensor)
	if runsOnPi() {
		go rfReceiver(sensorCh)
	} else {
		go udpReceiver(sensorCh)
	}

	go healthMonitor()
	go httpServer(getServerEndpoint())

	for {
       sensor := <- sensorCh
       updateZone(&sensor)
       log.Println("**** managePiAlarm: ", sensor)
	}
}

// Print the content of PiAlarm object
func printPiAlarm() {

	piAlarmInfo := "\n*** Pi Alarm ***"
	piAlarmInfo += "\n      MacId: " + piAlarm.MacId
	piAlarmInfo += "\n      Owner: " + piAlarm.Owner
	piAlarmInfo += "\n      Email: " + piAlarm.Email
	piAlarmInfo += "\n       Cell: " + piAlarm.Cell
	piAlarmInfo += "\n  NotifyVia: " + piAlarm.NotifyVia
	piAlarmInfo += "\n   CurState: " + piAlarm.CurState
	piAlarmInfo += "\nWantedState: " + piAlarm.WantedState
	piAlarmInfo += "\n    Updated: " + piAlarm.Updated.String()

	fmt.Println(piAlarmInfo)
	printZones(piAlarm.Zones)
}


func printZones(zones map[string]*Zone) {
	zonesInfo := fmt.Sprintf("\n*** Zones(%d) ***", len(zones))
	for _, zone := range zones {
		zonesInfo += printZone(*zone)
	}
	fmt.Println(zonesInfo)
}

func printZone(zone Zone) (zoneInfo string) {
	zoneInfo = ""
	zoneInfo += "\n\n     Id: " + zone.Id
	zoneInfo += "\nZoneName: " + zone.ZoneName
	zoneInfo += "\n    Type: " + zone.Type
	zoneInfo += "\n   State: " + zone.State
	zoneInfo += "\n Subunit: " + zone.Subunit
	zoneInfo += "\n Battery: " + zone.Battery
	zoneInfo += "\n    Data: " + zone.Data
	return
}

func lookupZoneName(id string) (zoneName string) {
	zone := piAlarm.Zones[id]
	
	if zone != nil {
		zoneName = zone.ZoneName
	} else {
		zoneName = "NONAME"
	}

	return
}

func soundAlarm() {
	fmt.Println("*** SOUND ALARMED ON ****")
}

func soundAlarmOff() {
	fmt.Println("*** SOUND ALARMED OFF ****")
}
