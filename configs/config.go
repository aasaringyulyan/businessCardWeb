package configs

type Config struct {
	DbHost   string `yaml:"dbHost"`
	DbPort   string `yaml:"dbPort"`
	User     string `yaml:"user"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
}

func NewConfig() *Config {
	return &Config{
		DbHost:   "",
		DbPort:   "",
		User:     "",
		Port:     "",
		Password: "",
	}
}
