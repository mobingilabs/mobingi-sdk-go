package sesha3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mobingilabs/mobingi-sdk-go/client"
	"github.com/pkg/errors"
	"net/http"
	"os"
)

type Notificate struct {
	Slack  bool
	client client.HttpClient
}

func (w *Notificate) WebhookNotification(v interface{}) error {
	type payload_t struct {
		Text string `json:"text"`
	}

	var urls []string
	//webhook URLs
	if w.Slack {
		urls = append(urls, os.Getenv("SLACK"))
	}

	var err_string string

	switch v.(type) {
	case string:
		err := fmt.Errorf("%s", v.(string))
		err_string = fmt.Sprintf("%+v", errors.WithStack(err))
	case error:
		err_string = fmt.Sprintf("%+v", errors.WithStack(v.(error)))
	default:
		err_string = fmt.Sprintf("%s", v)
	}

	payload := payload_t{
		Text: err_string,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "payload marshal failed")
	}

	for _, ep := range urls {
		req, err := http.NewRequest(http.MethodGet, ep, bytes.NewBuffer(b))
		req.Header.Add("Content-Type", "application/json")
		_, _, err = w.client.Do(req)
		if err != nil {
			return errors.Wrap(err, "notification client do failed")
		}
	}

	return errors.New(err_string)
}
