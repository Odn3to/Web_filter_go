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

type TokenBody struct {
	Token string `json:"token"`
}

// swagger: Response
type Response struct {
	Message string
	Data string
}

// swagger: ResponseSquid
type ResponseSquid struct {
	Class string
	Text string
}