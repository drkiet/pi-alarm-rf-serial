package main 


import (
	"gopkg.in/mgo.v2"
    // "gopkg.in/mgo.v2/bson"
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
func getAllAlarmUnits() {
	connectMongoDb()
	defer dbClient.Close()

	pialarmDb := dbClient.DB(mongoDbName)
	alarmunits := pialarmDb.C(ALARM_UNITS_COLLECTION)
	count, _ := alarmunits.Count()

	LogMsg(fmt.Sprintf("# records: %d", count))

	iter := alarmunits.Find(nil).Iter()

	var event Event
	for iter.Next(&event) {
		fmt.Println("ID: ", event.ID, "; type: ", event.Type, "; Time: ", event.Time)
	}

}