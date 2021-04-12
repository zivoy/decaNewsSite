package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func failedSimplePage() {
	// serve a basic page

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Encountered an error when getting configuration")
	})

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		panic(err)
	}
}

func getConfiguration(configUrl string, serverPassword string, aesKey string) (string, error) {
	// get the config data from the server
	req, err := http.NewRequest("GET", configUrl, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth("serverPoint", serverPassword)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New("invalid credentials")
	}
	defer resp.Body.Close()

	// read to a buffer and decrypt
	dataRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// clean newlines and spaces
	data := strings.Trim(string(dataRaw), " \n")

	return decrypt(data, aesKey)
}

func decrypt(encryptedString string, keyString string) (string, error) {

	key, err := hex.DecodeString(keyString)
	if err != nil {
		return "", err
	}
	enc, err := hex.DecodeString(encryptedString)
	if err != nil {
		return "", err
	}

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", plaintext), nil
}
