package rewrite

import (
	"fmt"
	"net/http"
	"time"
        "encoding/json"

	"github.com/karmanyaahm/up_rewrite/utils"
)

type Lomiri struct {
	Enabled bool `env:"UP_REWRITE_LOMIRI_ENABLE"`
	APIURL  string
}

func (l Lomiri) Path() string {
	if l.Enabled {
		return "/lomiri"
	}
	return ""
}

type lomiriSend struct {
	Token    string `json:"token"`
	AppId    string `json:"appid"`
	ExpireOn string `json:"expire_on"`
	Data     string `json:"data"`
	// clear all pending messages for appid
	ClearPending bool `json:"clear_pending,omitempty"`
	// replace pending messages with the same replace_tag
	ReplaceTag string `json:"replace_tag,omitempty"`
}

func (l Lomiri) Req(body []byte, req http.Request) (*http.Request, error) {
	token := req.URL.Query().Get("token")
	appid := req.URL.Query().Get("appid")

	var payload map[string]interface{}
        err := json.Unmarshal(body, &payload)
	newBody, err := utils.EncodeJSON(lomiriSend{
		Token:    token,
		AppId:    appid,
		ExpireOn: time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
		Data:     paylod,
	})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	newReq, err := http.NewRequest(http.MethodPost, l.APIURL, newBody)
	newReq.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	return newReq, nil
}

func (l Lomiri) RespCode(resp *http.Response) int {
	switch resp.StatusCode {
	case 429:
		return 429
	default:
		return 202
	}
}

func (l *Lomiri) Defaults() (failed bool) {
	if !l.Enabled {
		return
	}
	l.APIURL = "https://push.ubports.com/notify"
	return
}
