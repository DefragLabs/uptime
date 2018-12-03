package tasks

// MINUTE is the minute const.
const MINUTE = "minute"

// SECOND is the second const.
const SECOND = "second"

// MonitoringConfig stores the acceptable values for frequency & unit.
var MonitoringConfig = map[string][]int32{
	SECOND: []int32{30},
	MINUTE: []int32{1, 5, 15, 30},
}

// GetMonitoringConfigKeys returns the keys from monitoringConfig.
func GetMonitoringConfigKeys() []string {
	keys := make([]string, len(MonitoringConfig))

	for key := range MonitoringConfig {
		keys = append(keys, key)
	}

	return keys
}

// FrequencyInMonitoringConfig checks if the given frequency exists in MonitoringConfig
func FrequencyInMonitoringConfig(a int32, list []int32) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
