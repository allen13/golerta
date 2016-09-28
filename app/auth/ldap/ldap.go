package ldap

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/allen13/golerta/app/auth/token"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/ldap.v2"
	"time"
)

type LDAPAuthProvider struct {
	conn         *ldap.Conn
	signingKey   string
	Host         string   `toml:"host"`
	Port         int      `toml:"port"`
	UseSSL       bool     `toml:"use_ssl"`
	BaseDN       string   `toml:"base_dn"`
	BindDN       string   `toml:"bind_dn"`
	BindPassword string   `toml:"bind_password"`
	UserFilter   string   `toml:"user_filter"`
	Attributes   []string `toml:"attributes"`
}

func (lc *LDAPAuthProvider) createToken(username string) (string, error) {
	expirationTimestamp := time.Now().Add(time.Hour * 48).Unix()
	claims := jwt.MapClaims{
		"jti":  username,
		"iss":  "ldap",
		"exp":  expirationTimestamp,
		"name": username,
		//Everyone who logs in is an admin by default for now. Could check ldap groups for this.
		"role": "admin",
	}

	return token.CreateToken(lc.signingKey, claims)
}

func (lc *LDAPAuthProvider) SetSigningKey(key string) {
	lc.signingKey = key
}

// Connect connects to the ldap backend
func (lc *LDAPAuthProvider) Connect() error {
	if lc.conn == nil {
		var l *ldap.Conn
		var err error
		address := fmt.Sprintf("%s:%d", lc.Host, lc.Port)
		if !lc.UseSSL {
			l, err = ldap.Dial("tcp", address)
			if err != nil {
				return err
			}

			// Reconnect with TLS
			err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
			if err != nil {
				return err
			}
		} else {
			l, err = ldap.DialTLS("tcp", address, &tls.Config{InsecureSkipVerify: true, ServerName: lc.Host})
			if err != nil {
				return err
			}
		}

		lc.conn = l
	}
	return nil
}

// Close closes the ldap backend connection
func (lc *LDAPAuthProvider) Close() {
	if lc.conn != nil {
		lc.conn.Close()
		lc.conn = nil
	}
}

// Authenticate authenticates the user against the ldap backend
func (lc *LDAPAuthProvider) Authenticate(username, password string) (authenticated bool, token string, err error) {
	err = lc.Connect()
	defer lc.Close()

	if err != nil {
		return
	}

	// First bind with a read only user
	if lc.BindDN != "" && lc.BindPassword != "" {
		err = lc.conn.Bind(lc.BindDN, lc.BindPassword)
		if err != nil {
			return
		}
	}

	attributes := append(lc.Attributes, "dn")
	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		lc.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(lc.UserFilter, username),
		attributes,
		nil,
	)

	sr, err := lc.conn.Search(searchRequest)
	if err != nil {
		return
	}

	if len(sr.Entries) < 1 {
		err = errors.New("User does not exist")
		return
	}

	if len(sr.Entries) > 1 {
		err = errors.New("Too many entries returned")
		return
	}

	userDN := sr.Entries[0].DN
	user := map[string]string{}
	for _, attr := range lc.Attributes {
		user[attr] = sr.Entries[0].GetAttributeValue(attr)
	}

	// Bind as the user to verify their password
	err = lc.conn.Bind(userDN, password)
	if err != nil {
		return
	}

	token, err = lc.createToken(username)
	if err != nil {
		return
	}

	//We authenticated and we have our token
	authenticated = true

	// Rebind as the read only user for any further queries
	if lc.BindDN != "" && lc.BindPassword != "" {
		err = lc.conn.Bind(lc.BindDN, lc.BindPassword)
		if err != nil {
			return
		}
	}

	return
}
