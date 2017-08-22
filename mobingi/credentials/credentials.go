package credentials

import (
	"github.com/mobingilabs/mobingi-sdk-go/client"
	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
)

type creds struct {
	session *session.Session
	client  client.HttpClient
}

func (c *creds) List() (*client.Response, []byte, error) {
	return nil, nil, nil
}
