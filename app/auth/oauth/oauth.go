package oauth

import (
	"bytes"
	"crypto/tls"
	"fmt"
	tk "github.com/allen13/golerta/app/auth/token"
	"net/http"
	"time"
	"io/ioutil"
	"errors"
)

type OAuthAuthProvider struct {
	Host         string `toml:"host"`
	ClientID     string `toml:"client_id"`
	ResponseType string `toml:"response_type"`
	signingKey   string
}

func (op *OAuthAuthProvider) SetSigningKey(key string) {
	op.signingKey = key
}

func (op *OAuthAuthProvider) Authenticate(username, password string) (authenticated bool, token string, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	url := fmt.Sprintf("https://%s/oauth/authorize", op.Host)

	body := bytes.NewBuffer([]byte(`{}`))
	req, err := http.NewRequest("POST", url, body)

	if err != nil {
		OAuthFailed(err)
		return
	}

	q := req.URL.Query()
	q.Add("client_id", op.ClientID)
	q.Add("response_type", op.ResponseType)
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(username, password)

	err = req.ParseForm()
	if err != nil {
		OAuthFailed(err)
		return
	}

	resp, err := client.Do(req)

	if err != nil {
		OAuthFailed(err)
		return
	}

	if resp.StatusCode == 302 {
		authenticated = true
		token = tk.CreateExpiringToken(username, op.signingKey, time.Hour*48, "oauth")
		return
	} else {
		message, _ := ioutil.ReadAll(resp.Body)
		OAuthFailed(errors.New(fmt.Sprintf("response code: %d message: %s", resp.StatusCode, message)))
		return
	}
}

func OAuthFailed(err error) {
	fmt.Println("oauth failed: ", err.Error())
}

func (op *OAuthAuthProvider) Connect() error {
	return nil
}
func (op *OAuthAuthProvider) Close() {

}
