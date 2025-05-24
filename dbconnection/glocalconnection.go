package dbconnection

import (
	"database/sql"
	"salesdataanalysis/helpers"

	"gorm.io/gorm"
)

var (
	GMaria      *sql.DB
	GRMMaria    *gorm.DB
	GPostgres   *sql.DB
	GRMPostgres *gorm.DB
	lErr        error
)

func OpenPoolConnections(pDebug *helpers.HelperStruct) error {
	pDebug.StartFunc()

	// GMaria, GRMMaria, lErr = DBConnect(DBTypeMariaDB)
	// if lErr != nil {
	// 	pDebug.Log(helpers.Elog, "Error connecting to MariaDB:"+lErr.Error())
	// 	return helpers.ErrReturn(lErr)
	// }

	GPostgres, GRMPostgres, lErr = DBConnect(DBTypePostGresSQlDB)
	if lErr != nil {
		pDebug.Log(helpers.Elog, "Error connecting to PostGresSQlDB:"+lErr.Error())
		return helpers.ErrReturn(lErr)
	}
	pDebug.ExitFunc()
	return nil
}

func ClosePoolConnections(pDebug *helpers.HelperStruct) {
	pDebug.StartFunc()

	if GPostgres != nil {
		lErr := GPostgres.Close()
		if lErr != nil {
			pDebug.Log(helpers.Elog, "Error connecting to MariaDB:"+lErr.Error())
		}
	}
	pDebug.ExitFunc()
}
