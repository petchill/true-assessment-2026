package error

type HttpResponseError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
