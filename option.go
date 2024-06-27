/*
Copyright (c) 2020, Dash

Licensed under the LGPL, Version 3.0 (the "License");
you may not use this file except in compliance with the License.
*/

package proxyclient

type Option struct {
	Mux   *MuxOption
	User  *UserOption
	SSH   *SSHOption
	SS    *SSOption
	VMess *VMessOption
}

type MuxOption struct {
	Enabled        bool           `json:"enabled"`
	Protocol       string         `json:"protocol,omitempty"`
	MaxConnections int            `json:"max_connections,omitempty"`
	MinStreams     int            `json:"min_streams,omitempty"`
	MaxStreams     int            `json:"max_streams,omitempty"`
	Padding        bool           `json:"padding,omitempty"`
	Brutal         *BrutalOptions `json:"brutal,omitempty"`
}

type BrutalOptions struct {
	Enabled  bool `json:"enabled,omitempty"`
	UpMbps   int  `json:"up_mbps,omitempty"`
	DownMbps int  `json:"down_mbps,omitempty"`
}

type UserOption struct {
	Name   string `json:"name,omitempty"`
	Passwd string `json:"passwd,omitempty"`
}

type SSHOption struct {
	Name     string `json:"name"`
	Passwd   string `json:"passwd,omitempty"`
	PkFile   string `json:"pk_file,omitempty"`
	PkPasswd string `json:"pk_passwd,omitempty"`
}

type SSOption struct {
	Method   string `json:"method"`
	Password string `json:"password"`
}

type VMessOption struct {
	UUID                string `json:"uuid"`
	Security            string `json:"security"`
	AlterId             int    `json:"alter_id,omitempty"`
	GlobalPadding       bool   `json:"global_padding,omitempty"`
	AuthenticatedLength bool   `json:"authenticated_length,omitempty"`
}
