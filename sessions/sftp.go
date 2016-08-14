package sessions

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jmsdnns/omnitool/hosts"
	"github.com/pkg/sftp"
)

// SFTPConnnection is a container for the pieces necessary to hold an SFTP connection
// open in a goroutine
type SFTPConnnection struct {
	ssh    *SSHConnection
	client *sftp.Client
	error  error
}

func (s *SFTPConnnection) init(ss *SSHConnection) error {
	s.ssh = ss

	client, err := sftp.NewClient(ss.client)
	if err != nil {
		log.Fatal(err)
		return err
	}
	s.client = client

	return nil
}

// TearDown closes all of the connections in an SFTP connection
func (s *SFTPConnnection) TearDown() {
	s.client.Close()
}

func (s *SFTPConnnection) copyFile(localPath string, remotePath string) error {
	// Create remote file handle
	f, err := s.client.Create(remotePath)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// Open input file
	fi, err := os.Open(localPath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer fi.Close()

	fileChunk := make([]byte, 1024)
	for {
		// read a chunk
		n, err := fi.Read(fileChunk)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if n == 0 {
			break
		}

		// Write data to file handle
		if _, err := f.Write(fileChunk); err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	return nil
}

// InitSFTPConnection returns an SFTP client from an existing SSHSession
func InitSFTPConnection(s *SSHConnection) (*SFTPConnnection, error) {
	sftpConn := &SFTPConnnection{}

	err := sftpConn.init(s)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return sftpConn, nil
}

// SFTPResponse contains the result of copying a file to a host via SFTP
type SFTPResponse struct {
	Host   string
	Result string
	Err    error
}

// MapCopy takes the details for a file transfer and maps the file transfer
// across a list of hosts
func MapCopy(hosts hosts.Group, username string, keypath string, localPath string, remotePath string, results chan SFTPResponse) {
	for _, host := range hosts {
		go func(host string) {
			response := SFTPResponse{Host: host}

			sshConn, err := InitSSHConnection(host, username, keypath)
			if err != nil {
				log.Fatal(err)
				response.Err = err
			}
			sftpConn, err := InitSFTPConnection(sshConn)
			defer sshConn.TearDown()
			defer sftpConn.TearDown()

			err = sftpConn.copyFile(localPath, remotePath)
			response.Result = "ok"
			if err != nil {
				log.Fatal(err)
				response.Result = "error"
				response.Err = err
			}

			results <- response
		}(host)
	}
}
