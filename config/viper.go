package config

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"log/slog"
	"os"
)

var logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))

type DBConfig struct {
	Engine   string `json:"engine"`
	Server   string `json:"server"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Schema   string `json:"schema"`
}
type ServiceConfig struct {
	Port int `json:"port"`
}

type YamlConfig struct {
	DBConfig      DBConfig      `yaml:"dbdetails"`
	ServiceConfig ServiceConfig `yaml:"servicedetails"`
}

func init() {
	ReadConfig("devconfig")
}

func ReadConfig(filename string) {
	viper.AddConfigPath("./env")
	viper.SetConfigName(filename)
	viper.SetConfigType("toml")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Error(err.Error())
	}
}

func GetConfigValues() *YamlConfig {

	db := &DBConfig{
		Engine:   viper.GetString("DATABASE.engine"),
		Server:   viper.GetString("DATABASE.server"),
		Username: viper.GetString("DATABASE.username"),
		Password: viper.GetString("DATABASE.password"),
		Port:     viper.GetInt("DATABASE.port"),
		Schema:   viper.GetString("DATABASE.schema"),
	}
	server := &ServiceConfig{
		Port: viper.GetInt("SERVICE.port"),
	}
	return &YamlConfig{*db, *server}
}
