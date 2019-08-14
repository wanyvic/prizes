package massgrid

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/massgrid/btcd/chaincfg"
	"github.com/massgrid/btcutil"
	"github.com/wanyvic/gobtclib/client"
	"github.com/wanyvic/prizes/api/types/order"
)

var (
	DefaultNetParams     = DefaultMainNetParams
	DefaultMainNetParams = chaincfg.Params{
		Net:              0x22016535,
		Name:             "mainnet",
		DefaultPort:      "9443",
		PubKeyHashAddrID: 0x32, // starts with 1
		ScriptHashAddrID: 0x26, // starts with 3
		PrivateKeyID:     0x19, // starts with 5 (uncompressed) or K (compressed)

		HDPrivateKeyID: [4]byte{0x04, 0x88, 0xad, 0xe4}, // starts with xprv
		HDPublicKeyID:  [4]byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xpub

	}
	DefaultTestNetParams = chaincfg.Params{
		Net:              0x25358187,
		Name:             "testnet4",
		DefaultPort:      "19443",
		PubKeyHashAddrID: 0x6f, // starts with 1
		ScriptHashAddrID: 0xc4, // starts with 3
		PrivateKeyID:     0xef, // starts with 5 (uncompressed) or K (compressed)

		HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
		HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub
	}
	DefaultClientConfig = &client.Config{
		Host: "localhost:19442",
		User: "user",
		Pass: "pwd",
	}
)

func NewClient(host string, username string, password string) *client.Client {
	cfg := &client.Config{
		Host: host,
		User: username,
		Pass: password,
	}
	cli := client.NewClient(cfg).Startup()
	return cli
}
func NewClientWithoutOpts() *client.Client {
	cli := client.NewClient(DefaultClientConfig).Startup()
	return cli
}
func SendMany(v interface{}) (*string, error) {
	to := make(map[string]interface{})
	feeAddress := []string{}
	switch value := v.(type) {
	case *order.Statement:
		for _, payment := range value.Payments {
			if _, err := LocalNormalizePublicKey(payment.ReceiveAddress); err != nil {
				logrus.Warning(payment.ReceiveAddress, " not massgrid wallet import address")
				continue
			}
			amount := float64(payment.Amount) / 100000000
			if v, ok := to[payment.ReceiveAddress].(float64); ok {
				to[payment.ReceiveAddress] = v + amount
			} else {
				to[payment.ReceiveAddress] = amount
			}
		}
	case *order.RefundInfo:
		for _, payment := range value.Statement.Payments {
			if _, err := LocalNormalizePublicKey(payment.ReceiveAddress); err != nil {
				logrus.Warning(payment.ReceiveAddress, " not massgrid wallet import address")
				continue
			}
			amount := float64(payment.Amount) / 100000000
			if v, ok := to[payment.ReceiveAddress].(float64); ok {
				to[payment.ReceiveAddress] = v + amount
			} else {
				to[payment.ReceiveAddress] = amount
			}
		}
		for _, refund := range *value.RefundPay {
			if _, err := LocalNormalizePublicKey(refund.Drawee); err != nil {
				logrus.Warning(refund.Drawee, " not massgrid wallet import address")
				continue
			}
			amount := float64(refund.TotalAmount) / 100000000
			if v, ok := to[refund.Drawee].(float64); ok {
				to[refund.Drawee] = v + amount
			} else {
				to[refund.Drawee] = amount
			}
		}
	default:
		logrus.Error("types error")
		return nil, errors.New("types error")
	}
	for k, _ := range to {
		feeAddress = append(feeAddress, k)
	}
	cli := NewClientWithoutOpts()
	defer cli.Shutdown()
	hash, err := cli.SendManyEntire("", to, 6, false, "massgrid statement", feeAddress, false)
	if err != nil {
		logrus.Error("sendmany request error ", err)
		return nil, err
	}
	logrus.Info("sendmany successful ", *hash)
	return hash, nil
}
func LocalNormalizePublicKey(address string) (string, error) {
	address = strings.TrimSpace(address)
	btcAddress, err := btcutil.DecodeAddress(address, &DefaultNetParams)
	if err != nil {
		return "", err
	}
	if btcAddress.String() != address {
		return "", fmt.Errorf("Bitcoin NormalizeAddress mismatch %s", address)
	}
	return btcAddress.String(), nil
}
