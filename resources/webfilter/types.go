package webfilter

type WebFilterRequest struct {
    Nome string `json:"nome"`
    URL  string `json:"url"`
}

type WebFilterReturn struct {
    ID uint 
    Nome string `json:"nome"`
    URL  string `json:"url"`
}