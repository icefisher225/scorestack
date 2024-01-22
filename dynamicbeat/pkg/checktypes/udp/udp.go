package udp

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/scorestack/scorestack/dynamicbeat/pkg/check"
)

type Definition struct {
	Config  check.Config
	IP      string `optiontype:"required"`
	Port    string `optiontype:"required"`
	Payload	string `optiontype:"required"`
	Content string `optiontype:"required"`
}

// Run a single intance of the check
func (d *Definition) Run(ctx context.Context) check.Result {
	// Initialize Empty Result
	result := check.Result{Timestamp: time.Now(), Metadata: d.Config.Metadata}

	// Setup for UDP Connection
	// Combine IP and Port
	addr := d.IP + ":" + d.Port

	// Convert IP address to UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		result.Message = fmt.Sprintf("Problem converting IP/port %s:%s to UDP address: %s", d.IP, d.Port, err)
	}

	// Dial to the address with UDP
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		result.Message = fmt.Sprintf("Problem dialing to %s:%s with UDP: %s", d.IP, d.Port, err)
	}

	// Send a message to the server
	_, err = conn.Write([]byte(d.Payload + "\n"))
	if err != nil {
		result.Message = fmt.Sprintf("Problem sending message to %s:%s with UDP: %s", d.IP, d.Port, err)
	}

	// Read from the connection untill a new line is send
	data, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		result.Message = fmt.Sprintf("Problem recieving data from %s:%s via UDP: %s", d.IP, d.Port, err)
	}

	// Check if returned data is equal to expected data
	if data == d.Content {
		result.Passed = true
		return result
	}

	// Check failed due to incorrect content recieved
	result.Message = fmt.Sprintf("Incorrect data received from %s:%s. Expected %s, got %s", d.IP, d.Port, d.Content, data)
	return result
}

// GetConfig returns the current CheckConfig struct this check has been initialized with
func (d *Definition) GetConfig() check.Config {
	return d.Config
}

// SetConfig reconfigures this check with a new CheckConfig struct
func (d *Definition) SetConfig(c check.Config) {
	d.Config = c
}
