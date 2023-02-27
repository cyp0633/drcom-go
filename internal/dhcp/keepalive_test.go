package dhcp

import (
	"encoding/hex"
	"testing"
)

func TestGenFilePacket(t *testing.T) {
	pkt:=genKeepalive2Packet(true,1,0)
	t.Log(hex.EncodeToString(pkt))
	if(len(pkt)!=40) {
		t.Error("len(pkt)!=40")
	}
}
