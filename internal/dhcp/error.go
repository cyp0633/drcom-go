package dhcp

import "errors"

// 登录失败
var ErrorLogin = errors.New("login failed")

// Challenge 失败
var ErrorChallenge = errors.New("challenge failed")
