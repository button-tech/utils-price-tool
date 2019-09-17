package slip0044

import (
	"fmt"
)


func AddTrustHexToSlip() (SlipElements, error) {
	slip, err := createConvertedSlip()
	if err != nil {
		return nil, fmt.Errorf("can not createConvertedSlip: %v", err)
	}

	var address = "0x0000000000000000000000000000000000000000"
	for i:=0; i < len(slip); i++ {
		h := fmt.Sprintf("%x", i)
		hexed := address[:len(address)-len(h)] + h
		slip[i].TrustConstant = hexed
	}

	return slip, nil
}


