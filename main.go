package main

import (
	"log"
	"net/http"
	"os"
	"salesdataanalysis/apigate"
	"salesdataanalysis/dbconnection"
	"salesdataanalysis/errorcode"
	"salesdataanalysis/handlers"
	"salesdataanalysis/helpers"
	"salesdataanalysis/refreshmechanism"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	lDebug := new(helpers.HelperStruct)
	lDebug.Init()

	lDebug.StartFunc()

	//Creations log file
	lFile, lErr := os.OpenFile("./log/logfile"+time.Now().Format("02012006.15.04.05.000000000")+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if lErr != nil {
		log.Fatalf("error opening file:  %v", lErr)
	}
	defer lFile.Close()

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(lFile)

	if err := dbconnection.OpenPoolConnections(lDebug); err != nil {
		log.Fatalf("failed to open pool connections: %v", err)
		dbconnection.ClosePoolConnections(lDebug)
		return
	}
	defer dbconnection.ClosePoolConnections(lDebug)

	// go errorcode.WriteintoString()

	lProceedValue, lAcceptValue := apigate.AssignRateLimitValue()

	lRouter := mux.NewRouter()

	lRouter.HandleFunc("/fetchlogascsv", errorcode.FetchLogCSV)
	lRouter.Handle("/customeranalysis", apigate.RateLimiter(handlers.ProvideCustomerAnalysis, lProceedValue, lAcceptValue))
	lRouter.Handle("/revenuereport", apigate.RateLimiter(handlers.RevenueReportAPI, lProceedValue, lAcceptValue))
	lRouter.Handle("/refreshDataOnDemand", apigate.RateLimiter(refreshmechanism.RefreshDataOnDemand, lProceedValue, lAcceptValue))
	lHandler := apigate.RequestMiddleWare(lRouter)

	// lRouter.HandleFunc("/uploadsalesdata")

	refreshmechanism.StartDailyRollover(lDebug)

	//Initiate Queue to process API Incoming and outgoing Log to this service
	apigate.ApiCallLogChannel = apigate.InitiateApiCallLog()
	refreshmechanism.LogChannel = refreshmechanism.InitiateRefreshDataCallLog()

	lSrv := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      lHandler,
		Addr:         ":19998",
	}
	lDebug.ExitFunc()
	lSrv.ListenAndServe()
}
