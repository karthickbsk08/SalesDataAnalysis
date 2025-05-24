package helpers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"runtime"
	"salesdataanalysis/tomlutil"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type ContextKey string

const RequestIDKey ContextKey = "requestID"

const (
	Elog      = 0
	Statement = 1
	Details   = 2
)

type HelperStruct struct {
	Sid       string
	Reference string
}

func (h *HelperStruct) Init() {
	KeyValue := ""
	session := uuid.NewV4()
	sessionSHA256 := session.String()
	KeyValue = strings.ReplaceAll(sessionSHA256, "-", "")
	h.Sid = KeyValue
}
func (h *HelperStruct) SetUid(r *http.Request) {
	requestID, ok := r.Context().Value(RequestIDKey).(string)
	if !ok {
		KeyValue := ""
		session := uuid.NewV4()
		sessionSHA256 := session.String()
		KeyValue = strings.ReplaceAll(sessionSHA256, "-", "")
		h.Sid = KeyValue
	} else {
		h.Sid = requestID
	}
}

func (h *HelperStruct) SetReference(pReference interface{}) {
	h.Reference = fmt.Sprintf("%v", pReference)
}

func (h *HelperStruct) RemoveReference() {
	h.Reference = ""
}

func ErrReturn(pErr error) error {
	lErr := ""
	lPc, lFile, lLine, _ := runtime.Caller(1)
	lFuncname := runtime.FuncForPC(lPc).Name()
	lStrArray := strings.Split(lFile, "/")
	lFilename := lStrArray[len(lStrArray)-2] + "/" + lStrArray[len(lStrArray)-1]

	if strings.Contains(pErr.Error(), " @@ ") && strings.Contains(pErr.Error(), " @@ ln ") {
		lErr = ""
	} else {

		lErr = lFilename + " @@ " + lFuncname + " @@ ln " + fmt.Sprintf("%d", lLine) + " @@ "
	}
	return errors.New(lErr + pErr.Error())
}

func (h *HelperStruct) ExitFunc() {
	lPc, lFilename, lLine, _ := runtime.Caller(1)
	lFuncname := runtime.FuncForPC(lPc).Name()
	funcName := strings.Split(lFuncname, ".")[1]
	h.Log(Statement, funcName, " (-)")
	log.Println("@@", h.Sid, "@@ (", h.Reference, ") @@", lFilename, "@@", lFuncname, "@@ ln", lLine, "@@", funcName, "(-)")

}
func (h *HelperStruct) StartFunc() {
	lPc, lFilename, lLine, _ := runtime.Caller(1)
	lFuncname := runtime.FuncForPC(lPc).Name()
	funcName := strings.Split(lFuncname, ".")[1]
	log.Println("@@", h.Sid, "@@ (", h.Reference, ") @@", lFilename, "@@", lFuncname, "@@ ln", lLine, "@@", funcName, "(+)")
}

func (h *HelperStruct) Log(pDebugLevel int, pMsg ...interface{}) {
	// find the mata data of printed value
	lPc, lFile, lLine, _ := runtime.Caller(1)
	lFuncname := runtime.FuncForPC(lPc).Name()

	//read the value from toml

	lConfigFile := tomlutil.ReadTomlConfig("./toml/debug.toml")
	lLevel := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["LogCategory"])
	intlevel, err := strconv.Atoi(lLevel)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	lBase64FileLevel := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["Base64Encode"])
	lBase64level, err := strconv.Atoi(lBase64FileLevel)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	lReference := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["LogReference"])

	//file name over writting
	lStrArray := strings.Split(lFile, "/")
	lFilename := lStrArray[len(lStrArray)-2] + "/" + lStrArray[len(lStrArray)-1]

	//find the Debug level
	if (pDebugLevel <= intlevel && intlevel != 0) || pDebugLevel == Elog {
		str := fmt.Sprintf("%v", pMsg)
		lFinal := ReplaceBase64String(str[1:len(str)-1], lBase64level)

		//check the sid will be set
		if strings.EqualFold(h.Sid, "") {
			log.Println("Set the id before <-- debug := new(helpers.HelperStruct) after debug.Init(reference value) -->")
		} else {

			//print the O/P based on reference value
			if lReference == "" || strings.EqualFold(lReference, h.Reference) {
				if h.Reference == "" || lReference != "" {
					h.Reference = lReference
				}

				if strings.Contains(lFinal, " @@ ") && strings.Contains(lFinal, "@@ ln") {
					log.Println("@@", h.Sid, "@@ (", h.Reference, ") @@", lFinal)
				} else {
					log.Println("@@", h.Sid, "@@ (", h.Reference, ") @@", lFilename, "@@", lFuncname, "@@ ln", lLine, "@@", lFinal)
				}
			}
		}

	}
}

