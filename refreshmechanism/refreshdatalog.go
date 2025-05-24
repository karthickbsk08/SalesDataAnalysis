package refreshmechanism

import (
	"salesdataanalysis/dbconnection"
	"salesdataanalysis/helpers"
)

var LogChannel chan<- RefreshLogActivity

// This method is used to initiate the channel for capture the logs
func InitiateRefreshDataCallLog() chan<- RefreshLogActivity {
	LogRefChannel := make(chan RefreshLogActivity, 100)
	go func() {
		count := 0
		for value := range LogRefChannel {
			count++
			RefreshDataLogCapture(count, value)
			//time.Sleep(time.Millisecond * 500)
		}
	}()
	return LogRefChannel
}

// // This method is used to insert the record by GORM
func RefreshDataLogCapture(count int, pRequest RefreshLogActivity) {
	pRequest.pDebug.Log(helpers.Statement, "RefreshDataLogCapture (+)", count, dbconnection.GPostgres.Stats())
	pRequest.pDebug.Log(helpers.Details, pRequest, "pRequest")
	// Insert the record into the database using GORM
	// Replace the <table_name>
	lResult := dbconnection.GRMPostgres.Table("refresh_log_activity").Create(&pRequest)
	if lResult.Error != nil {
		pRequest.pDebug.Log(helpers.Elog, "LogCapture Error", lResult.Error.Error())
	}

	pRequest.pDebug.Log(helpers.Statement, "RefreshDataLogCapture (-)")

}
