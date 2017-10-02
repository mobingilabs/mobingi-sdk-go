package sesha3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

type Notificate struct {
	Slack  bool
	Cred   string
	Region string
	URLs   EventN
}

type EventN struct {
	Server_name string `dynamo:"server_name"`
	Slack       string `dynamo:"slack"`
}

func (n *Notificate) Dynamoget() (EventN, error) {
	serverName := "sesha3"
	var results []EventN
	cred := credentials.NewSharedCredentials("/root/.aws/credentials", n.Cred)
	db := dynamo.New(session.New(), &aws.Config{Region: aws.String(n.Region),
		Credentials: cred,
	})
	table := db.Table("SESHA3")
	err := table.Get("server_name", serverName).All(&results)
	url := results[0]
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
		NotificateURL := w.URLs
		urls = append(urls, NotificateURL.Slack)
	}

	log.Println("finish get slack url")
	var err_string string

	switch v.(type) {
	case string:
		err := v.(string)
		err_string = fmt.Sprintf("%v", err)
	case error:
		err_string = fmt.Sprintf("%+v", errors.WithStack(v.(error)))
	default:
		err_string = fmt.Sprintf("%s", v)
	}

	err_string = "```" + err_string + "```"
	payload := payload_t{
		Text: err_string,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "payload marshal failed")
	}

	client := &http.Client{}
	for _, ep := range urls {
		req, err := http.NewRequest(http.MethodPost, ep, bytes.NewBuffer(b))
		req.Header.Add("Content-Type", "application/json")
		_, err = client.Do(req)
		if err != nil {
			return errors.Wrap(err, "notification client do failed")
		}
	}

	return err
}
