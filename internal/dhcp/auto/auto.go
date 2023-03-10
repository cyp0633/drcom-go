package dhcpauto

import (
	"github.com/cyp0633/drcom-go/internal/util"
	"github.com/manifoldco/promptui"
	"go.uber.org/zap"
)

// Auto 自动生成 DHCP 配置文件
func Auto() {
	util.Logger.Info("We'll go through some steps to generate configuration automatically." +
		"Please note that this is not a well-developed function." +
		"The first thing is to input your credentials.")
	ver := selectVersion()
	util.Conf.Username = inputAccount()
	util.Conf.Password = inputPassword()

	util.Logger.Info("Probing auth server...")
	sendProbe()
	recvProbe()

	// Host part
	util.Logger.Info("Gathering host information...")
	getHostInfo()

	// Mechanism part
	util.Logger.Info("Analyzing auth mechanism...")
	guess(ver)
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
			util.Logger.Error("Select version failed", zap.Error(err))
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

// inputAccount 输入账号
func inputAccount() string {
	var result string
	for err := error(nil); ; {
		prompt := promptui.Prompt{
			Label:    "Username",
			Validate: func(input string) error { return nil },
		}
		result, err = prompt.Run()
		if err != nil {
			util.Logger.Error("Input username prompt failed", zap.Error(err))
		} else {
			break
		}
	}
	return result
}

// inputPassword 输入密码
func inputPassword() string {
	var result string
	for err := error(nil); ; {
		prompt := promptui.Prompt{
			Label:    "Password",
			Validate: func(input string) error { return nil },
			Mask:     '*',
		}
		result, err = prompt.Run()
		if err != nil {
			util.Logger.Error("Input password prompt failed", zap.Error(err))
		} else {
			break
		}
	}
	return result
}
