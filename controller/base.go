package controller

type ApiJson struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func ApiResource(status int, objects interface{}, msg string) (apijson *ApiJson) {
	apijson = &ApiJson{Status: status, Data: objects, Msg: msg}
	return
}
