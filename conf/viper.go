package conf

import (
	"github.com/spf13/viper"
)

func NewViper() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.AddConfigPath(".")

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	v.AutomaticEnv()

	return v, nil
}
