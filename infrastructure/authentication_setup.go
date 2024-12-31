package infrastructure

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"io/ioutil"
	"log"
)

const Alogirthm = "RS256"

type JWTAuth struct {
	alg       jwa.SignatureAlgorithm
	signKey   interface{}
	verifyKey interface{}
	verifier  jwt.ParseOption
}

func loadAuthToken() error {
	// Load private key
	privateReader, err := ioutil.ReadFile(privatePath)
	if err != nil {
		log.Println("No RSA private pem file")
		return err
	}
	privatePem, _ := pem.Decode(privateReader)

	if privatePem.Type != "PRIVATE KEY" {
		return errors.New("invalid private key type")
	}

	// Parse PKCS#8 private key
	privateKey, err := x509.ParsePKCS8PrivateKey(privatePem.Bytes)
	if err != nil {
		log.Println(err)
		return err
	}

	// Nếu là private key RSA, ép kiểu về *rsa.PrivateKey
	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return errors.New("private key is not RSA")
	}

	// Load public key
	publicReader, err := ioutil.ReadFile(publicPath)
	if err != nil {
		log.Println("No RSA public pem file")
		return err
	}
	publicPem, _ := pem.Decode(publicReader)
	publicKey, err := x509.ParsePKIXPublicKey(publicPem.Bytes)
	if err != nil {
		log.Println(err)
		return err
	}

	encodeAuth = New(Alogirthm, rsaPrivateKey, publicKey)
	decodeAuth = New(Alogirthm, nil, publicKey)

	return nil
}

func RsaEncrypt(decrypt string) (string, error) {
	// Load public key
	publicKey, err := ioutil.ReadFile(publicPath)
	if err != nil {
		InfoLog.Println("No RSA public pem file")
		return "", err
	}
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return "", errors.New("public key error!")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Println(err)
		return "", err
	}
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, pub.(*rsa.PublicKey), []byte(decrypt))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherText), nil
}
func RsaDecrypt(encrypt string) ([]byte, error) {
	// Load private key
	privateKey, err := ioutil.ReadFile(privatePath)
	if err != nil {
		log.Println("No RSA private pem file")
		return nil, err
	}
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rsaPrivateKey, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not RSA")
	}
	cipherText, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		log.Println("Error decoding base64:", err)
		return nil, err
	}
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, rsaPrivateKey, cipherText)
	if err != nil {
		log.Println("Error decrypting RSA:", err)
		return nil, err
	}
	return plainText, nil
}

// New creates a new JWTAuth instance
func New(alg string, signKey interface{}, verifyKey interface{}) *JWTAuth {
	ja := &JWTAuth{alg: jwa.SignatureAlgorithm(alg), signKey: signKey, verifyKey: verifyKey}

	if ja.verifyKey != nil {
		ja.verifier = jwt.WithVerify(ja.alg, ja.verifyKey)
	} else {
		ja.verifier = jwt.WithVerify(ja.alg, ja.signKey)
	}

	return ja
}
func (ja *JWTAuth) Encode(claims map[string]interface{}) (t jwt.Token, tokenString string, err error) {
	t = jwt.New()
	for k, v := range claims {
		t.Set(k, v)
	}
	payload, err := ja.sign(t)
	if err != nil {
		return nil, "", err
	}
	tokenString = string(payload)
	return
}
func (ja *JWTAuth) sign(token jwt.Token) ([]byte, error) {
	return jwt.Sign(token, ja.alg, ja.signKey)
}
