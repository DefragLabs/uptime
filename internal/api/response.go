package api

// Response format for api's
type Response struct {
	Success bool              `json:"success"`
	Data    map[string]string `json:"data"`
	Error   map[string]string `json:"error"`
}

// StructResponse format for api's returning structs.
type StructResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
	Error   map[string]string      `json:"error"`
}
