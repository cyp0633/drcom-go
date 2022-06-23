package network

import (
	"github.com/cyp0633/go-drcom/util"
	"log"
	"net"
)

var IfName = "eth0" // interface name

// TrySmartEaplogin try to authenticate via 802.11x
func TrySmartEaplogin() {
	IfList, err := net.Interfaces()
	if err != nil {
		log.Fatalln("Failed to acquire net interface list,", err.Error())
	}
	for i, ifs := range IfList {
		eapLogin(ifs)
	}
}

// eapLogin authenticate via EAP
func eapLogin(ifs net.Interface) {
	log.Printf("Use user %s to login...\n", util.Conf.Username)
	log.Printf("[EAP:0] Initialize interface...")
	log.Fatalln("802.11x not implemented")
}

// eapolInit initialize a buffer and socket
func eapolInit()