func ErrPrint(err error) string {
	lStrArray := strings.Split(err.Error(), "@@")
	return lStrArray[len(lStrArray)-1]
}

// func ReplaceBase64String(pBase46String string, pFilterCategory int) string {
// 	if pFilterCategory == 1 || pBase46String == "" {
// 		return pBase46String
// 	}

// 	// Regular expression to match base64 encoded strings within double quotes
// 	base64Regex := regexp.MustCompile(`("[^"]*?":[ ]*")((?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)*)(")`)
// 	// Find all matches in the JSON string
// 	matches := base64Regex.FindAllStringSubmatch(pBase46String, -1)

// 	// Iterate over matches
// 	if len(matches) == 0 {
// 		return pBase46String
// 	}
// 	for _, match := range matches {
// 		// Append the base64 encoded string to the array
// 		decoded, err := base64.StdEncoding.DecodeString(match[2])
// 		if err != nil || len(decoded) == 0 {
// 			continue
// 		}
// 		if !strings.Contains(strings.ToLower(http.DetectContentType(decoded)), "text/plain") {

// 			pBase46String = strings.ReplaceAll(pBase46String, match[2], "base64 encoded file")
// 		}
// 	}
// 	return pBase46String
// }

func ReplaceString(pBase46String string, pThreshold int) string {

	// Regular expression to match base64 encoded strings within double quotes
	base64Regex := regexp.MustCompile(`("[^"]*?":[ ]*")((?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)*)(")`)
	// Find all matches in the JSON string
	matches := base64Regex.FindAllStringSubmatch(pBase46String, -1)

	// Iterate over matches
	if len(matches) == 0 {
		return pBase46String
	}
	for _, match := range matches {
		// Append the base64 encoded string to the array

		if len(match[2]) < pThreshold {
			continue
		}

		decoded, err := base64.StdEncoding.DecodeString(match[2])
		if err != nil || len(decoded) == 0 || strings.Contains(strings.ToLower(http.DetectContentType(decoded)), "text/plain") {
			continue
		}
		// if !strings.Contains(strings.ToLower(http.DetectContentType(decoded)), "text/plain") {
		pBase46String = strings.ReplaceAll(pBase46String, match[2], "base64 encoded file")
		// }
	}
	return pBase46String
}

func ReplaceBase64String(pBase64String string, pFilterCategory int) string {
	// Try to decode the string
	if pFilterCategory == 1 || pBase64String == "" {
		return pBase64String
	}

	pattern := `("[^"]*?":[ ]*")((?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)*)(")`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// lThreshold := tomlconfig.GtomlConfigLoader.GetValueString("debug", "Threshold")
	lConfigFile := tomlutil.ReadTomlConfig("./toml/debug.toml")
	lThreshold := fmt.Sprintf("%v", lConfigFile.(map[string]interface{})["Threshold"])

	lThresholdInt, lErr := strconv.Atoi(lThreshold)

	if lErr != nil {
		fmt.Println("Error:", lErr.Error())
		return lErr.Error()
	}

	// lThreshold := tomlconfig.GtomlConfigLoader.GetValueString("debug", "Threshold")

	if re.MatchString(pBase64String) {
		pBase64String = ReplaceString(pBase64String, lThresholdInt)
		return pBase64String
	}

	DecodeString, err := base64.StdEncoding.DecodeString(pBase64String)
	if err == nil {
		if re.MatchString(string(DecodeString)) {
			pBase64String = ReplaceString(string(DecodeString), lThresholdInt)
			return pBase64String
		} else if !strings.Contains(strings.ToLower(http.DetectContentType([]byte(DecodeString))), "text/plain") && len(DecodeString) > lThresholdInt {
			return "Base 64 Encode File"
		}
	}
	if !strings.Contains(strings.ToLower(http.DetectContentType([]byte(pBase64String))), "text/plain") && len(pBase64String) > lThresholdInt {
		return "Raw File"
	}
	return pBase64String
}
