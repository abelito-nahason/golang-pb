package genericresponse

type GenericResponse struct {
	Data            any    `json:"data"`
	ResponseMessage string `json:"responseMessage"`
	ResponseCode    uint8  `json:"responseCode"`
}

type GenericErrorResponse struct {
	Error           error  `json:"error"`
	ResponseMessage string `json:"responseMessage"`
	ResponseCode    int    `json:"responseCode"`
}
