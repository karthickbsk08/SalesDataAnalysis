package dbconnection

import (
	"database/sql"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var GDB AllDatabases

const (
	// DBTypeMariaDB is the database type for MariaDB
	DBTypeMariaDB       = "MARIADB"
	DBTypePostGresSQlDB = "POSTGRESSQLDB"
	// DBTypeMySQL is the database type for MySQL
)

func DBConnect(pDBFlag string) (*sql.DB, *gorm.DB, error) {
	// Load environment variables from .env file
	GDB.Init()

	// if pDBFlag == GDB.MariaDB.DBflag {
	// 	return ConnectToDB(GDB.MariaDB, GDB.MariaDB.DBType)
	// }
	if pDBFlag == GDB.PostGresSQL.DBflag {
		return ConnectToDB(GDB.PostGresSQL, GDB.PostGresSQL.DBType)
	}
	return nil, nil, fmt.Errorf("unsupported database type: %s", pDBFlag)
}
func ConnectToDB(dbConfig *DBConfig, DbType string) (*sql.DB, *gorm.DB, error) {

	var lDialector gorm.Dialector
	var lErr error
	var lGORMDb *gorm.DB
	var lSqlDb *sql.DB

	var dsn string
	switch DbType {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbConfig.User, dbConfig.Password, dbConfig.Server, dbConfig.Port, dbConfig.Database)
		lDialector = mysql.New(mysql.Config{
			DSN: dsn})
	case "mssql":
		// Driver={ODBC Driver 17 for SQL Server};Server=your_server_name;Database=your_database_name;Uid=your_user_id;Pwd=your_password;
		dsn = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;timeout=2;", dbConfig.Server, dbConfig.User, dbConfig.Password, dbConfig.Port, dbConfig.Database)
		lDialector = sqlserver.New(sqlserver.Config{
			DSN: dsn})

	case "postgres":
		// Driver={PostgreSQL ODBC Driver};Server=your_server_name;Port=5432;Database=your_database_name;Uid=your_user_id;Pwd=your_password;
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbConfig.Server, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Database)
		lDialector = postgres.New(postgres.Config{
			DSN: dsn})
	default:
		return lSqlDb, lGORMDb, fmt.Errorf("unsupported database type: %s", DbType)
	}

	lGORMDb, lErr = gorm.Open(lDialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if lErr != nil {
		log.Println("Error @ Sql DB connection (LGFLDBC005) : ", lErr.Error())
		return lSqlDb, lGORMDb, fmt.Errorf(" Invalid DB Details")
	}
	lSqlDb, lErr = lGORMDb.DB()
	if lErr != nil {
		log.Println("Error @ Sql DB connection (LGFLDBC005) : ", lErr.Error())
		return lSqlDb, lGORMDb, fmt.Errorf(" Invalid DB Details")
	}

	return lSqlDb, lGORMDb, nil

}
