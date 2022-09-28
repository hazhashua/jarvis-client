package utils

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SftpClient sftp客户端
type SftpClient struct {
	Session *sftp.Client
}

// NewSessionWithPassword 使用密码登录
func NewSessionWithPassword(host string, port int, user, password string) (sftpClient *SftpClient, err error) {
	sftpClient, err = NewSession(host, port, user, ssh.Password(password))
	return
}

// NewSessionWithKey 通过key登录sftp
func NewSessionWithKey(host string, port int, user, keyPath, keyPassword string) (sftpClient *SftpClient, err error) {
	fmt.Println(keyPath)
	keyData, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return
	}
	fmt.Println(string(keyData))
	var singer ssh.Signer
	if keyPassword == "" {
		singer, err = ssh.ParsePrivateKey(keyData)
	} else {
		singer, err = ssh.ParsePrivateKeyWithPassphrase(keyData, []byte(keyPassword))
	}
	if err != nil {
		return
	}
	sftpClient, err = NewSession(host, port, user, ssh.PublicKeys(singer))
	return
}

// NewSession 开启一个会话
func NewSession(host string, port int, user string, authMehtods ...ssh.AuthMethod) (sftpClient *SftpClient, err error) {
	clientConfig := &ssh.ClientConfig{
		User:    user,
		Auth:    authMehtods,
		Timeout: 4 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), clientConfig)
	if err != nil {
		return nil, err
	}
	sftpClient = &SftpClient{}
	sftpClient.Session, err = sftp.NewClient(sshClient)
	if err != nil {
		return nil, err
	}
	return
}

// ScopyRmoteFile 拷贝远程文件
func (sftpClient *SftpClient) ScopyRmoteFile(remoteFilePath, localFilePath string) (err error) {
	srcFile, err := sftpClient.Session.Open(remoteFilePath)
	if err != nil {
		return
	}
	defer srcFile.Close()
	dstFile, err := os.Create(localFilePath)
	if err != nil {
		return
	}
	defer dstFile.Close()
	if _, err = srcFile.WriteTo(dstFile); err != nil {
		return
	}
	return
}

// SendFile 发送文件到远程
func (sftpClient *SftpClient) SendFile(localFilePath, remoteFilePath string) (err error) {
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		return
	}
	defer srcFile.Close()
	dstFile, err := sftpClient.Session.Create(localFilePath)
	if err != nil {
		return
	}
	defer dstFile.Close()
	buf := make([]byte, 1024)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf)
	}
	return
}

// Close close session
func (sftpClient *SftpClient) Close() (err error) {
	err = sftpClient.Session.Close()
	return
}
