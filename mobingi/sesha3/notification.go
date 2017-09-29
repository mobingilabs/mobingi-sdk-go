package sesha3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/mobingilabs/mobingi-sdk-go/client"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

type Notificate struct {
	Slack  bool
	client client.HttpClient
	Cred   string
	Region string
}

type Event struct {
	Sesha3 string `dynamo:"server_name"`
	Slack  string `dynamo:"slack"`
}

func (w *Notificate) Dynamoget(key string) (string, error) {
	var results []Event
	cred := credentials.NewSharedCredentials("/root/.aws/credentials", w.Cred)
	db := dynamo.New(session.New(), &aws.Config{Region: aws.String(w.Region),
		Credentials: cred,
	})
	table := db.Table("SESHA3")
	err := table.Get("server_name", key).All(&results)
	url := results[0].Slack
	return url, err
}

func (w *Notificate) WebhookNotification(v interface{}) error {
	log.Println("start webhook")
	type payload_t struct {
		Text string `json:"text"`
	}

	var urls []string
	//webhook URLs
	log.Println("start get slack url")
	if w.Slack {
		slackURL, _ := w.Dynamoget("slack")
		urls = append(urls, slackURL)
	}

	log.Println("finish get slack url")
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
		req, err := http.NewRequest(http.MethodPost, ep, bytes.NewBuffer(b))
		req.Header.Add("Content-Type", "application/json")
		_, _, err = w.client.Do(req)
		if err != nil {
			return errors.Wrap(err, "notification client do failed")
		}
	}

	return errors.New(err_string)
}
