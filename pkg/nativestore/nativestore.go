package nativestore

import dcred "github.com/docker/docker-credential-helpers/credentials"

func Set(url, user, secret string) {
	cr := &dcred.Credentials{
		ServerURL: url,
		Username:  user,
		Secret:    secret,
	}

	ns.Add(cr)
}

func Get(url string) (string, string, error) {
	return ns.Get(url)
}
