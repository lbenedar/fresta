package json

import "github.com/lbenedar/fresta/internal/models/db"

type Item struct {
	Img       string         `json:"img"`
	Name      string         `json:"name"`
	Type      string         `json:"type"`
	Folder    string         `json:"folder"`
	Ownership map[string]int `json:"ownership"`
	Stats     Stats          `json:"_stats"`
	ID        string         `json:"_id"`
	Sort      int            `json:"sort"`
	// System           any              `json:"system,omitempty"`
	// Flags            any              `json:"flags,omitempty"`
	// Effects          []any            `json:"effects"`
}

func (i *Item) ToDB(dest **db.Item) bool {
	if dest == nil {
		return false
	}

	item := &db.Item{
		Img:    i.Img,
		Name:   i.Name,
		Type:   i.Type,
		Folder: i.Folder,
		ID:     i.ID,
		Sort:   i.Sort,
	}

	i.Stats.ToDB(&item.Stats)

	*dest = item

	return true
}
