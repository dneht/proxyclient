/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package proxyclient

import "C"
import (
	"context"
	"errors"
	M "github.com/sagernet/sing/common/metadata"
	"github.com/sagernet/sing/protocol/http"
	"github.com/sagernet/sing/protocol/socks"
	"net"
)

func (p *Proxy) initHttp(ctx context.Context) error {
	user, passwd := "", ""
	if nil != p.option.User {
		user, passwd = p.option.User.Name, p.option.User.Passwd
	}
	client := http.NewClient(http.Options{
		Dialer:   p.dialer,
		Server:   p.addr,
		Username: user,
		Password: passwd,
	})
	p.http = client
	return nil
}

func (p *Proxy) singHttp(ctx context.Context, network string, addr M.Socksaddr) (net.Conn, error) {
	if nil == p.http {
		return nil, errors.New("invalid http client")
	}
	return p.http.DialContext(ctx, network, addr)
}

func (p *Proxy) initSocks(ctx context.Context, version socks.Version) error {
	user, passwd := "", ""
	if nil != p.option.User {
		user, passwd = p.option.User.Name, p.option.User.Passwd
	}
	client := socks.NewClient(p.dialer, p.addr, version, user, passwd)
	p.socks = client
	return nil
}

func (p *Proxy) singSocks(ctx context.Context, network string, addr M.Socksaddr) (net.Conn, error) {
	if nil == p.socks {
		return nil, errors.New("invalid socks client")
	}
	return p.socks.DialContext(ctx, network, addr)
}
