package main 


import (
	"gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
	"net/url"
	"fmt"
	"log"
)

var mongoDbHost, mongoDbName, mongoDbUsername, mongoDbPassword string
var dbClient *mgo.Session

const ALARM_UNITS_COLLECTION = "alarmunits"

func makeDbUrl() (dbUrl string){
	var mongoDbUrl url.URL
	mongoDbUrl.Scheme = "mongodb"
	mongoDbUrl.Host = mongoDbHost
	mongoDbUrl.User = url.UserPassword(mongoDbUsername, mongoDbPassword)
	mongoDbUrl.Path = mongoDbName
	return mongoDbUrl.String()
}

func connectMongoDb() {
	dbUrl := makeDbUrl()
	LogMsg("Connecting: " + dbUrl)

	var err error
	dbClient, err = mgo.Dial(dbUrl)
	if err != nil {
		fmt.Println("Error in open database")
		log.Println(err)
	}
}

// Get all alarm units from the collection
func getAllAlarmUnits() (alarmUnits []AlarmUnit) {
	connectMongoDb()
	defer dbClient.Close()

	pialarmDb := dbClient.DB(mongoDbName)
	dbAlarmUnits := pialarmDb.C(ALARM_UNITS_COLLECTION)
	count, _ := dbAlarmUnits.Count()

	LogMsg(fmt.Sprintf("# records: %d", count))

	dbAlarmUnits.Find(nil).All(&alarmUnits)

	return
}

func getAnAlarmUnit(macid string) (alarmUnit *AlarmUnit) {
	connectMongoDb()
	defer dbClient.Close()

	pialarmDb := dbClient.DB(mongoDbName)
	dbAlarmUnits := pialarmDb.C(ALARM_UNITS_COLLECTION)

	var all []AlarmUnit
	dbAlarmUnits.Find(bson.M{"macid": macid}).All(&all)

	alarmUnit = &all[0]
	return
}

func updateOrCreateIfNotExistAlarmUnit(id string, alarmUnit AlarmUnit) {
	return
}