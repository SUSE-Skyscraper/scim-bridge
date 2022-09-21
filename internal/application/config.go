package application

import (
	"fmt"

	"github.com/spf13/viper"
)

const DefaultConfigDir = "."

type Configurator struct {
	viper *viper.Viper
}

type DBConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Port     int64  `mapstructure:"port"`
	Host     string `mapstructure:"host"`
}

type Argon2Config struct {
	MemoryCost  uint32 `mapstructure:"memory_cost"`
	TimeCost    uint32 `mapstructure:"time_cost"`
	Parallelism uint8  `mapstructure:"parallelism"`
}

type FGAConfig struct {
	APIScheme string `mapstructure:"api_scheme"`
	APIHost   string `mapstructure:"api_host"`
	StoreID   string `mapstructure:"store_id"`
}

func (db *DBConfig) GetDSN() string {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", db.User, db.Password, db.Host, db.Port, db.Database)
	return dsn
}

type ServerConfig struct {
	BaseURL string `mapstructure:"base_url"`
}

type Config struct {
	DB           DBConfig     `mapstructure:"db"`
	ServerConfig ServerConfig `mapstructure:"server"`
	Argon2Config Argon2Config `mapstructure:"argon2"`
	FGAConfig    FGAConfig    `mapstructure:"fga"`
}

func NewConfigurator(configDir string) Configurator {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configDir)

	return Configurator{
		viper: v,
	}
}

func (c *Configurator) Parse() (Config, error) {
	err := c.viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	conf := defaultConfig()
	err = c.viper.Unmarshal(&conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}

func defaultConfig() Config {
	return Config{
		FGAConfig: FGAConfig{
			APIScheme: "http",
			APIHost:   "127.0.0.1",
			StoreID:   "",
		},
		ServerConfig: ServerConfig{
			BaseURL: "http://localhost:8080",
		},
		Argon2Config: Argon2Config{
			MemoryCost:  64 * 1024,
			TimeCost:    30,
			Parallelism: 4,
		},
	}
}
