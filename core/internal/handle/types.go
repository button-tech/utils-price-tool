package handle

type Data struct {
	Tokens     []string `json:"tokens"`
	Currencies []string `json:"currencies"`
	Change     string   `json:"change"`
	API        string   `json:"api"`
}

type UniqueData struct {
	Tokens     map[string]struct{}
	Currencies map[string]struct{}
	Change     string
	API        string
}

type response struct {
	Currency string              `json:"currency"`
	Rates    []map[string]string `json:"rates"`
}

type APIs struct {
	Name             string         `json:"name"`
	SupportedChanges []string       `json:"supported_changes"`
	SupportedFiats   map[string]int `json:"supported_fiats"`
}
