package util

import "encoding/binary"

// 进行 CRC 加密
//
// 来自于 drcom-generic 的 checksum 函数
func Checksum(s []byte) []byte {
	ret := 1234
	for i := 0; i < len(s); i += 4 {
		ret ^= int(binary.BigEndian.Uint32(s[i : i+4]))
	}
	ret = (1968 * ret) & 0xffffffff
	retBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(retBytes, uint32(ret))
	return retBytes
}

// 位移加密
//
// 来自于 drcom-generic 的 ror 函数
func Ror(md5, pwd []byte) []byte {
	ret := make([]byte, 0)
	for i := 0; i < len(pwd); i++ {
		x := md5[i] ^ pwd[i]
		ret = append(ret, (x<<3)&0xff+(x>>5))
	}
	return ret
}