package utils

import "github.com/spf13/viper"

type StorageConfig struct {
	PgUsername string `mapstructure:"POSTGRES_USER"`
	PgPassword string `mapstructure:"POSTGRES_PASSWORD"`
	PgHost     string `mapstructure:"POSTGRES_HOST"`
	PgPort     string `mapstructure:"POSTGRES_PORT"`
	PgDatabase string `mapstructure:"POSTGRES_DB"`
	Host       string `mapstructure:"HOST"`
	PortGrpc   string `mapstructure:"PORTGRPC"`
	PortHttp   string `mapstructure:"PORTHTTP"`
}

// LoadConfig Конструктор для создания StorageConfig, который содержит считанные из .env файла данные.
func LoadConfig(path string) (config StorageConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
