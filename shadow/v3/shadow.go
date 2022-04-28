/*
	MIT License

	Copyright (c) 2020 Operator Foundation

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE.
*/

// Package shadow provides a PT 2.1 Go API wrapper around the connections used by Shadowsocks
package shadow

import (
	"github.com/OperatorFoundation/go-shadowsocks2/darkstar"
	"net"
	"strconv"
	"strings"
)

//Config contains the necessary command like arguments to run shadow
type ClientConfig struct {
	Password   string `json:"password"`
	CipherName string `json:"cipherName"`
	Address    string `json:"address"`
}

type ServerConfig struct {
	Password   string `json:"password"`
	CipherName string `json:"cipherName"`
}

//Transport contains the arguments to be used with Optimizer
type Transport struct {
	Password   string
	CipherName string
	Address    string
}

type ShadowListener struct {
	Password string
	Address  string
	Listener net.Listener
}

func (s ShadowListener) Accept() (net.Conn, error) {
	addressArray := strings.Split(s.Address, ":")
	host := addressArray[0]
	port, stringErr := strconv.Atoi(addressArray[1])
	if stringErr != nil {
		return nil, stringErr
	}

	server := darkstar.NewDarkStarServer(s.Password, host, port)
	c, err := s.Listener.Accept()
	return server.StreamConn(c), err
}

func (s ShadowListener) Close() error {
	return s.Listener.Close()
}

func (s ShadowListener) Addr() net.Addr {
	return s.Listener.Addr()
}

//NewConfig is used to create a config for testing
func NewClientConfig(password string, cipherName string, address string) ClientConfig {
	return ClientConfig{
		Password:   password,
		CipherName: cipherName,
		Address:    address,
	}
}

func NewServerConfig(password string, cipherName string) ServerConfig {
	return ServerConfig{
		Password:   password,
		CipherName: cipherName,
	}
}

//NewTransport is used for creating a transport for Optimizer
func NewTransport(password string, cipherName string, address string) Transport {
	return Transport{
		Password:   password,
		CipherName: cipherName,
		Address:    address,
	}
}

//Listen checks for a working connection
func (config ServerConfig) Listen(address string) (net.Listener, error) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	shadowListener := ShadowListener{
		Password: config.Password,
		Address:  address,
		Listener: l,
	}

	return shadowListener, nil
}

//Dial connects to the address on the named network
func (config ClientConfig) Dial(address string) (net.Conn, error) {
	addressArray := strings.Split(address, ":")
	//portArray := strings.SplitAfter(address, ":")
	host := addressArray[0]
	port, stringErr := strconv.Atoi(addressArray[1])
	if stringErr != nil {
		return nil, stringErr
	}
	client := darkstar.NewDarkStarClient(config.Password, host, port)

	netConn, dialError := net.Dial("tcp", address)
	if dialError != nil {
		return nil, dialError
	}

	return client.StreamConn(netConn), nil
}

// Dial creates outgoing transport connection
func (transport *Transport) Dial() (net.Conn, error) {
	addressArray := strings.Split(transport.Address, ":")
	host := addressArray[0]
	port, stringErr := strconv.Atoi(addressArray[1])
	if stringErr != nil {
		return nil, stringErr
	}

	client := darkstar.NewDarkStarClient(transport.Password, host, port)
	netConn, dialError := net.Dial("tcp", transport.Address)
	if dialError != nil {
		return nil, dialError
	}

	return client.StreamConn(netConn), nil
}

func (transport *Transport) Listen() (net.Listener, error) {
	listener, err := net.Listen("tcp", transport.Address)
	if err != nil {
		return nil, err
	}

	shadowListener := ShadowListener{
		Password: transport.Password,
		Address:  transport.Address,
		Listener: listener,
	}

	return shadowListener, nil
}
