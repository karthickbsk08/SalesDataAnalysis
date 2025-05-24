package handlers

import "salesdataanalysis/helpers"

func LoadData(pDebug *helpers.HelperStruct) (int64, error) {
	pDebug.StartFunc()
	var lErr error
	var lRecordCnt int64

	lRecordCnt, lErr = LoadCSVInStagingTable(pDebug)
	if lErr != nil {
		pDebug.Log(helpers.Elog, lErr.Error())
		return lRecordCnt, helpers.ErrReturn(lErr)
	}
	lErr = InsertIntoAllTablesUniquely(pDebug)
	if lErr != nil {
		pDebug.Log(helpers.Elog, lErr.Error())
		return lRecordCnt, helpers.ErrReturn(lErr)

	}
	pDebug.ExitFunc()
	return lRecordCnt, nil
}
