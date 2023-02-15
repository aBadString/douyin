package conf

var Properties *ApplicationProperties

type ApplicationProperties struct {
	DataPath    string `json:"dataPath"`
	DataUrl     string `json:"dataUrl"`
	SecretKey   string `json:"secretKey"`
	Hostname    string `json:"hostname"`
	DatabaseUrl string `json:"databaseUrl"`
}
