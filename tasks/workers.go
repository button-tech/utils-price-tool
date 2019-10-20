package tasks

import (
	"github.com/button-tech/utils-price-tool/services"
	"log"
	"sync"
)

func cmcWorker(wg *sync.WaitGroup, cont *DuiCont, list map[string]string) {
	tokens := services.CreateCMCRequestData(list)
	tokensWG := sync.WaitGroup{}

	for _, t := range tokens.Tokens {
		tokensWG.Add(1)
		go func(token services.TokensWithCurrency, tWG *sync.WaitGroup) {
			got, err := cont.Service.GetPricesCMC(token)
			if err != nil {
				log.Println(err)
				return
			}
			cont.Store.Set("cmc", got)

			tWG.Done()
		}(t, &tokensWG)
	}
	tokensWG.Wait()
	wg.Done()
}

func crcWorker(wg *sync.WaitGroup, cont *DuiCont, list map[string]string) {
	res := cont.Service.GetPricesCRC(list)
	cont.Store.Set("crc", res)
	wg.Done()
}

func huobiWorker(wg *sync.WaitGroup, cont *DuiCont, list map[string]string) {
	res := cont.Service.GetPricesHUOBI(list)
	cont.Store.Set("huobi", res)
	wg.Done()
}
