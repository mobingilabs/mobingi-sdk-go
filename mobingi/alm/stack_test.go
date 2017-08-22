package alm

import (
	"testing"

	"github.com/mobingilabs/mobingi-sdk-go/mobingi/session"
)

func TestNew(t *testing.T) {
	sess, _ := session.NewSession(&session.Config{})
	alm := New(sess)
	if alm == nil {
		t.Errorf("Expecting non-nil")
	}
}
