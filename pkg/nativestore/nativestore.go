package nativestore

import (
	"errors"

	dcred "github.com/docker/docker-credential-helpers/credentials"
)

const CliUrl = "github.com/mobingi/mobingi-cli"

func Set(url, user, secret string) error {
	pmyns := nativeStore()
	if pmyns == nil {
		return errors.New("native store not supported yet")
	}

	cr := &dcred.Credentials{
		ServerURL: url,
		Username:  user,
		Secret:    secret,
	}

	myns := *pmyns
	myns.Add(cr)
	return nil
}

func Get(url string) (string, string, error) {
	pmyns := nativeStore()
	if pmyns == nil {
		return "", "", errors.New("native store not supported yet")
	}

	myns := *pmyns
	return myns.Get(url)
}