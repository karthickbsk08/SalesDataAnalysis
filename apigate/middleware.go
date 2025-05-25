package apigate

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"salesdataanalysis/helpers"
	"strings"
	"time"
)

/* Step-by-Step Middleware Process
1. Receive the incoming request
2. Log request details (method, URL, headers, IP, timestamp) &&&&&&&&&&
3. Parse request body (JSON, form-data, etc.) &&&&&&&&&&
4. Authenticate the user (check token, session, API key)
5. Authorize the user (check permissions for the requested action)
6. Validate input data (query params, headers, body fields)
7. Apply rate limiting to prevent abuse &&&&&&&&&&
8. Handle CORS (set Access-Control-Allow-* headers) &&&&&&&&&
9. Add security headers (e.g., Content-Security-Policy, X-Frame-Options)
10. Pass request to the next middleware or route handler &&&&&&&&&&&
11. (Optional) Modify or compress the response &&&&&&&&
12. Send the final response to the client  &&&&&&&&&&& */

func RequestMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		lDebug := new(helpers.HelperStruct)
		lDebug.Init()
		lDebug.Log(helpers.Statement, "RequestMiddleWare(+)")

		var lLogDetails ApiLogCapture
		var capturewriter = ResponseCaptureWriter{ResponseWriter: w}
		lrequestTime := time.Now().Unix()

		// 1. Receive the incoming request
		// 2. Parse request body (JSON, form-data, etc.)
		lrequestDetails := GetRequestorDetail(r)
		// log.Println("lrequestDetails : ", lrequestDetails)
		r.Body = io.NopCloser(strings.NewReader(lrequestDetails.Body))

		//Handle CORS (set Access-Control-Allow-* headers)
		(w).Header().Set("Access-Control-Allow-Origin", "*")
		(w).Header().Set("Access-Control-Allow-Credentials", "true")
		(w).Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == http.MethodOptions {
			capturewriter.WriteHeader(http.StatusOK)
			capturewriter.Status()
			return
		}

		ctx := context.WithValue(r.Context(), helpers.RequestIDKey, lDebug.Sid)
		r = r.WithContext(ctx)

		next.ServeHTTP(&capturewriter, r)

		//respone time prep
		lResponseTime := time.Now().Unix()
		lResonseDateTime := time.Now()

		//2. Log request details (method, URL, headers, IP, timestamp)
		lLogDetails.RequestId = lDebug.Sid
		lLogDetails.Endpoint = lrequestDetails.EndPoint
		lLogDetails.ForwardedIP = lrequestDetails.ForwardedIP
		lLogDetails.Header = lrequestDetails.Header
		lLogDetails.Host = lrequestDetails.Host
		lLogDetails.Method = lrequestDetails.Method
		lLogDetails.Path = lrequestDetails.Path
		lLogDetails.RealIP = lrequestDetails.RealIP
		lLogDetails.RemoteAddr = lrequestDetails.RemoteAddr
		lLogDetails.ReqBody = lrequestDetails.Body
		lLogDetails.RequestUnixTime = lrequestTime
		lLogDetails.ResponseUnixTime = lResponseTime
		lLogDetails.RespDateTime = lResonseDateTime
		lLogDetails.RespBody = string(capturewriter.Body())
		lLogDetails.ResponseStatus = int64(capturewriter.Status())
		lLogDetails.PDebug = lDebug

		lData, _ := json.Marshal(lLogDetails)
		log.Println("DATA : ", string(lData))

		ApiCallLogChannel <- lLogDetails
		lDebug.Log(helpers.Statement, "RequestMiddleWare(-)")

	})

}
