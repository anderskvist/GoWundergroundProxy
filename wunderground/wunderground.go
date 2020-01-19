package wunderground

// ResponseOutputSensor contains sensor data
// We need to submit all values as strings to MQTT, so we just save them as strings here
type WundergroundData struct {
	Indoortempf    string `json:"indoortempf.string"`
	Tempf          string `json:"tempf,string"`
	Dewptf         string `json:"dewptf,string"`
	Windchillf     string `json:"windchillf,string"`
	Indoorhumidity string `json:"indoorhumidity,string"`
	Humidity       string `json:"humidity,string"`
	Windspeedmph   string `json:"windspeedmph,string"`
	Windgustmph    string `json:"windgustmph,string"`
	Winddir        string `json:"winddir,string"`
	Absbaromin     string `json:"absbaromin,string"`
	Baromin        string `json:"baromin,string"`
	Rainin         string `json:"rainin,string"`
	Dailyrainin    string `json:"dailyrainin,string"`
	Weeklyrainin   string `json:"weeklyrainin,string"`
	Monthlyrainin  string `json:"monthlyrainin,string"`
	Yearlyrainin   string `json:"yearlyrainin,string"`
	Solarradiation string `json:"solarradiation,string"`
	UV             string `json:"UV,string"`
}

func Parse(keys map[string][]string) WundergroundData {

	data := WundergroundData{
		Indoortempf:    keys["indoortempf"][0],
		Tempf:          keys["tempf"][0],
		Dewptf:         keys["dewptf"][0],
		Windchillf:     keys["windchillf"][0],
		Indoorhumidity: keys["indoorhumidity"][0],
		Humidity:       keys["humidity"][0],
		Windspeedmph:   keys["windspeedmph"][0],
		Windgustmph:    keys["windgustmph"][0],
		Winddir:        keys["winddir"][0],
		Absbaromin:     keys["absbaromin"][0],
		Baromin:        keys["baromin"][0],
		Rainin:         keys["rainin"][0],
		Dailyrainin:    keys["dailyrainin"][0],
		Weeklyrainin:   keys["weeklyrainin"][0],
		Monthlyrainin:  keys["monthlyrainin"][0],
		Yearlyrainin:   keys["yearlyrainin"][0],
		Solarradiation: keys["solarradiation"][0],
		UV:             keys["UV"][0]}
	return data
}
