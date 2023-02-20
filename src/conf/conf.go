package conf

var Properties *ApplicationProperties

type ApplicationProperties struct {
	DataPath    string `json:"dataPath"`
	DataUrl     string `json:"dataUrl"`
	SecretKey   string `json:"secretKey"`
	Hostname    string `json:"hostname"`
	DatabaseUrl string `json:"databaseUrl"`
	Port        int    `json:"port"`
	Https       Https  `json:"https"`
	Redis       Redis  `json:"redis"`
}

type Https struct {
	Enable   bool   `json:"enable"`
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}

type Redis struct {
	Enable   bool   `json:"enable"`
	Addr     string `json:"addr"`
	Password string `json:"password"`
	Db       int    `json:"db"`
}
