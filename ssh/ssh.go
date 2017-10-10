package ssh

import (
    "golang.org/x/crypto/ssh"
    "golang.org/x/crypto/ssh/terminal"
    "io/ioutil"
    "os"
    "time"
    "net"
)

func client(user, addr, password, keyFile string) (*ssh.Client, error) {
    key, err := ioutil.ReadFile(keyFile)
    if err != nil {
        return nil, err
    }

    // Create the Signer for this private key.
    signer, err := ssh.ParsePrivateKey(key)
    if err != nil {
        return nil, err
    }

    config := &ssh.ClientConfig{
        User: user,
        Auth: []ssh.AuthMethod{
            ssh.Password(password),
            // Use the PublicKeys method for remote authentication.
            ssh.PublicKeys(signer),
        },
        HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
            return nil
        },
        Timeout: 5 * time.Second,
    }

    // Connect to the remote server and perform the SSH handshake.
    client, err := ssh.Dial("tcp", addr, config)
    if err != nil {
        return nil, err
    }
    return client, err
}

func Session(user, addr, password, key, cmd string) ([]byte, error) {
    client, err := client(user, addr, password, key)
    if err != nil {
        return nil, err
    }
    defer client.Close()
    session, err := client.NewSession()
    if err != nil {
        return nil, err
    }
    defer session.Close()
    return session.CombinedOutput(cmd)
}

func Term(user, addr, password, key, cmd string) (string,error) {
    client, err := client(user, addr, password, key)
    if err != nil {
        return "",err
    }
    defer client.Close()
    session, err := client.NewSession()
    if err != nil {
        return "",err
    }
    defer session.Close()

    fd := int(os.Stdin.Fd())
    oldState, err := terminal.MakeRaw(fd)
    if err != nil {
        return "",err
    }
    defer terminal.Restore(fd, oldState)

    termWidth, termHeight, err := terminal.GetSize(fd)
    if err != nil {
        return "",err
    }
    // Set up terminal modes
    modes := ssh.TerminalModes{
        ssh.ECHO:          1,     // enable echoing
        ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
        ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
    }

    // Request pseudo terminal
    if err := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
        return "",err
    }

    res, _ := session.Output(cmd)
    out := addr + "\n" + string(res) + "\n"
    return string(out),nil
}

