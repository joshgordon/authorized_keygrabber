//Copied from stack overflow and modified
//http://stackoverflow.com/questions/16007695/verifying-a-signature-using-go-crypto-openpgp
package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/openpgp/packet"
	"io/ioutil"
	"os"
)

// gpg --export YOURKEYID --export-options export-minimal,no-export-attributes |
// hexdump /dev/stdin -v -e '/1 "%02X"'

func CheckSig(fileName string, sigFileName string, publicKeyHex string) error {
	// First, get the content of the file we have signed
	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	// Get a Reader for the signature file
	sigFile, err := os.Open(sigFileName)
	if err != nil {
		return err
	}

	defer func() {
		if err := sigFile.Close(); err != nil {
			panic(err)
		}
	}()

	// Read the signature file
	pack, err := packet.Read(sigFile)
	if err != nil {
		return err
	}

	// Was it really a signature file ? If yes, get the Signature
	signature, ok := pack.(*packet.Signature)
	if !ok {
		return errors.New(os.Args[2] + " is not a valid signature file.")
	}

	// For convenience, we have the key in hexadecimal, convert it to binary
	publicKeyBin, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return err
	}

	// Read the key
	pack, err = packet.Read(bytes.NewReader(publicKeyBin))
	if err != nil {
		return err
	}

	// Was it really a public key file ? If yes, get the PublicKey
	publicKey, ok := pack.(*packet.PublicKey)
	if !ok {
		return errors.New("Invalid public key.")
	}

	// Get the hash method used for the signature
	hash := signature.Hash.New()

	// Hash the content of the file (if the file is big, that's where you have
	// to change the code to avoid getting the whole file in memory, by reading
	// and writting in small chunks)
  _, err = hash.Write(fileContent)

	// Check the signature
	err = publicKey.VerifySignature(hash, signature)
	if err != nil {
		return err
	}

	return nil
}
