package config

type Config struct {
	Wappsto struct {
		Username string `toml:"username"`
		Password string `toml:"password"`
		Server   string `toml:"server"`
	} `toml:"wappsto"`

	Kafka struct {
		Connect string `toml:"connect"`
	}
}

var C Config
