package dhcp

// 登录
func login() (tail []byte, salt []byte, err error) {
	salt, err = challenge()
	return
}

// 发送 challenge 并获得 salt
func challenge() (salt []byte, err error) {
	return
}
