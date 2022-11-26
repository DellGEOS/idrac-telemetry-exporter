package redfishmetricreport

type Config struct {
	GlobalConfig GlobalConfig `yaml:"global"`
	Idracs []IdracConfig `yaml:"idracs"`
}

type GlobalConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type IdracConfig struct {
	IpAddress string `yaml:"ipAddress"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
