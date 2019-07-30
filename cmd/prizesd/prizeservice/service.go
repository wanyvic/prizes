package prizeservice

import (
	"bytes"
	"crypto/rand"
	"math/big"
	mathRand "math/rand"
	"net"
	"time"
)

func CreateRandomNumberString(len int) string {
	var container string
	var str = "1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}

func CreateRandomString(len int) string {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyz1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}
func GetFreeIp() net.IP {
	mathRand.Seed(time.Now().UnixNano())
	int1 := mathRand.Intn(254)
	mathRand.Seed(time.Now().UnixNano())
	int2 := mathRand.Intn(255)

	var bytes [4]byte
	bytes[0] = byte((int1) & 0xFF)
	bytes[1] = byte((int2) & 0xFF)
	bytes[2] = byte(0x00)
	bytes[3] = byte(0x0A)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}
