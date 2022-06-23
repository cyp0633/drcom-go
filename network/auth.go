package network

import (
	"fmt"
	"github.com/cyp0633/go-drcom/util"
)

const BindPort = 61440

func Drcom(tryTimes int) {
	if util.Opt.Verbose {
		fmt.Printf("You are binding at %s", "test")
	}
}
