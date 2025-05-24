package refreshmechanism

import (
	"salesdataanalysis/errorcode"
	"salesdataanalysis/handlers"
	"salesdataanalysis/helpers"
	"salesdataanalysis/tomlutil"
	"time"
)

// StartDailyRollover kicks off a goroutine to run CompressAndStore at midnight daily
func StartDailyRollover(pDebug *helpers.HelperStruct) {
	go func() {
		var lErrMsg string
		var RefreshFrequency int

		lErr := tomlutil.DecodeTOMLWithTypeCheck("./toml/config.toml", &RefreshFrequency)
		if lErr != nil {
			pDebug.Log(helpers.Elog, "Error in DecodeTOMLWithTypeCheck : ", lErr.Error())
			lErrMsg = lErr.Error()
		}
		ticker := time.NewTicker(time.Duration(RefreshFrequency) * time.Minute)
		defer ticker.Stop()

		for {
			<-ticker.C
			errorcode.Mu.Lock()

			pDebug.Init()
			var lRequestLogStTime = time.Now()
			errorcode.Mu.Unlock()

			var lStatus string
			var lLogDetails RefreshLogActivity
			lRecordCnt, lErr := handlers.LoadData(pDebug)
			if lErr != nil { // RegisterError(err.Error())
				pDebug.Log(helpers.Elog, "Error in LoadData : ", lErr.Error())
				// fmt.Println("Error compressing/storing logs:", lErr)
				lErrMsg = lErr.Error()
			}

			if lErrMsg != "" {
				lStatus = "E"
			} else {
				lStatus = "S"
			}

			//2. Log request details (method, URL, headers, IP, timestamp)
			lLogDetails.RequestID = pDebug.Sid
			lLogDetails.pDebug = pDebug
			lLogDetails.CreatedBy = "Autobot"
			lLogDetails.CreatedAt = lRequestLogStTime
			lLogDetails.UpdatedAt = time.Now()
			lLogDetails.UpdatedBy = "Autobot"
			lLogDetails.DurationSeconds = int(lLogDetails.UpdatedAt.Sub(lRequestLogStTime))
			lLogDetails.ErrorMessage = lErrMsg
			lLogDetails.Status = lStatus
			lLogDetails.RefreshType = "Autobot"
			lLogDetails.TotalRecordsAffected = lRecordCnt

			LogChannel <- lLogDetails

		}
	}()
}
