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

func New(in *SeshaClientInput) (*Creds, error) {
	var err error
	if len(in.URL) < 1 {
		err = errors.Wrap(err, "url should not be empty")
		return &Creds{}, err
	}
	client, err := gottyclient.NewClient(in.URL)
	if err != nil {
		err = errors.Wrap(err, "sesha3 client creation failed")
		return &Creds{}, err
	}
	return &Creds{client: client}, err
}

func (c *Creds) Run() error {
	err := c.client.Loop()
	err = errors.Wrap(err, "can't connect sesha3")
	return err
}
