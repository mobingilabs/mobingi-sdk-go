package sesha3

import (
	"github.com/moul/gotty-client"
	"github.com/pkg/errors"
)

type SeshaClientInput struct {
	URL string
}

type Creds struct {
	client *gottyclient.Client
}

func New(in *SeshaClientInput) (ret *Creds, err error) {
	if len(in.URL) < 1 {
		err = errors.Wrap(err, "can't create sesha3 client")
		return
	}
	client, err := gottyclient.NewClient(in.URL)
	if err != nil {
		err = errors.Wrap(err, "can't create sesha3 client")
		return
	}
	ret = &Creds{client: client}
	return ret, err
}

func (c *Creds) Run() (err error) {
	err = c.client.Loop()
	return err
}
