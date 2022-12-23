package dhcp

import "errors"

// 登录失败
var ErrorLogin = errors.New("login failed")

// Challenge 失败
var ErrorChallenge = errors.New("challenge failed")

// 生成登录数据包失败
var ErrorGenLoginPacket = errors.New("generate login packet failed")

// 心跳包 1 失败
var ErrorKeepalive1 = errors.New("keepalive 1 failed")

// 心跳包 2 失败
var ErrorKeepalive2 = errors.New("keepalive 2 failed")