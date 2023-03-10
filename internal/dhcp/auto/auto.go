package dhcpauto

import (
	"fmt"

	"github.com/cyp0633/drcom-go/internal/util"
	"github.com/manifoldco/promptui"
	"go.uber.org/zap"
)

// Auto 自动生成 DHCP 配置文件
func Auto() {
	fmt.Printf("We'll go through some steps to generate configuration automatically.\n")
	_ = selectVersion()
}

// selectVersion 选择 Dr.com 客户端版本
func selectVersion() int {
	var result string
	for err := error(nil); ; {
		prompt := promptui.Select{
			Label: "Select Dr.com client version",
			Items: []string{"5.2D", "6.0D"},
		}
		_, result, err = prompt.Run()
		if err != nil {
			util.Logger.Error("Prompt failed", zap.Error(err))
		} else {
			break
		}
	}
	switch result {
	case "5.2D":
		return Drcom52D
	case "6.0D":
		return Drcom60D
	default:
		return -1
	}
}

// Drcom version enum
const (
	Drcom52D = iota
	Drcom60D = iota
)
