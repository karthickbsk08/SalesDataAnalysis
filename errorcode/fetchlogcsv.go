package errorcode

import (
	"compress/gzip"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"salesdataanalysis/dbconnection"
	"strings"
	"time"
)

// FetchLogCSV handles HTTP requests for CSV logs, date as query param "date" (yyyy-mm-dd)
// if no date or today, serve from memory; else fetch from DB, decode and serve

func FetchLogCSV(w http.ResponseWriter, r *http.Request) {

	// Uncomment the following lines if you need to handle these common headers.
	// Set up CROS credentails
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")

	// Set up the Header API for common
	// If you need additional header you need to config here itself
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	requestedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		RegisterError(err.Error())
		http.Error(w, "invalid date format, use yyyy-mm-dd", http.StatusBadRequest)
		return
	}

	today := time.Now().Format("2006-01-02")

	filename := fmt.Sprintf("logfile%s.%s.csv",
		requestedDate.Format("02012006"),
		time.Now().Format("15.04.05.000000000"))

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "text/csv")

	if dateStr == today {
		// Serve from in-memory
		err = ExportCSVToWriter(w)
		if err != nil {
			RegisterError("failed to export logs" + err.Error())
			http.Error(w, "failed to export logs", 500)
		}
		return
	}

	base64data, err := getErrorLogFromDB(dateStr)
	if err != nil {
		RegisterError(err.Error())
		http.Error(w, "decode error", 500)
		return
	}

	compressedBytes, err := base64.StdEncoding.DecodeString(base64data)
	if err != nil {
		RegisterError("decode error" + err.Error())
		http.Error(w, "decode error", 500)
		return
	}

	// Decompress and write to response
	err = gzipDecompressToWriter(compressedBytes, w)
	if err != nil {
		RegisterError("decompress error" + err.Error())
		http.Error(w, "decompress error", 500)
		return
	}
}

func gzipDecompressToWriter(compressed []byte, w io.Writer) error {
	r := strings.NewReader(string(compressed))
	gz, err := gzip.NewReader(r)
	if err != nil {
		RegisterError(err.Error())
		return err
	}
	defer gz.Close()

	_, err = io.Copy(w, gz)
	return err
}

// ExportCSVToWriter writes current in-memory registry as CSV to an io.Writer
func ExportCSVToWriter(w io.Writer) error {
	Mu.Lock()
	defer Mu.Unlock()

	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write([]string{"Error Code", "Package", "Method", "Description", "Line/File"})
	for _, e := range Registry {
		if err := writer.Write([]string{e.Code, e.Package, e.Method, e.Description, e.LineFile}); err != nil {
			RegisterError(err.Error())
			return err
		}
	}
	return nil
}

func getErrorLogFromDB(pDate string) (string, error) {
	// Fetch from DB compressed base64 data

	var base64data string
	lResult := dbconnection.GRMPostgres.Table("daily_logs").Select("compressed_log_base64").Where("log_date = ?", pDate).Scan(&base64data)
	if lResult.Error != nil {
		RegisterError(lResult.Error.Error())
		log.Println("Error fetching log from DB:", lResult.Error)
		return base64data, lResult.Error
	}
	log.Println("base64data : ", base64data)
	return base64data, lResult.Error
}
