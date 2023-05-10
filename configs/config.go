package configs

import "github.com/spf13/viper"

type Config struct {
	GroupID string `mapstructure:"GROIP_ID"`
	Token   string `mapstructure:"TOKEN"`
	UserID  string `mapstructure:"USER_ID"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("conf")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
