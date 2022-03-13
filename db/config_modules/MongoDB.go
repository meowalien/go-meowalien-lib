package config_modules

type MangoDB struct {
	Host     string `json:"Host"`
	Database string `json:"Database"`
	User     string `json:"User"`
	Password string `json:"Password"`
	Port     string `json:"Port"`
}