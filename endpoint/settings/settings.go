package settings

import "net"

type Settings struct {
	Port  int `extractor:"-"`
	Peers []net.TCPAddr
}
