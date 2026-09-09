package json

import "github.com/lbenedar/fresta/internal/models/db"

type DocumentTypes struct {
	Actor          DocumentTypeData `json:"Actor"`
	Item           DocumentTypeData `json:"Item"`
	RegionBehavior RegionBehavior   `json:"RegionBehavior,omitempty"`
}

func (d *DocumentTypes) ToDB(dest *db.DocumentTypes) bool {
	if dest == nil {
		return false
	}

	dest.Data = make([]*db.DocumentTypeData, 2)

	dest.Data[0] = &db.DocumentTypeData{Type: "Actor"}
	d.Actor.ToDB(dest.Data[0])

	dest.Data[1] = &db.DocumentTypeData{Type: "Item"}
	d.Item.ToDB(dest.Data[1])

	return true
}

type DocumentTypeData struct {
	HtmlFields []string `json:"htmlFields"`
}

func (d *DocumentTypeData) ToDB(dest *db.DocumentTypeData) bool {
	if dest == nil {
		return false
	}

	dest.HtmlFields = d.HtmlFields

	return true
}
