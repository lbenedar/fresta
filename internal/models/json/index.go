package json

import "github.com/lbenedar/fresta/internal/models/db"

type Index struct {
	Id     string `json:"_id"`
	Folder string `json:"folder"`
	Img    string `json:"img"`
	Name   string `json:"name"`
	Type   string `json:"type"`
}

func (i *Index) ToDB(dest **db.Index) bool {
	if dest == nil {
		return false
	}

	index := &db.Index{
		ID:     i.Id,
		Folder: i.Folder,
		Img:    i.Img,
		Name:   i.Name,
		Type:   i.Type,
	}

	*dest = index

	return true
}
