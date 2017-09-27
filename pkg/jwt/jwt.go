package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/mobingilabs/mobingi-sdk-go/pkg/private"
	"github.com/pkg/errors"
)

var rsainit bool

func init() {
	tmpdir := os.TempDir() + "/sesha3/rsa/"
	debug.Info("tmp:", tmpdir)
	pub := tmpdir + "token.pem.pub"
	prv := tmpdir + "token.pem"

	// create dir if necessary
	if !private.Exists(tmpdir) {
		err := os.MkdirAll(tmpdir, 0700)
		if err != nil {
			return
		}
	}

	// create public and private pem files
	if !private.Exists(pub) || !private.Exists(prv) {
		priv, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return
		}

		privder := x509.MarshalPKCS1PrivateKey(priv)
		pubkey := priv.Public()
		pubder, err := x509.MarshalPKIXPublicKey(pubkey)
		if err != nil {
			return
		}

		pubblock := &pem.Block{Type: "RSA PUBLIC KEY", Bytes: pubder}
		pemblock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privder}
		pubfile, err := os.OpenFile(pub, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return
		}

		defer pubfile.Close()
		prvfile, err := os.OpenFile(prv, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return
		}

		defer prvfile.Close()
		pem.Encode(pubfile, pubblock)
		pem.Encode(prvfile, pemblock)
	}

	rsainit = true
}

type CustomClaims struct {
	Username string
	Passwd   string
	jwt.StandardClaims
}

type jwtctx struct {
	name   string
	rsa    string
	pub    []byte
	PemPub string
	PemPrv string
}

func (j *jwtctx) GenerateToken() (*jwt.Token, string, error) {
	var clms CustomClaims
	var stoken string

	expire := time.Hour * 24
	clms.ExpiresAt = time.Now().Add(expire).Unix()
	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS512"), clms)
	defkey, err := ioutil.ReadFile(j.PemPrv)
	if err != nil {
		return token, stoken, errors.Wrap(err, "readfile failed")
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(defkey)
	if err != nil {
		return token, stoken, errors.Wrap(err, "parse priv key from pem failed")
	}

	stoken, err = token.SignedString(key)
	return token, stoken, nil
}

func NewCtx() (*jwtctx, error) {
	if !rsainit {
		return nil, errors.New("failed in rsa init")
	}

	var ctx jwtctx
	tmpdir := os.TempDir() + "/sesha3/rsa/"
	ctx.PemPub = tmpdir + "token.pem.pub"
	ctx.PemPrv = tmpdir + "token.pem"
	return &ctx, nil
}
