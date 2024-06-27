/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package proxyclient

import (
	"context"
	"errors"
	vmess "github.com/sagernet/sing-vmess"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
	"github.com/sagernet/sing/common/ntp"
	"net"
)

func (o *VMessOption) check() bool {
	if o.Security == "" {
		o.Security = "auto"
	}
	return o.UUID != ""
}

func (p *Proxy) initVMess(ctx context.Context) error {
	opt := p.option.VMess
	if nil == opt || !opt.check() {
		return errors.New("invalid vmess option")
	}
	options := make([]vmess.ClientOption, 0, 4)
	if timeFunc := ntp.TimeFuncFromContext(ctx); timeFunc != nil {
		options = append(options, vmess.ClientWithTimeFunc(timeFunc))
	}
	if opt.GlobalPadding {
		options = append(options, vmess.ClientWithGlobalPadding())
	}
	if opt.AuthenticatedLength {
		options = append(options, vmess.ClientWithAuthenticatedLength())
	}
	client, err := vmess.NewClient(opt.UUID, opt.Security, opt.AlterId, options...)
	if err != nil {
		return err
	}
	p.vmess = client
	p.mux, err = NewMux(&ProxyDialer{p}, p.logger, p.option.Mux)
	if err != nil {
		return err
	}
	return nil
}

func (p *Proxy) singVMess(ctx context.Context, network string, addr M.Socksaddr) (net.Conn, error) {
	if nil == p.vmess {
		return nil, errors.New("invalid vmess client")
	}
	connect, err := p.dialer.DialContext(ctx, N.NetworkTCP, p.addr)
	if err != nil {
		return nil, err
	}
	switch network {
	case N.NetworkTCP:
		return p.vmess.DialEarlyConn(connect, addr), nil
	case N.NetworkUDP:
		return p.vmess.DialEarlyPacketConn(connect, addr), nil
	default:
		return nil, errors.New("unknown network " + network)
	}
}
