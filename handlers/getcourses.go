package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"log"
	"net/http"
	"os"
	"sync"
)

// data what to get
type dataTokensAndCurrencies struct {
	Tokens      []string `json:"tokens"`
	Currencies []string `json:"currencies"`
	Change     string   `json:"change"`
	API        string   `json:"api"`
}

// data to trustwallet
type dataGetPrices struct {
	Currency string  `json:"currency"`
	Tokens    []Token `json:"tokens"`
}

type Token struct {
	Contract string `json:"contract"`
}

// from trustwallet
type gotPrices struct {
	Status bool `json:"status"`
	Docs   []DocsPrices`json:"docs"`
	Currency string `json:"currency"`
}

type DocsPrices struct {
	Price            string `json:"price"`
	Contract         string `json:"contract"`
	PercentChange24H string `json:"percent_change_24h"`
}

type resultPrices struct {
	mu   sync.Mutex
	Data []Prices `json:"data"`
}

type Prices struct {
	Currency string              `json:"currency"`
	Rates    []map[string]string `json:"rates"`
}

func getCourses(c *gin.Context) {
	resp := dataTokensAndCurrencies{}
	if err := c.ShouldBindJSON(&resp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}

	dGetPrice := dataGetPrices{}
	dGetPrice.Tokens = collect(resp.Tokens, collectHelper)



}

//	collect([]string{"1", "2"}, collectHelper)
func collect(list []string, f func(string) Token) []Token {
	result := make([]Token, len(list))
	for i, item := range list {
		result[i] = f(item)
	}

	return result
}

func collectHelper(token string) Token {
	return Token{
		Contract: token,
	}
}

func prepareAnswer(g gotPrices, ch <- chan gotPrices) {
	for range ch {
		gPrices := <- ch
		result := Prices{}
		result.Currency = gPrices.Currency
		for k, v := range gPrices.Docs {

		}
	}

}

func (data dataGetPrices) getPrices(currencies []string, ch chan gotPrices)   {
	for _, c := range currencies {
		data.Currency = c

		gPrices := gotPrices{}
		go func(){
			gPrices.getPrices(data)
			ch <- gPrices
		}()
	}
}

func (r *gotPrices) getPrices(d dataGetPrices) {
	url := os.Getenv("TRUST_URL")
	rq, err := req.Post(url, req.BodyJSON(d))
	if err != nil {
		log.Println(err)
		return
	}

	if err = rq.ToJSON(r); err != nil {
		log.Println(err)
		return
	}
}

func InitRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api")
	{
		v1.GET("/prices", getCourses)
	}

	return r
}
