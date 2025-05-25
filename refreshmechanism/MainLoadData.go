package refreshmechanism

import (
	"fmt"
	"salesdataanalysis/dbconnection"
	"salesdataanalysis/helpers"
)

func LoadData(pDebug *helpers.HelperStruct) (int64, error) {
	pDebug.StartFunc()
	var lErr error
	var lRecordCnt int64

	lRecordCnt, lErr = LoadCSVInStagingTable(pDebug)
	if lErr != nil {
		pDebug.Log(helpers.Elog, lErr.Error())
		return lRecordCnt, helpers.ErrReturn(lErr)
	}
	if lRecordCnt > 0 {
		err := dbconnection.GRMPostgres.Exec(`
						TRUNCATE TABLE 
							categories, 
							products, 
							customers, 
							orders, 
							regions 
						RESTART IDENTITY CASCADE;`).Error

		if err != nil {
			return lRecordCnt, helpers.ErrReturn(fmt.Errorf(" Failed to truncate tables: %v", err))
		}

		lErr = InsertIntoAllTablesUniquely(pDebug)
		if lErr != nil {
			pDebug.Log(helpers.Elog, lErr.Error())
			return lRecordCnt, helpers.ErrReturn(lErr)

		}
	}
	pDebug.ExitFunc()
	return lRecordCnt, nil
}
