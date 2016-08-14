package sessions

import (
	"io/ioutil"
	"log"

	"golang.org/x/crypto/ssh"
)

//
// Connection Setup
//

func loadPrivateKey(filepath string) (ssh.Signer, error) {
	pemBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return signer, nil
}

func generateConfig(username string, keypath string) (*ssh.ClientConfig, error) {
	signer, err := loadPrivateKey(keypath)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
	}

	return config, nil
}
