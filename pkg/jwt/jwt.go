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

type CustomClaims struct {
	Username string
	Passwd   string
	jwt.StandardClaims
}

type jwtctx struct {
	name    string
	rsa     string
	pub     []byte
	User    string
	PemPub  string
	PemPriv string
	Reuse   bool
}

func (j *jwtctx) GenerateToken() (*jwt.Token, string, error) {
	var clms CustomClaims
	var stoken string

	expire := time.Hour * 24
	clms.ExpiresAt = time.Now().Add(expire).Unix()
	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS512"), clms)
	defkey, err := ioutil.ReadFile(j.PemPriv)
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

func NewCtx(user string) (*jwtctx, error) {
	var ctx jwtctx

	tmpdir := os.TempDir() + "/token/rsa/"
	debug.Info("tmp:", tmpdir)
	ctx.User = user
	ctx.PemPub = tmpdir + user + ".pem.pub"
	ctx.PemPriv = tmpdir + user + ".pem"
	ctx.Reuse = true

	// create dir if necessary
	if !private.Exists(tmpdir) {
		err := os.MkdirAll(tmpdir, 0700)
		if err != nil {
			return nil, errors.Wrap(err, "mkdirall failed")
		}
	}

	// create public and private pem files
	if !private.Exists(ctx.PemPub) || !private.Exists(ctx.PemPriv) {
		ctx.Reuse = false // we are not reusing pem files
		priv, err := rsa.GenerateKey(rand.Reader, 2048)
		privder := x509.MarshalPKCS1PrivateKey(priv)
		if err != nil {
			return nil, errors.Wrap(err, "MarshalPKCS1PrivateKey failed")
		}

		pubkey := priv.Public()
		pubder, err := x509.MarshalPKIXPublicKey(pubkey)
		if err != nil {
			return nil, errors.Wrap(err, "MarshalPKIXPublicKey failed")
		}

		pubblock := &pem.Block{Type: "RSA PUBLIC KEY", Bytes: pubder}
		pemblock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privder}
		pubfile, err := os.OpenFile(ctx.PemPub, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return nil, errors.Wrap(err, "open file failed (pub)")
		}

		defer pubfile.Close()
		privfile, err := os.OpenFile(ctx.PemPriv, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return nil, errors.Wrap(err, "open file failed (priv)")
		}

		defer privfile.Close()
		pem.Encode(pubfile, pubblock)
		pem.Encode(privfile, pemblock)
	}

	return &ctx, nil
}
