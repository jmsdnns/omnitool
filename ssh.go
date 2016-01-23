package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/pkg/sftp"
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

func dialServer(hostname string, config *ssh.ClientConfig) (*ssh.Client, error) {
	conn, err := ssh.Dial("tcp", hostname, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// SSHSession is a container for the pieces necessary to hold an SSH connection
// open in a goroutine
type SSHSession struct {
	hostname string
	config   *ssh.ClientConfig
	conn     *ssh.Client
	session  *ssh.Session
	error    error
}

func (s *SSHSession) init(hostname string, username string, keypath string) error {
	s.hostname = hostname

	// Instantiate config
	config, err := generateConfig(username, keypath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	s.config = config

	// Establish connection
	conn, err := dialServer(hostname, config)
	if err != nil {
		log.Fatal(err)
		return err
	}
	s.conn = conn

	// Make a session
	session, err := conn.NewSession()
	if err != nil {
		log.Fatal(err)
		return err
	}
	s.session = session

	return nil
}

func (s *SSHSession) tearDown() {
	s.session.Close()
	s.conn.Close()
}

func (s *SSHSession) executeCmd(cmd string) string {
	var stdoutBuf bytes.Buffer
	s.session.Stdout = &stdoutBuf
	s.session.Run(cmd)

	return stdoutBuf.String()
}

// GetSFTPClient returns an SFTP client from an existing SSHSession
func (s *SSHSession) GetSFTPClient() (*sftp.Client, error) {
	return sftp.NewClient(s.conn)
}

// SSHResponse contains the restul of running a command on a host via SSH
type SSHResponse struct {
	Hostname string
	Result   string
	Err      error
}

// ConnectToMachine takes host details and returns a connected SSHSession
func ConnectToMachine(address string, username string, keypath string) (*SSHSession, error) {
	session := &SSHSession{}
	err := session.init(address, username, keypath)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// MapCmd takes the details for a command and maps it, via SSH, across a list
// of hosts
func MapCmd(hostnames HostGroup, username string, keypath string, command string, results chan SSHResponse) {
	for _, hostname := range hostnames {
		go func(hostname string) {
			response := SSHResponse{Hostname: hostname}

			session, err := ConnectToMachine(hostname, username, keypath)
			defer session.tearDown()

			if err != nil {
				fmt.Println(err.Error())
				response.Err = err
			} else {
				result := session.executeCmd(command)
				response.Result = result
			}

			results <- response
		}(hostname)
	}
}

// MapScp takes the details for a file transfer and maps the file transfer
// across a list of hosts
func MapScp(hostnames HostGroup, username string, keypath string, localPath string, remotePath string, results chan SSHResponse) {
	for _, hostname := range hostnames {
		go func(hostname string) {
			response := SSHResponse{Hostname: hostname}

			session, err := ConnectToMachine(hostname, username, keypath)
			defer session.tearDown()

			sftpc, err := session.GetSFTPClient()
			if err != nil {
				fmt.Println(err.Error())
				response.Err = err
			}
			defer sftpc.Close()

			fmt.Println("PARTH:", filepath.Base(remotePath))
			// w := sftp.Walk(remotepath)
			// for w.Step() {
			// 	if w.Err() != nil {
			// 		continue
			// 	}
			// 	log.Println(w.Path())
			// }

			f, err := sftpc.Create("hello.txt")
			if err != nil {
				fmt.Println(err.Error())
				response.Err = err

			}
			if _, err := f.Write([]byte("Hello world!")); err != nil {
				fmt.Println(err.Error())
				response.Err = err
			}

			results <- response
		}(hostname)
	}
}
