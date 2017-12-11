package ssh

import (
	"github.com/luopengift/types"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
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
	return &Endpoint{}
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
		addr = ep.Host + ":" + strconv.Itoa(ep.Port)
	} else {
		addr = ep.Ip + ":" + strconv.Itoa(ep.Port)
	}
	return addr
}

func (ep *Endpoint) Session() (*ssh.Session, error) {
	authMethods, err := ep.authMethods()
	if err != nil {
		return nil, err
	}
	config := &ssh.ClientConfig{
		User: ep.User,
		Auth: authMethods,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 5 * time.Second,
	}

	// Connect to the remote server and perform the SSH handshake.
	client, err := ssh.Dial("tcp", ep.Address(), config)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (ep *Endpoint) CmdOutBytes(cmd string) ([]byte, error) {
	session, err := ep.Session()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	return session.CombinedOutput(cmd)
}

func (ep *Endpoint) StartTerminal() error {
	session, err := ep.Session()
	if err != nil {
		return err
	}
	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer terminal.Restore(fd, oldState)

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	size := WindowSize{}

	go func() error {
		t := time.NewTimer(time.Second * 0)
		var err error
		for {
			select {
			case <-t.C:
				size.Width, size.Height, err = terminal.GetSize(fd)
				if err != nil {
					return err
				}
				err = session.WindowChange(size.Height, size.Width)
				if err != nil {
					return err
				}
				t.Reset(time.Second * 1)
			}
		}
	}()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1, //显示输入的命令
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err = session.RequestPty("xterm-256color", size.Height, size.Width, modes); err != nil {
		return err
	}
	if err = session.Shell(); err != nil {
		return err
	}
	if err = session.Wait(); err != nil {
		return err
	}
	return nil

}

