package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort              string `mapstructure:"port"`
	ServerMode              string `mapstructure:"mode"`
	ServerReadTimeout       int    `mapstructure:"read_timeout"`
	ServerWriteTimeout      int    `mapstructure:"write_timeout"`
	ServerCtxDefaultTimeout int    `mapstructure:"ctx_default_timeout"`
	ServerDebug             bool   `mapstructure:"debug"`

	DBHost            string `mapstructure:"db_host"`
	DBUser            string `mapstructure:"db_user"`
	DBPassword        string `mapstructure:"db_password"`
	DBName            string `mapstructure:"db_name"`
	DBSchema          string `mapstructure:"db_schema"`
	DBPort            string `mapstructure:"db_port"`
	DBSSLMode         string `mapstructure:"db_sslmode"`
	DBMaxIdleConns    int    `mapstructure:"postgresql_max_idle_conns"`
	DBMaxOpenConns    int    `mapstructure:"postgresql_max_open_conns"`
	DBConnMaxLifeTime int    `mapstructure:"postgresql_con_max_life_time"`
	DBTimeout         int    `mapstructure:"postgresql_timeout"`
	DBDebug           bool   `mapstructure:"postgresql_debug"`

	LoggerEncoding         string `mapstructure:"logger_encoding"`
	LoggerLevel            string `mapstructure:"logger_level"`
	LoggerIsFullPathCaller bool   `mapstructure:"logger_full_path_caller"`

	JWTSecret        string `mapstructure:"jwt_secret"`
	JWTExpireMinutes int    `mapstructure:"jwt_expire_minutes"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.BindEnv("port", "PORT")
	viper.BindEnv("mode", "MODE")
	viper.BindEnv("read_timeout", "READ_TIMEOUT")
	viper.BindEnv("write_timeout", "WRITE_TIMEOUT")
	viper.BindEnv("ctx_default_timeout", "CTX_DEFAULT_TIMEOUT")
	viper.BindEnv("debug", "DEBUG")

	viper.BindEnv("db_host", "DB_HOST")
	viper.BindEnv("db_user", "DB_USER")
	viper.BindEnv("db_password", "DB_PASSWORD")
	viper.BindEnv("db_name", "DB_NAME")
	viper.BindEnv("db_schema", "DB_SCHEMA")
	viper.BindEnv("db_port", "DB_PORT")
	viper.BindEnv("db_sslmode", "DB_SSLMODE")
	viper.BindEnv("postgresql_max_idle_conns", "POSTGRESQL_MAX_IDLE_CONNS")
	viper.BindEnv("postgresql_max_open_conns", "POSTGRESQL_MAX_OPEN_CONNS")
	viper.BindEnv("postgresql_con_max_life_time", "POSTGRESQL_CON_MAX_LIFE_TIME")
	viper.BindEnv("postgresql_timeout", "POSTGRESQL_TIMEOUT")
	viper.BindEnv("postgresql_debug", "POSTGRESQL_DEBUG")

	viper.BindEnv("jwt_secret", "JWT_SECRET")
	viper.BindEnv("jwt_expire_minutes", "JWT_EXPIRE_MINUTES")

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return cfg, nil
}
