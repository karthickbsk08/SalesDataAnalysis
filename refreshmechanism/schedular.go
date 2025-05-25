package refreshmechanism

import (
	"fmt"
	"salesdataanalysis/errorcode"
	"salesdataanalysis/helpers"
	"salesdataanalysis/tomlutil"
	"strconv"
	"time"
)

func StartDailyRollover(pDebug *helpers.HelperStruct) {
	pDebug.Log(helpers.Statement, "StartDailyRollover(+)")
	go func() {
		var lErrMsg string
		var refreshFrequency int

		// Load frequency from TOML
		lConfig := tomlutil.ReadTomlConfig("./toml/config.toml")
		refreshFrequency, _ = strconv.Atoi(fmt.Sprintf("%v", lConfig.(map[string]interface{})["RefreshPeriod"]))

		ticker := time.NewTicker(time.Duration(refreshFrequency) * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			pDebug.Log(helpers.Details, "Start with Current session : ", pDebug.Sid)
			startTime := time.Now()

			// Lock only for Init or shared state
			errorcode.Mu.Lock()
			pDebug.Init()
			requestID := pDebug.Sid // capture safely
			errorcode.Mu.Unlock()

			// Run the actual logic
			recordCount, lErr := LoadData(pDebug)
			if lErr != nil {
				lErrMsg = lErr.Error()
				pDebug.Log(helpers.Elog, "LoadData error: ", lErrMsg)
			} else {
				lErrMsg = ""
			}

			// Build log details
			endTime := time.Now()
			duration := int(endTime.Sub(startTime).Seconds())
			status := "S"
			if lErrMsg != "" {
				status = "E"
			}

			logDetails := RefreshLogActivity{
				RequestID:            requestID,
				pDebug:               pDebug,
				CreatedBy:            "Autobot",
				CreatedAt:            startTime,
				UpdatedAt:            endTime,
				UpdatedBy:            "Autobot",
				DurationSeconds:      duration,
				ErrorMessage:         lErrMsg,
				Status:               status,
				RefreshType:          "Autobot",
				TotalRecordsAffected: recordCount,
			}

			// Push to logging channel
			LogChannel <- logDetails
		}
	}()
}
