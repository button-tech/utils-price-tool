package slip0044

import (
	"fmt"
	"github.com/jeyldii/bip44"
	"log"
	"strconv"
	"strings"
)

type TrustWalletSlips struct {
	Contract string
	CoinSymbol string
}

func AddTrustHexBySlip() ([]*TrustWalletSlips, error) {
	slip := bip44.Create()
	constants := slip.Get()

	trustWalletSlips := make([]*TrustWalletSlips, 0)
	for i:=0; i<len(constants); i++ {
		slipHexed, err := makeHexString(constants[i].Constant)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		trustWalletSlip := TrustWalletSlips{
			Contract:   slipHexed,
			CoinSymbol: constants[i].CoinSymbol,
		}
		trustWalletSlips = append(trustWalletSlips, &trustWalletSlip)
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

	// Because cycle started by index 1
	ss = append(ss, "0x0000000000000000000000000000000000000000")

	for i:=1; i<len(splitter); i++ {
		toInt, err := strconv.Atoi(splitter[i])
		if err != nil {
			return "", fmt.Errorf("can not parseString: %v", err)
		}

		if toInt != 0 {
			ss = append(splitter[i:])
			break
		}
	}

	var address = "0x0000000000000000000000000000000000000000"

	joined := strings.Join(ss, "")

	hexString := address[:len(address)-len(ss)] + joined
	return hexString, nil
}


