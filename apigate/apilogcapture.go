package apigate

import (
	"log"
	"salesdataanalysis/dbconnection"
	"salesdataanalysis/helpers"
)

var ApiCallLogChannel chan<- ApiLogCapture

// This method is used to initiate the channel for capture the logs
func InitiateApiCallLog() chan<- ApiLogCapture {
	log.Println("InitiateApiCallLog(+)")
	LogRespChannel := make(chan ApiLogCapture, 100)
	go func() {
		count := 0
		for value := range LogRespChannel {
			count++
			ApiCallLogCapture(count, value)
			//time.Sleep(time.Millisecond * 500)
		}
	}()
	return LogRespChannel
}

// // This method is used to insert the record by GORM
func ApiCallLogCapture(count int, pRequest ApiLogCapture) {
	pRequest.PDebug.Log(helpers.Statement, "ApiCallLogCapture (+)", count, dbconnection.GPostgres.Stats())
	pRequest.PDebug.Log(helpers.Details, pRequest, "pRequest")
	// Insert the record into the database using GORM
	// Replace the <table_name>
	lResult := dbconnection.GRMPostgres.Table("api_log_capture").Create(&pRequest)
	if lResult.Error != nil {
		pRequest.PDebug.Log(helpers.Elog, "LogCapture Error", lResult.Error.Error())
	}

	pRequest.PDebug.Log(helpers.Statement, "ApiCallLogCapture (-)")

}
