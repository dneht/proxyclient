/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package proxyclient

import (
	"context"
	"errors"
	shadowsocks "github.com/sagernet/sing-shadowsocks2"
	"github.com/sagernet/sing/common/bufio"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
	"net"
)

func (o *SSOption) check() bool {
	return o.Method != "" && o.Password != ""
}

func (p *Proxy) initShadowsocks(ctx context.Context) error {
	opt := p.option.SS
	if nil == opt || !opt.check() {
		return errors.New("invalid shadowsocks option")
	}
	method, err := shadowsocks.CreateMethod(ctx, opt.Method, shadowsocks.MethodOptions{
		Password: opt.Password,
	})
	if err != nil {
		return err
	}
	p.ss = method
	p.mux, err = NewMux(&ProxyDialer{p}, p.logger, p.option.Mux)
	if err != nil {
		return err
	}
	return nil
}

func (p *Proxy) singShadowsocks(ctx context.Context, network string, addr M.Socksaddr) (net.Conn, error) {
	if nil == p.ss {
		return nil, errors.New("invalid shadowsocks client")
	}
	connect, err := p.dialer.DialContext(ctx, N.NetworkTCP, p.addr)
	if err != nil {
		return nil, err
	}
	switch network {
	case N.NetworkTCP:
		return p.ss.DialEarlyConn(connect, addr), nil
	case N.NetworkUDP:
		return bufio.NewBindPacketConn(p.ss.DialPacketConn(connect), addr), nil
	default:
		return nil, errors.New("unknown network " + network)
	}
}
