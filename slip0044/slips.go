package slip0044

import (
	"fmt"
	"github.com/button-tech/bip44"
	"log"
	"strconv"
	"strings"
)

type TrustWalletSlips struct {
	TWSlipsWithCrypto map[string]string
	Contract          string
	CoinSymbol        string
}

func AddTrustHexBySlip() (map[string]string, error) {
	slip := bip44.Create()
	constants := slip.Get()

	trustWalletSlips := make(map[string]string)
	trustWalletSlips["BTC"] = "0x0000000000000000000000000000000000000000"

	for i := 1; i < len(constants); i++ {
		slipHexed, err := makeHexString(constants[i].Constant)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		trustWalletSlips[constants[i].CoinSymbol] = slipHexed
	}

	return trustWalletSlips, nil
}

func makeHexString(s string) (string, error) {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return "", fmt.Errorf("can not parseInt: %v", err)
	}

	hexed := fmt.Sprintf("%x", n)
	splitter := strings.Split(hexed, "")

	ss := make([]string, 0)

	// Start from 1, because 1st item is 8. It should be skipped
	for i := 1; i < len(splitter); i++ {
		if splitter[i] != "0" {
			ss = append(splitter[i:])
			break
		}
	}

	address := "0x0000000000000000000000000000000000000000"

	joined := strings.Join(ss, "")

	hexString := address[:len(address)-len(ss)] + joined
	return hexString, nil
}
