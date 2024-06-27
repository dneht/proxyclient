/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package proxyclient

import (
	mux "github.com/sagernet/sing-mux"
	L "github.com/sagernet/sing/common/logger"
	N "github.com/sagernet/sing/common/network"
	"github.com/sagernet/sing/common/uot"
)

func NewMux(dialer N.Dialer, logger L.Logger, option *MuxOption) (*mux.Client, error) {
	if nil == option || !option.Enabled {
		return nil, nil
	}
	brutal := mux.BrutalOptions{
		Enabled:    option.Brutal.Enabled,
		SendBPS:    uint64(option.Brutal.UpMbps * MbpsToBps),
		ReceiveBPS: uint64(option.Brutal.DownMbps * MbpsToBps),
	}
	if brutal.SendBPS < mux.BrutalMinSpeedBPS {
		brutal.SendBPS = mux.BrutalMinSpeedBPS
	}
	if brutal.ReceiveBPS < mux.BrutalMinSpeedBPS {
		brutal.ReceiveBPS = mux.BrutalMinSpeedBPS
	}
	return mux.NewClient(mux.Options{
		Dialer:         dialer,
		Logger:         logger,
		Protocol:       option.Protocol,
		MaxConnections: option.MaxConnections,
		MinStreams:     option.MinStreams,
		MaxStreams:     option.MaxStreams,
		Padding:        option.Padding,
		Brutal:         brutal,
	})
}

func NewUot(dialer N.Dialer, version int) *uot.Client {
	return &uot.Client{
		Dialer:  dialer,
		Version: uint8(version),
	}
}
