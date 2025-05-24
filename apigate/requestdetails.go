package apigate

import (
	"io"
	"net/http"
)

// --------------------------------------------------------------------
// get request header details
// --------------------------------------------------------------------
func GetHeaderDetails(r *http.Request) string {
	value1 := ""
	for name, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			value1 = value1 + " " + name + "-" + value
		}
	}
	return value1
}

func GetRequestorDetail(r *http.Request) RequestorDetails {
	var reqDtl RequestorDetails
	// Prefer X-Original-Forwarded-For, then X-Forwarded-For, then X-Real-IP, then RemoteAddr
	if originalForwarded := r.Header.Get("X-Original-Forwarded-For"); originalForwarded != "" {
		reqDtl.ForwardedIP = originalForwarded
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		reqDtl.ForwardedIP = forwarded
	} else if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		reqDtl.ForwardedIP = realIP
	} else {
		reqDtl.ForwardedIP = r.RemoteAddr
	}
	reqDtl.RealIP = r.Header.Get("Referer")
	reqDtl.Method = r.Method
	reqDtl.Path = r.URL.Path + "?" + r.URL.RawQuery
	reqDtl.Host = r.Host
	reqDtl.RemoteAddr = r.RemoteAddr
	reqDtl.EndPoint = r.URL.Path
	reqDtl.RequestType = r.Header.Get("Content-Type")
	reqDtl.Header = GetHeaderDetails(r)
	if r.Body != nil {
		lBodyByte, _ := io.ReadAll(r.Body)
		reqDtl.Body = string(lBodyByte)
	}
	return reqDtl
}
