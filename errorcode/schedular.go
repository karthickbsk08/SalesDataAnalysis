package errorcode

import (
	"compress/gzip"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"salesdataanalysis/dbconnection"
	"strings"
	"time"
)

// StartDailyRollover kicks off a goroutine to run CompressAndStore at midnight daily
func StartDailyRollover() {
	go func() {
		ticker := time.NewTicker(2 * time.Minute)
		defer ticker.Stop()

		for {
			<-ticker.C
			err := CompressAndStore()
			if err != nil {
				RegisterError(err.Error())
				fmt.Println("Error compressing/storing logs:", err)
			}
		}
	}()
}

type DBLogLayout struct {
	LogDate             string `json:"log_date" gorm:"column:log_date"`
	CompressedLogBase64 string `json:"compressed_log_base64" gorm:"column:compressed_log_base64"`
}

// CompressAndStore compresses today’s logs, encodes base64, inserts into DB, clears registry
func CompressAndStore() error {
	Mu.Lock()
	defer Mu.Unlock()

	var buf strings.Builder
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Write([]string{"Error Code", "Package", "Method", "Description", "Line/File"})
	for _, e := range Registry {
		csvWriter.Write([]string{e.Code, e.Package, e.Method, e.Description, e.LineFile})
	}
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		RegisterError(err.Error())
		return err
	}

	// Compress CSV text
	compressedBytes, err := gzipCompress([]byte(buf.String()))
	if err != nil {
		RegisterError(err.Error())
		return err
	}

	encoded := base64.StdEncoding.EncodeToString(compressedBytes)

	// Store in DB, assuming today's date as key
	today := time.Now().Format("2006-01-02")

	var lLogInfo DBLogLayout
	lLogInfo.LogDate = today
	lLogInfo.CompressedLogBase64 = encoded

	lResult := dbconnection.GRMMaria.Table("daily_logs").Create(&lLogInfo)
	if lResult.Error != nil {
		RegisterError(lResult.Error.Error())
		return lResult.Error
	}

	// Clear registry and counters
	Registry = map[string]*ErrorEntry{}
	Counters = map[string]int{}
	return nil
}

func gzipCompress(data []byte) ([]byte, error) {
	var b strings.Builder
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(data); err != nil {
		RegisterError(err.Error())
		return nil, err
	}
	if err := gz.Close(); err != nil {
		RegisterError(err.Error())
		return nil, err
	}
	return []byte(b.String()), nil
}
