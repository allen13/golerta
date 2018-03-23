package config

import (
	"github.com/BurntSushi/toml"
	"github.com/allen13/golerta/app/algorithms"
	"github.com/allen13/golerta/app/auth/ldap"
	"github.com/allen13/golerta/app/auth/oauth"
	"github.com/allen13/golerta/app/auth/noop"
	"github.com/allen13/golerta/app/db/rethinkdb"
	"github.com/allen13/golerta/app/notifiers"
	"log"
	"time"
)

type GolertaConfig struct {
	Golerta       golerta
	Ldap          ldap.LDAPAuthProvider
	OAuth         oauth.OAuthAuthProvider
	Noop          noop.NoopAuthProvider
	Rethinkdb     rethinkdb.RethinkDB
	Notifiers     notifiers.Notifiers
	FlapDetection algorithms.FlapDetection
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

type golerta struct {
	BindAddr                string   `toml:"bind_addr"`
	SigningKey              string   `toml:"signing_key"`
	AuthProvider            string   `toml:"auth_provider"`
	ContinuousQueryInterval duration `toml:"continuous_query_interval"`
	LogAlertRequests        bool     `toml:"log_alert_requests"`
	TLSEnabled              bool     `toml:"tls_enabled"`
	TLSCert                 string   `toml:"tls_cert"`
	TLSKey                  string   `toml:"tls_key"`
	TLSAutoEnabled          bool     `toml:"tls_auto_enabled"`
	TLSAutoHosts            string   `toml:"tls_auto_hosts"`
}

func BuildConfig(configFile string) (config GolertaConfig) {
	_, err := toml.DecodeFile(configFile, &config)

	if err != nil {
		log.Fatal("config file error: " + err.Error())
	}

	setDefaultConfigs(&config)
	return
}

func setDefaultConfigs(config *GolertaConfig) {
	if config.Golerta.AuthProvider == "" {
		config.Golerta.AuthProvider = "noop"
	}
	if config.Golerta.ContinuousQueryInterval.Duration == 0 {
		config.Golerta.ContinuousQueryInterval.Duration = time.Second * 5
	}
}
