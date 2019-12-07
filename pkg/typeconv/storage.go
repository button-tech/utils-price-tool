package typeconv

import (
	"github.com/button-tech/utils-price-tool/pkg/storage/cache"
)

func StorageFiat(f string) cache.Fiat {
	return cache.Fiat(f)
}

func StorageApi(a string) cache.Api {
	return cache.Api(a)
}

func StorageCC(c string) cache.CryptoCurrency {
	return cache.CryptoCurrency(c)
}
