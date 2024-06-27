/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package proxyclient

import (
	"context"
	"errors"
	mux "github.com/sagernet/sing-mux"
	shadowsocks "github.com/sagernet/sing-shadowsocks2"
	vmess "github.com/sagernet/sing-vmess"
	L "github.com/sagernet/sing/common/logger"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
	"github.com/sagernet/sing/protocol/http"
	"github.com/sagernet/sing/protocol/socks"
	"golang.org/x/crypto/ssh"
	"net"
	"strings"
)

type Proxy struct {
	scheme string
	option *Option
	logger L.Logger
	dialer N.Dialer
	addr   M.Socksaddr
	http   *http.Client
	socks  *socks.Client
	ssh    *ssh.Client
	ss     shadowsocks.Method
	vmess  *vmess.Client
	mux    *mux.Client
}

type ProxyDialer struct {
	proxy *Proxy
}

func New(url string, option *Option) (*Proxy, error) {
	return newProxy(url, nil, nil, option)
}

func NewWithLogger(url string, log L.Logger, option *Option) (*Proxy, error) {
	return newProxy(url, log, nil, option)
}

func newProxy(url string, log L.Logger, dialer N.Dialer, option *Option) (*Proxy, error) {
	p := Proxy{
		option: option,
	}
	if nil == log {
		p.logger = L.NOP()
	}
	if nil == dialer {
		p.dialer = N.SystemDialer
	}
	idx := strings.Index(url, "://")
	if idx <= 0 {
		return nil, errors.New("invalid server scheme")
	}
	ctx, addr := context.Background(), url[idx+3:]
	if len(addr) < 3 {
		return nil, errors.New("invalid server address")
	}
	p.scheme, p.addr = strings.Clone(url[0:idx]), M.ParseSocksaddr(url[idx+3:])
	switch p.scheme {
	case SchemeHTTP, SchemeHTTPS:
		return &p, p.initHttp(ctx)
	case SchemeSocks, SchemeSocks5:
		return &p, p.initSocks(ctx, socks.Version5)
	case SchemeSocks4:
		return &p, p.initSocks(ctx, socks.Version4)
	case SchemeSocks4a:
		return &p, p.initSocks(ctx, socks.Version4A)
	case SchemeSSH:
		return &p, p.initSsh(ctx)
	case SchemeSS:
		return &p, p.initShadowsocks(ctx)
	case SchemeVMess:
		return &p, p.initVMess(ctx)
	default:
		return nil, errors.New("unsupported scheme " + p.scheme)
	}
}

func (p *Proxy) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	return p.dialContext(ctx, network, M.ParseSocksaddr(addr))
}

func (p *Proxy) dialContext(ctx context.Context, network string, addr M.Socksaddr) (net.Conn, error) {
	switch p.scheme {
	case SchemeHTTP, SchemeHTTPS:
		return p.singHttp(ctx, network, addr)
	case SchemeSocks, SchemeSocks5, SchemeSocks4, SchemeSocks4a:
		return p.singSocks(ctx, network, addr)
	case SchemeSSH:
		return p.singSsh(ctx, network, addr)
	case SchemeSS:
		return p.singShadowsocks(ctx, network, addr)
	case SchemeVMess:
		return p.singVMess(ctx, network, addr)
	default:
		return nil, errors.New("unsupported scheme " + p.scheme)
	}
}

func (d *ProxyDialer) DialContext(ctx context.Context, network string, dest M.Socksaddr) (net.Conn, error) {
	return d.proxy.dialContext(ctx, network, dest)
}

func (d *ProxyDialer) ListenPacket(ctx context.Context, dest M.Socksaddr) (net.PacketConn, error) {
	return nil, errors.New("unimplemented listen")
}
