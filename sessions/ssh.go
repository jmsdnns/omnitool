package sessions

import (
	"bytes"
	"fmt"
	"log"

	"github.com/jmsdnns/omnitool/hosts"
	"golang.org/x/crypto/ssh"
)

func dialServer(hostname string, config *ssh.ClientConfig) (*ssh.Client, error) {
	client, err := ssh.Dial("tcp", hostname, config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// SSHConnection is a container for the pieces necessary to hold an SSH connection
// open in a goroutine
type SSHConnection struct {
	hostname string
	config   *ssh.ClientConfig
	client   *ssh.Client
	session  *ssh.Session
}

func (s *SSHConnection) init(hostname string, username string, keypath string) error {
	s.hostname = hostname

	// Instantiate config
	config, err := generateConfig(username, keypath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	s.config = config

	// Establish connection
	client, err := dialServer(hostname, config)
	if err != nil {
		log.Fatal(err)
		return err
	}
	s.client = client

	// Make a session
	session, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
		return err
	}
	s.session = session

	return nil
}

// TearDown closes all of the connections in an SSH connection
func (s *SSHConnection) TearDown() {
	s.session.Close()
	s.client.Close()
}

func (s *SSHConnection) executeCmd(cmd string) string {
	var stdoutBuf bytes.Buffer
	s.session.Stdout = &stdoutBuf
	s.session.Run(cmd)

	return stdoutBuf.String()
}

// InitSSHConnection takes host details and returns an SSHConnection
func InitSSHConnection(address string, username string, keypath string) (*SSHConnection, error) {
	sshConn := &SSHConnection{}
	err := sshConn.init(address, username, keypath)
	if err != nil {
		return nil, err
	}

	return sshConn, nil
}

// SSHResponse contains the result of running a command on a host via SSH
type SSHResponse struct {
	Host   string
	Result string
	Err    error
}

// MapCmd takes the details for a command and maps it, via SSH, across a list
// of hosts
func MapCmd(hosts hosts.Group, username string, keypath string, command string, results chan SSHResponse) {
	for _, host := range hosts {
		go func(host string) {
			response := SSHResponse{Host: host}

			sshConn, err := InitSSHConnection(host, username, keypath)
			defer sshConn.TearDown()

			if err != nil {
				fmt.Println(err.Error())
				response.Err = err
			} else {
				result := sshConn.executeCmd(command)
				response.Result = result
			}

			results <- response
		}(host)
	}
}
