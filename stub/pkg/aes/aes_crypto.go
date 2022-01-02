package aes

import (
	"errors"
	"crypto/aes"
	"crypto/cipher"

)

/*
	chave = Chave simétrica usada no processo de criptografia/descriptografia
	nonce = Número aleatório gerado pelo sistema para cada cifra (i.e IV)
	gcm = Galois/Counter Mode - Modo de operação

*/ 
type AesCrypto struct {
	Chave string
	nonce []byte
	gcm cipher.AEAD
}

func (a *AesCrypto) NewAesCrypto(nonce []byte) error {
	c, err := aes.NewCipher([]byte(a.Chave))
	if err != nil {
		return err
	}

	a.gcm, err = cipher.NewGCM(c)
	if err != nil {
		return err
	}

	if len(nonce) != a.gcm.NonceSize() {
		return errors.New("nonce size is not valid")
	}

	a.nonce = nonce 
	return nil
}

func (a *AesCrypto) Encriptar(texto []byte) []byte {
	cifrado := a.gcm.Seal(nil, a.nonce, texto, nil)
	return cifrado
}

func (a *AesCrypto) Decriptar(cifrado []byte) ([]byte, error) {
	texto, err := a.gcm.Open(nil, a.nonce, cifrado, nil)
	if err != nil {
		return nil, err 
	}

	return texto, nil
}