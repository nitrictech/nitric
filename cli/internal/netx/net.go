// Copyright Nitric Pty Ltd.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package netx

import (
	"fmt"
	"net"
	"sync"
)

type getNextListenerOptions struct {
	minPort int
	maxPort int
}

type getNextListenerOption = func(opts *getNextListenerOptions)

func defaultGetNextListenerOptions() *getNextListenerOptions {
	// Defaults to IANA recommended ephemeral port range of 49152–65535
	return &getNextListenerOptions{
		minPort: 4000,
		maxPort: 6000,
	}
}

func MaxPort(maxPort int) getNextListenerOption {
	return func(opts *getNextListenerOptions) {
		opts.maxPort = maxPort
	}
}

func MinPort(minPort int) getNextListenerOption {
	return func(opts *getNextListenerOptions) {
		opts.minPort = minPort
	}
}

type ReservedPort int

var portReservationsLock = &sync.Mutex{}
var portReservations map[ReservedPort]struct{}

func (r ReservedPort) Release() {
	portReservationsLock.Lock()
	defer portReservationsLock.Unlock()

	if _, ok := portReservations[r]; ok {
		delete(portReservations, r)
	}
}

func GetNextPort(opts ...getNextListenerOption) (ReservedPort, error) {
	// reserve the port
	lis, err := GetNextListener(opts...)
	if err != nil {
		return 0, err
	}

	// close the listener to free the port
	defer func() {
		if err := lis.Close(); err != nil {
			fmt.Printf("Error closing listener: %v\n", err)
		}
	}()

	// return the port number
	port, ok := lis.Addr().(*net.TCPAddr)
	if !ok {
		return 0, fmt.Errorf("failed to get port from listener address")
	}

	return ReservedPort(port.Port), nil
}

func isPortReserved(port int) bool {
	if portReservations == nil {
		portReservations = make(map[ReservedPort]struct{})
	}

	// Check if the port is reserved and not expired
	if _, exists := portReservations[ReservedPort(port)]; exists {
		return true
	}

	return false
}

func reservePort(port int) {
	if portReservations == nil {
		portReservations = make(map[ReservedPort]struct{})
	}

	// Reserve the port with the current time
	portReservations[ReservedPort(port)] = struct{}{}
}

// GetNextListener - Gets the next available free port starting from a predefined minimum port
// Up to a pre-defined maximum port
func GetNextListener(opts ...getNextListenerOption) (net.Listener, error) {
	portReservationsLock.Lock()
	defer portReservationsLock.Unlock()

	// default and apply options
	// this allows the use of single or default options
	// without having to include parameters
	options := defaultGetNextListenerOptions()
	for _, opt := range opts {
		opt(options)
	}

	currentPort := options.minPort

	for currentPort < options.maxPort {
		if isPortReserved(currentPort) {
			currentPort = currentPort + 1
			continue
		}

		// attempt to get listener for port
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", currentPort))
		if err != nil {
			// increment the port and continue
			currentPort = currentPort + 1
			continue
		}

		// return the listener
		reservePort(currentPort)
		return lis, nil
	}

	return nil, fmt.Errorf("no ports available in range [%d-%d]", options.minPort, options.maxPort)
}

func GetInterfaceIpv4Addr(interfaceName string) (string, error) {
	ief, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return "", err
	}

	addrs, err := ief.Addrs()
	if err != nil {
		return "", err
	}

	var ipv4Addr net.IP
	for _, addr := range addrs {
		if ipv4Addr = addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
			break
		}
	}

	if ipv4Addr == nil {
		return "", fmt.Errorf("interface %s don't have an ipv4 address", interfaceName)
	}

	return ipv4Addr.String(), nil
}
