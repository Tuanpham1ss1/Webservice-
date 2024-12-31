package infrastructure

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
)

type RSAService interface {
	RsaEncrypt(decrypt string) (string, error)
	RsaDecrypt(encrypt string) (string, error)
}
type rsaService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func (r *rsaService) RsaEncrypt(decrypt string) (string, error) {
	// Encrypt
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, []byte(decrypt))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return base64.StdEncoding.EncodeToString(cipherText), nil
}
func (r *rsaService) RsaDecrypt(encrypt string) (string, error) {
	// Decrypt
	cipherText, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		log.Println("Error decoding base64:", err) // Log lỗi nếu không thể giải mã base64
		return "", err
	}
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, r.privateKey, cipherText)
	if err != nil {
		log.Println("Error decrypting RSA:", err) // Log lỗi nếu không thể giải mã RSA
		return "", err
	}
	return string(plainText), nil
}

func NewRSAService() RSAService {
	object := &rsaService{}
	err := object.setPrivateKey()
	if err != nil {
		log.Println(err)
	}
	err = object.setPublicKey()
	if err != nil {
		log.Println(err)
	}
	return object
}

func (r *rsaService) setPrivateKey() error {
	// Read private key
	key, err := ioutil.ReadFile(privatePath)
	if err != nil {
		InfoLog.Println("NO RSA private pem file")
		return errors.New("no RSA private pem file")
	}
	// Decode private key
	block, _ := pem.Decode(key)
	if block == nil {
		return errors.New("private key error!")
	}
	// Parse private key
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Println(err)
		return err
	}
	r.privateKey = priv.(*rsa.PrivateKey)
	return nil
}
func (r *rsaService) setPublicKey() error {
	// Read public key
	key, err := ioutil.ReadFile(publicPath)
	if err != nil {
		InfoLog.Println("No RSA public pem file")
		return errors.New("no RSA public pem file")
	}
	// Decode public key
	block, _ := pem.Decode(key)
	if block == nil {
		return errors.New("public key error!")
	}
	// Parse public key
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	r.publicKey = pub.(*rsa.PublicKey)
	return nil
}
