package configs

type MQ struct {
	Addresses []string `yaml:"addresses"`
	Port      int      `yaml:"port"`
	Username  string   `yaml:"username"`
	Password  string   `yaml:"password"`
	ClientID  string   `yaml:"client_id"`
}
