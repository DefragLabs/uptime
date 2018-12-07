package utils

const (
	// StatusUp service status up
	StatusUp = "UP"

	// StatusDown service status down
	StatusDown = "DOWN"
)

// GetServiceStatus returns StatusUp or StatusDown depending
// on the response status code.
func GetServiceStatus(responseStatusCode int) string {
	if responseStatusCode >= 400 {
		return StatusDown
	}
	return StatusUp
}
