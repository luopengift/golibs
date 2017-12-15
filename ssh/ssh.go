package ssh

import (
	"fmt"
	//	"syscall"
	//	"github.com/luopengift/golibs/logger"
	"github.com/luopengift/types"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"time"
)

type Endpoint struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Ip       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Key      string `yaml:"key"`
}

type WindowSize struct {
	Width  int
	Height int
}

func NewEndpoint() *Endpoint {
	return new(Endpoint)
}

func NewEndpointWithValue(name, host, ip string, port int, user, password, key string) *Endpoint {
	return &Endpoint{
		Name:     name,
		Host:     host,
		Ip:       ip,
		Port:     port,
		User:     user,
		Password: password,
		Key:      key,
	}
}

func (ep *Endpoint) Init(filename string) error {
	return types.ParseConfigFile(filename, ep)
}

// 解析登录方式
func (ep *Endpoint) authMethods() ([]ssh.AuthMethod, error) {
	authMethods := []ssh.AuthMethod{
		ssh.Password(ep.Password),
	}

	if ep.Key == "" {
		return authMethods, nil
	}
	keyBytes, err := ioutil.ReadFile(ep.Key)
	if err != nil {
		return authMethods, err
	}
	// Create the Signer for this private key.
	var signer ssh.Signer
	if ep.Password == "" {
		signer, err = ssh.ParsePrivateKey(keyBytes)
	} else {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(keyBytes, []byte(ep.Password))
	}
	if err != nil {
		return authMethods, err
	}
	// Use the PublicKeys method for remote authentication.
	authMethods = append(authMethods, ssh.PublicKeys(signer))
	return authMethods, nil
}

func (ep *Endpoint) Address() string {
	addr := ""
	if ep.Host != "" {
		addr = fmt.Sprintf("%s:%d", ep.Host, ep.Port)
	} else {
		addr = ep.Ip + ":" + strconv.Itoa(ep.Port)
	}
	return addr
}

func (ep *Endpoint) InitSshClient() (*ssh.Client, error) {
	auths, err := ep.authMethods()

	if err != nil {
		return nil, fmt.Errorf("鉴权出错:", err)
	}

	config := &ssh.ClientConfig{
		User: ep.User,
		Auth: auths,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 5 * time.Second,
	}

	client, err := ssh.Dial("tcp", ep.Address(), config)
	if err != nil {
		return nil, fmt.Errorf("建立SSH连接出错:", err)
	}
	return client, nil
}

func (ep *Endpoint) Upload(src, dest string) ([]byte, error) {
	client, err := ep.InitSshClient()
	if err != nil {
		return nil, fmt.Errorf("建立SSH连接出错:", err)
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)

	if err != nil {
		return nil, fmt.Errorf("建立sftp出错:", err)
	}
	defer sftpClient.Close()
	srcFile, err := os.Open(src)
	if err != nil {
		return nil, fmt.Errorf("读取本地文件出错:", err)
	}
	defer srcFile.Close()

	destFile, err := sftpClient.Create(dest)
	if err != nil {
		return nil, fmt.Errorf("创建远程文件出错:", err)
	}
	defer destFile.Close()
	size := 0
	buf := make([]byte, 1024)
	for {
		n, err := srcFile.Read(buf)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("上传文件出错1:", err)
		}
		if n == 0 {
			break
		}
		if _, err := destFile.Write(buf[:n]); err != nil {
			return nil, fmt.Errorf("上传文件出错2:", err)
		}
		size += n
	}
	return []byte(fmt.Sprintf("文件上传成功,%dkb↑", size)), nil
}

func (ep *Endpoint) Download(src, dest string) ([]byte, error) {
	client, err := ep.InitSshClient()
	if err != nil {
		return nil, fmt.Errorf("建立SSH连接出错:", err)
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)

	if err != nil {
		return nil, fmt.Errorf("建立sftp出错:", err)
	}
	defer sftpClient.Close()

	srcFile, err := sftpClient.Open(src)
	if err != nil {
		return nil, fmt.Errorf("读取远程文件出错:", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return nil, fmt.Errorf("创建本地文件出错:", err)
	}
	defer destFile.Close()

	size := 0
	buf := make([]byte, 1024)
	for {
		n, err := srcFile.Read(buf)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("下载文件出错1:", err)
		}
		if n == 0 {
			break
		}
		if _, err := destFile.Write(buf[:n]); err != nil {
			return nil, fmt.Errorf("下载文件出错2:", err)
		}
		size += n
	}
	return []byte(fmt.Sprintf("文件下载成功,%dkb↓", size)), nil
}

func (ep *Endpoint) CmdOutBytes(cmd string) ([]byte, error) {
	auths, err := ep.authMethods()

	if err != nil {
		return nil, fmt.Errorf("鉴权出错:", err)
	}

	config := &ssh.ClientConfig{
		User: ep.User,
		Auth: auths,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 5 * time.Second,
	}

	client, err := ssh.Dial("tcp", ep.Address(), config)
	if err != nil {
		return nil, fmt.Errorf("建立连接出错:", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("创建Session出错:", err)
	}
	defer session.Close()
	return session.CombinedOutput(cmd)
}

func (ep *Endpoint) StartTerminal() error {
	auths, err := ep.authMethods()

	if err != nil {
		return fmt.Errorf("鉴权出错:", err)
	}

	config := &ssh.ClientConfig{
		User: ep.User,
		Auth: auths,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 5 * time.Second,
	}

	client, err := ssh.Dial("tcp", ep.Address(), config)
	if err != nil {
		return fmt.Errorf("建立连接出错:", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("创建Session出错:", err)
	}

	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("创建文件描述符出错:", err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	size := &WindowSize{}
	go func() error {
		t := time.NewTimer(time.Millisecond * 0)
		for {
			select {
			case <-t.C:
				size.Width, size.Height, err = terminal.GetSize(fd)
				if err != nil {
					return fmt.Errorf("获取窗口宽高出错:", err)
				}
				err = session.WindowChange(size.Height, size.Width)
				if err != nil {
					return fmt.Errorf("改变窗口大小出错:", err)
				}
				t.Reset(500 * time.Millisecond)
			}
		}
	}()
	defer terminal.Restore(fd, oldState)

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm-256color", size.Height, size.Width, modes); err != nil {
		return fmt.Errorf("创建终端出错:", err)
	}

	err = session.Shell()
	if err != nil {
		return fmt.Errorf("执行Shell出错:", err)
	}

	err = session.Wait()
	if err != nil {
		return fmt.Errorf("执行Wait出错:", err)
	}
	return nil
}
