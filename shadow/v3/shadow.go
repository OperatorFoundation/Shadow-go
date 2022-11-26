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
	"errors"
	"net"
	"strconv"
	"strings"

	"github.com/OperatorFoundation/locket-go"
	"github.com/OperatorFoundation/go-shadowsocks2/darkstar"
)

type ClientConfig struct {
	ServerAddress   string  `json:"serverAddress"`  
	ServerPublicKey string  `json:"serverPublicKey"`
	CipherName      string  `json:"cipherName"`     
	Transport       string  `json:"transport"`
	LogDir     		*string `json:"logDir"`
}

type ServerConfig struct {
	ServerAddress    string  `json:"serverAddress"`   
	ServerPrivateKey string  `json:"serverPrivateKey"`
	CipherName       string  `json:"cipherName"`      
	Transport        string  `json:"transport"`  
	LogDir    		 *string `json:"logDir"`
}

type Transport struct {
	ServerAddress string
	ServerKey  	  string
	CipherName 	  string
	LogDir    	  *string
}

type ShadowListener struct {
	Address  		 string
	ServerPrivateKey string
	CipherName		 string
	Listener 		 net.Listener
	LogDir   		 *string
}

func (s ShadowListener) Accept() (net.Conn, error) {
	addressArray := strings.Split(s.Address, ":")
	host := addressArray[0]
	port, stringErr := strconv.Atoi(addressArray[1])
	if stringErr != nil {
		return nil, stringErr
	}

	c, err := s.Listener.Accept()
	if err != nil {
		return nil, err
	}

	if s.LogDir != nil {
		c, err = locketgo.NewLocketConn(c, *s.LogDir, "ShadowServer")
		if err != nil {
			return nil, err
		}
	}

	if s.CipherName == "darkstar" {
	server := darkstar.NewDarkStarServer(s.ServerPrivateKey, host, port)
	
	return server.StreamConn(c)
	} else {
		return nil, errors.New("invalid cipher name")
	}
}

func (s ShadowListener) Close() error {
	return s.Listener.Close()
}

func (s ShadowListener) Addr() net.Addr {
	return s.Listener.Addr()
}

func NewClientConfig(serverAddress string, serverPublicKey string, cipherName string, transport string, logDir *string) ClientConfig {
	return ClientConfig{
		ServerAddress: 	 serverAddress,
		ServerPublicKey: serverPublicKey,
		CipherName:    	 cipherName,
		Transport: 		 transport,
		LogDir:     	 logDir,
	}
}

func NewServerConfig(serverAddress string, serverPrivateKey string, cipherName string, transport string, logDir *string) ServerConfig {
	return ServerConfig{
		ServerAddress: 	  serverAddress,
		ServerPrivateKey: serverPrivateKey,
		CipherName:    	  cipherName,
		Transport: 		  transport,
		LogDir:     	  logDir,
	}
}

func NewTransport(serverAddress string, serverKey string, cipherName string, logDir *string) Transport {
	return Transport{
		ServerAddress: serverAddress,
		ServerKey:     serverKey,
		CipherName:	   cipherName,
		LogDir:        logDir,
	}
}

// Listen checks for a working connection
func (config ServerConfig) Listen() (net.Listener, error) {
	// Verify the transport name on the config
	if config.Transport != "shadow" {
		return nil, errors.New("incorrect transport name")
	}

	l, err := net.Listen("tcp", config.ServerAddress)
	if err != nil {
		return nil, err
	}

	shadowListener := ShadowListener{
		Address:  		  config.ServerAddress,
		ServerPrivateKey: config.ServerPrivateKey,
		CipherName:		  config.CipherName,
		Listener: 		  l,
		LogDir:   		  config.LogDir,
	}

	return shadowListener, nil
}

// Dial connects to the server and returns a DarkStar connection if the handshake was successful
func (config ClientConfig) Dial() (net.Conn, error) {
	// Verify the transport name on the config
	if config.Transport != "shadow" {
		return nil, errors.New("incorrect transport name")
	}

	// Get a host and port from the provided address string
	addressArray := strings.Split(config.ServerAddress, ":")
	host := addressArray[0]
	port, stringErr := strconv.Atoi(addressArray[1])
	if stringErr != nil {
		return nil, stringErr
	}

	// Create a network connection
	netConn, dialError := net.Dial("tcp", config.ServerAddress)
	if dialError != nil {
		return nil, dialError
	}

	if config.LogDir != nil {
		netConn, dialError = locketgo.NewLocketConn(netConn, *config.LogDir, "ShadowClient")
		if dialError != nil {
			return nil, dialError
		}
	}

	if config.CipherName == "darkstar" {
	// Create a new  DarkStarClient
	darkStarClient := darkstar.NewDarkStarClient(config.ServerPublicKey, host, port)
	
	// Attempts to connect with the server and complete a handshake
	// If the handshake is successful, returns a DarkStar connection
	return darkStarClient.StreamConn(netConn)
	} else {
		return nil, errors.New("invalid cipher name")
	}
}

// Dial connects to the server and returns a DarkStar connection if the handshake was successful
func (transport *Transport) Dial() (net.Conn, error) {

	// Get a host and port from the transport address string
	addressArray := strings.Split(transport.ServerAddress, ":")
	host := addressArray[0]
	port, stringErr := strconv.Atoi(addressArray[1])
	if stringErr != nil {
		return nil, stringErr
	}

	// Create a new  DarkStarClient
	darkStarClient := darkstar.NewDarkStarClient(transport.ServerKey, host, port)
	if darkStarClient == nil {
		return nil, errors.New("failed to create a DarkStarClient with the provided password")
	}

	// Create a network connection
	netConn, dialError := net.Dial("tcp", transport.ServerAddress)
	if dialError != nil {
		return nil, dialError
	}

	if transport.LogDir != nil {
		netConn, dialError = locketgo.NewLocketConn(netConn, *transport.LogDir, "ShadowClient")
		if dialError != nil {
			return nil, dialError
		}
	}

	// Attempts to connect with the server and complete a handshake
	// If the handshake is successful, returns a DarkStar connection
	return darkStarClient.StreamConn(netConn)
}

func (transport *Transport) Listen() (net.Listener, error) {
	listener, err := net.Listen("tcp", transport.ServerAddress)
	if err != nil {
		return nil, err
	}

	shadowListener := ShadowListener{
		Address:  		  transport.ServerAddress,
		ServerPrivateKey: transport.ServerKey,
		CipherName:		  transport.CipherName,
		Listener: 		  listener,
		LogDir:   		  transport.LogDir,
	}

	return shadowListener, nil
}
