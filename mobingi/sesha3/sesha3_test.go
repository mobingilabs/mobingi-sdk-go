package sesha3

import (
	"testing"
)

func TestNew(t *testing.T) {
	url := &SeshaClientInput{URL: "https://sesha3.labs.mobingi.com:32914/6h6fcxo233gby42bdvu93nzhl3rgorq6t7yd/"}
	cli, _ := New(url)
	err := cli.Run()
	_ = err
}
