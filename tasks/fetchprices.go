package tasks

import (
	"github.com/utils-tool_prices/services"
	"github.com/utils-tool_prices/storage"
	"time"
)

func NewGetGroupTask(timeout time.Duration, s services.Service, store storage.Storage) {
	ticker := time.Tick(timeout)

	go func() {
		for _ = range ticker {
			prices := s.GetAllPricesCMC()
			store.Update(prices)
		}
	}()

	store.Get()
}
