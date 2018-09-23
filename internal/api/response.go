package api

// Response format for api's
type Response struct {
	Success bool              `json:"success"`
	Data    map[string]string `json:"data"`
	Error   map[string]string `json:"error"`
}
