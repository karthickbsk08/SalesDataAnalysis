package refreshmechanism

import (
	"encoding/json"
	"net/http"
	"salesdataanalysis/apigate"
	"salesdataanalysis/helpers"
	"time"
)

func RefreshDataOnDemand(w http.ResponseWriter, r *http.Request) {
	lDebug := new(helpers.HelperStruct)
	lDebug.SetUid(r)
	lDebug.StartFunc()

	var lRespRec apigate.ResponseStruct
	lRespRec.Status = "S"
	lRespRec.Msg = ""

	if r.Method == http.MethodGet {
		startTime := time.Now()

		TotalRecCnt, lErr := LoadData(lDebug)
		if lErr != nil {
			lDebug.Log(helpers.Elog, lErr.Error())
			lRespRec.Msg = lErr.Error()
			lRespRec.Status = "E"
		} else {
			// Build log details
			endTime := time.Now()
			duration := int(endTime.Sub(startTime).Seconds())

			logDetails := RefreshLogActivity{
				RequestID:            lDebug.Sid,
				CreatedBy:            "ondemandapi",
				CreatedAt:            startTime,
				UpdatedAt:            endTime,
				UpdatedBy:            "ondemandapi",
				DurationSeconds:      duration,
				ErrorMessage:         lRespRec.Msg,
				Status:               lRespRec.Status,
				RefreshType:          "Autobot",
				TotalRecordsAffected: TotalRecCnt,
				pDebug:               lDebug,
			}
			// Push to logging channel
			LogChannel <- logDetails

		}

		lErr = json.NewEncoder(w).Encode(lRespRec)
		if lErr != nil {
			lDebug.Log(helpers.Elog, lErr.Error())
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}
