package typeconv

import "github.com/button-tech/utils-price-tool/pkg/storage"

func StorageFiat(f string) storage.Fiat {
	return storage.Fiat(f)
}

func StorageApi(a string) storage.Api {
	return storage.Api(a)
}

func StorageCC(c string) storage.CryptoCurrency {
	return storage.CryptoCurrency(c)
}
