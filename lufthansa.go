package lufthansa

import (
	"encoding/json"
	"log"
	"net/url"
	"time"

	"github.com/cmodk/go-simplehttp"
	"github.com/sirupsen/logrus"
)

type Lufthansa struct {
	lg     *logrus.Logger
	sh     simplehttp.SimpleHttp
	key    string
	secret string
	token  string
	debug  bool
}

func New(k string, s string, logger *logrus.Logger) *Lufthansa {
	lufthansa := Lufthansa{
		lg:     logger,
		sh:     simplehttp.New("https://api.lufthansa.com/v1", logger),
		key:    k,
		secret: s,
	}

	lufthansa.sh.AddHeader("Accept", "application/json")

	go handleToken(&lufthansa)

	return &lufthansa

}

func (lufthansa *Lufthansa) SetDebug(d bool) {
	lufthansa.debug = d
	lufthansa.sh.SetDebug(d)
}

func handleToken(lh *Lufthansa) {
	data := url.Values{}
	data.Set("client_id", lh.key)
	data.Set("client_secret", lh.secret)
	data.Set("grant_type", "client_credentials")
	for {
		resp, err := lh.sh.Post("/oauth/token", data)
		if err != nil {
			lh.lg.WithField("error", err).Error("Could not get token")
			continue
		}

		log.Printf("Resp: %s\n", resp)

		t := struct {
			Token     string `json:"access_token"`
			TokenType string `json:"token_type"`
			ExpiresIn int    `json:"expires_in"`
		}{}

		if err := json.Unmarshal([]byte(resp), &t); err != nil {
			lh.lg.WithField("error", err).Error("Could not parse json")
			continue
		}

		lh.token = t.Token

		log.Printf("Sleeping for %d seconds\n", t.ExpiresIn)
		time.Sleep(time.Second * time.Duration(t.ExpiresIn))

	}

}
