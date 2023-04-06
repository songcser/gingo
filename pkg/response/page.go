package response

import "github.com/songcser/gingo/pkg/model"

type Page struct {
	Total   int64 `json:"total"`
	Size    int   `json:"size"`
	Current int   `json:"current"`
	Results any   `json:"results"`
}

func NewPage(page model.Page) Page {
	return Page{
		Size:    page.GetSize(),
		Total:   page.GetTotal(),
		Current: page.GetCurrent(),
	}
}
