package helpers

import (
	"encoding/json"
	"log"
)

const (
	lSuccessCode = "S"
	lErrorCode   = "E"
)

type Error_Response struct {
	Status    string `json:"status"`
	ErrorCode string `json:"statusCode"`
	ErrMsg    string `json:"errMsg"`
}

type Msg_Response struct {
	Status      string `json:"status"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func GetMsg_String(Msg_Title string, Msg_Description string) string {

	var Msg_Res Msg_Response

	Msg_Res.Status = lSuccessCode
	Msg_Res.Title = Msg_Title
	Msg_Res.Description = Msg_Description

	result, err := json.Marshal(Msg_Res)

	if err != nil {
		log.Println(err)
	}

	return string(result)

}

func GetErrorString(Err_Title string, Err_Description string) string {

	var Err_Response Error_Response

	Err_Response.Status = lErrorCode
	Err_Response.ErrorCode = Err_Title
	Err_Response.ErrMsg = Err_Title + "/" + Err_Description

	lResult, err := json.Marshal(Err_Response)

	if err != nil {
		log.Println(err)
	}

	return string(lResult)

}
