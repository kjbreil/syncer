package settings

import "net"

// Settings contains the configuration for the server.
type Settings struct {
	// Port is the port the server listens on.
	Port int `json:"port"`
	// Peers is a list of peers the server connects to.
	Peers []net.TCPAddr `json:"peers"`
	// AutoUpdate determines if the server should update itself automatically.
	AutoUpdate bool `json:"auto_update"`
}
