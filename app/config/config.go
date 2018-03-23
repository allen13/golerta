package config

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/allen13/golerta/app/algorithms"
	"github.com/allen13/golerta/app/auth/ldap"
	"github.com/allen13/golerta/app/auth/oauth"
	"github.com/allen13/golerta/app/auth/noop"
	"github.com/allen13/golerta/app/db/rethinkdb"
	"github.com/allen13/golerta/app/notifiers"
)

type GolertaConfig struct {
	App           app
	Ldap          ldap.LDAPAuthProvider
	OAuth         oauth.OAuthAuthProvider
	Noop          noop.NoopAuthProvider
	Rethinkdb     rethinkdb.RethinkDB
	Notifiers     notifiers.Notifiers
	FlapDetection algorithms.FlapDetection
}

type duration struct {
	time.Duration `mapstructure:"continuous_query_interval"`
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

type app struct {
	BindAddr                string   `mapstructure:"bind_addr"`
	SigningKey              string   `mapstructure:"signing_key"`
	AuthProvider            string   `mapstructure:"auth_provider"`
	ContinuousQueryInterval duration
	LogAlertRequests        bool     `mapstructure:"log_alert_requests"`
	TLSEnabled              bool     `mapstructure:"tls_enabled"`
	TLSCert                 string   `mapstructure:"tls_cert"`
	TLSKey                  string   `mapstructure:"tls_key"`
	TLSAutoEnabled          bool     `mapstructure:"tls_auto_enabled"`
	TLSAutoHosts            string   `mapstructure:"tls_auto_hosts"`
}

func BuildConfig(configFile string) (config GolertaConfig) {
	viper.SetEnvPrefix("golerta")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName(configFile)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/golerta/")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal("could not read config: " + err.Error())
	}

	err = viper.Unmarshal(&config)

	if err != nil {
		log.Fatal("config file error: " + err.Error())
	}

	setDefaultConfigs(&config)
	return
}

func setDefaultConfigs(config *GolertaConfig) {
	if config.App.AuthProvider == "" {
		config.App.AuthProvider = "noop"
	}
	if config.App.ContinuousQueryInterval.Duration == 0 {
		config.App.ContinuousQueryInterval.Duration = time.Second * 5
	}
}
