package config_modules

type MangoDBConfiguration struct {
	User     string `json:"User"`
	Password string `json:"Password"`
	Host     string `json:"Host"`
	Port     string `json:"Port"`
	Database string `json:"Database"`
}