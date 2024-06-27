/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package proxyclient

import (
	"context"
	"errors"
	M "github.com/sagernet/sing/common/metadata"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
)

func (o *SSHOption) check() bool {
	return "" != o.Name
}

func (p *Proxy) initSsh(ctx context.Context) error {
	opt := p.option.SSH
	if nil == opt || !opt.check() {
		return errors.New("invalid ssh option")
	}
	auths, err := sshAuthMethod(opt.Passwd, opt.PkFile, opt.PkPasswd)
	if err != nil {
		return err
	}
	config := &ssh.ClientConfig{
		User: opt.Name,
		Auth: auths,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	p.ssh, err = ssh.Dial("tcp", p.addr.String(), config)
	if err != nil {
		return err
	}
	return nil
}

func (p *Proxy) singSsh(ctx context.Context, network string, addr M.Socksaddr) (net.Conn, error) {
	if nil == p.ssh {
		return nil, errors.New("invalid ssh client")
	}
	return p.ssh.DialContext(ctx, network, addr.String())
}

func sshAuthMethod(passwd, pkFile, pkPasswd string) ([]ssh.AuthMethod, error) {
	auths := make([]ssh.AuthMethod, 0, 4)
	if pkFile != "" {
		auth, err := sshPrivateKeyMethod(pkFile, pkPasswd)
		if err != nil {
			return nil, err
		}
		auths = append(auths, auth)
		return auths, nil
	}
	if passwd != "" {
		auths = append(auths, sshPasswordMethod(passwd))
		return auths, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	if auth, err := sshPrivateKeyMethod(home+"/.ssh/id_rsa", pkPasswd); nil == err {
		auths = append(auths, auth)
	}
	if auth, err := sshPrivateKeyMethod(home+"/.ssh/id_ed25519", pkPasswd); nil == err {
		auths = append(auths, auth)
	}
	return auths, nil
}

func sshPrivateKeyMethod(pkFile, pkPasswd string) (ssh.AuthMethod, error) {
	pkData, err := os.ReadFile(pkFile)
	if err != nil {
		return nil, err
	}
	var pk ssh.Signer
	if pkPasswd == "" {
		pk, err = ssh.ParsePrivateKey(pkData)
		if err != nil {
			return nil, err
		}
	} else {
		pk, err = ssh.ParsePrivateKeyWithPassphrase(pkData, []byte(pkPasswd))
		if err != nil {
			return nil, err
		}
	}
	return ssh.PublicKeys(pk), nil
}

func sshPasswordMethod(passwd string) ssh.AuthMethod {
	return ssh.Password(passwd)
}
