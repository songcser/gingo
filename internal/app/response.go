package app

type Response struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Level       string `json:"level"`
	Type        string `json:"type"`
}
