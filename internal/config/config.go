package config

import "github.com/spf13/viper"

type Config struct {
	EnvServerPort          string `mapstructure:"TODO_PORT"`
	EnvDatabasePath        string `mapstructure:"TODO_DBFILE"`
	EnvPassword            string `mapstructure:"TODO_PASSWORD"`
	EnvDatabasePathForTest string `mapstructure:"TODO_DBFILE_TEST"`
	EnvPassHashForTest     string `mapstructure:"TODO_HASH_FOR_TEST"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
