package json

import "github.com/lbenedar/fresta/internal/models/db"

type Table struct {
	Name        string         `json:"name"`
	Results     []*TableResult `json:"results"`
	Description string         `json:"description"`
	Formula     string         `json:"formula"`
	ID          string         `json:"_id"`
	Img         string         `json:"img"`
	Replacement bool           `json:"replacement"`
	DisplayRoll bool           `json:"displayRoll"`
	Folder      string         `json:"folder"`
	Sort        int            `json:"sort"`
	Ownership   map[string]int `json:"ownership,omitempty"`
	Stats       Stats          `json:"_stats"`
	// TablesFlags      any              `json:"flags,omitempty"`
}

func (t *Table) ToDB(dest **db.Table) bool {
	if dest == nil {
		return false
	}

	table := &db.Table{
		Name:        t.Name,
		Description: t.Description,
		Formula:     t.Formula,
		ID:          t.ID,
		Img:         t.Img,
		Replacement: t.Replacement,
		DisplayRoll: t.DisplayRoll,
		Folder:      t.Folder,
		Sort:        t.Sort,
	}

	t.Stats.ToDB(&table.Stats)

	OwnershipToDB(&table.Ownership, t.Ownership)

	CopySliceToDB(&table.Results, t.Results)

	*dest = table

	return true
}

type TableResult struct {
	Type        string `json:"type"`
	Weight      int    `json:"weight"`
	Range       []int  `json:"range"`
	Drawn       bool   `json:"drawn"`
	ID          string `json:"_id"`
	Img         string `json:"img"`
	Stats       Stats  `json:"_stats"`
	Description string `json:"description"`
	Name        string `json:"name"`
	// ResultsFlags any `json:"flags"`
}

func (t *TableResult) ToDB(dest **db.TableResult) bool {
	if dest == nil {
		return false
	}

	tableResult := &db.TableResult{
		Type:        t.Type,
		Weight:      t.Weight,
		Drawn:       t.Drawn,
		ID:          t.ID,
		Img:         t.Img,
		Description: t.Description,
		Name:        t.Name,
	}

	tableResult.Range = make([]int, len(t.Range))
	copy(tableResult.Range, t.Range)

	t.Stats.ToDB(&tableResult.Stats)

	*dest = tableResult

	return true
}
