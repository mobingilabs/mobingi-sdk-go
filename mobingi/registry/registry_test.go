package registry

import (
	"os"
	"testing"

	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
)

func TestGetRegistryToken(t *testing.T) {
	return
	if os.Getenv("MOBINGI_CLIENT_ID") != "" && os.Getenv("MOBINGI_CLIENT_SECRET") != "" &&
		os.Getenv("MOBINGI_USERNAME") != "" && os.Getenv("MOBINGI_PASSWORD") != "" {
		sess, _ := session.New(&session.Config{
			BaseApiUrl: "https://apidev.mobingi.com",
		})

		reg := New(sess)
		in := &GetRegistryTokenInput{
			Scope: "repository:" + os.Getenv("MOBINGI_USERNAME") + "/hello:*",
		}

		resp, body, token, err := reg.GetRegistryToken(in)
		if err != nil {
			t.Errorf("expecting nil error, received %v", err)
		}

		_, _, _ = resp, body, token
	}
}
