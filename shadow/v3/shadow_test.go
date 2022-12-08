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

package shadow

import (
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/kataras/golog"
)

const data = "test"

func TestMain(m *testing.M) {
	logDir := "/Users/bluesaxorcist/Desktop"
	config := NewServerConfig("127.0.0.1:1236", "RaHouPFVOazVSqInoMm8BSO9o/7J493y4cUVofmwXAU=", "darkstar", "Shadow", &logDir)
	listener, listenErr := config.Listen()
	if listenErr != nil {
		return
	}
	go acceptConnections(listener)

	os.Exit(m.Run())
}

func acceptConnections(listener net.Listener) {
	serverBuffer := make([]byte, 4)
	for {
		serverConn, err := listener.Accept()
		if err != nil {
			return
		}
		go func() {
			//read on server side
			_, serverReadErr := serverConn.Read(serverBuffer)
			if serverReadErr != nil {
				return
			}

			//write data from serverConn for client to read
			_, serverWriteErr := serverConn.Write([]byte(data))
			if serverWriteErr != nil {
				return
			}
		}()
	}
}

func TestShadow(t *testing.T) {
	logDir := "/Users/bluesaxorcist/Desktop"
	config := NewClientConfig("127.0.0.1:1236", "6LukZ8KqZLQ7eOdaTVFkBVqMA8NS1AUxwqG17L/kHnQ=", "darkstar", "Shadow", &logDir)

	//create client buffer
	clientBuffer := make([]byte, 4)

	//call dial on client and check error
	clientConn, dialErr := config.Dial()
	if dialErr != nil {
		fmt.Println("clientConn Dial error")
		t.Fail()
		return
	}

	//write data from clientConn for server to read
	_, clientWriteErr := clientConn.Write([]byte(data))
	if clientWriteErr != nil {
		fmt.Println("client write error")
		t.Fail()
		return
	}

	//read on client side
	_, clientReadErr := clientConn.Read(clientBuffer)
	if clientReadErr != nil {
		fmt.Println("client read error")
		t.Fail()
		return
	}
}
func TestShadowFactory(t *testing.T) {
	MakeLog()
	logDir := "/Users/bluesaxorcist/Desktop"
	transport := NewTransport("127.0.0.1:1236", "6LukZ8KqZLQ7eOdaTVFkBVqMA8NS1AUxwqG17L/kHnQ=", "darkstar", &logDir)

	//create client buffer
	clientBuffer := make([]byte, 4)

	//call dial on client and check error
	clientConn, dialErr := transport.Dial()
	if dialErr != nil {
		fmt.Println("clientConn Dial error")
		t.Fail()
		return
	}

	//write data from clientConn for server to read
	_, clientWriteErr := clientConn.Write([]byte(data))
	if clientWriteErr != nil {
		fmt.Println("client write error")
		t.Fail()
		return
	}

	//read on client side
	_, clientReadErr := clientConn.Read(clientBuffer)
	if clientReadErr != nil {
		fmt.Println("client read error")
		t.Fail()
		return
	}
}

func MakeLog() {
	golog.SetLevel("debug")
	golog.SetOutput(os.Stderr)
}
