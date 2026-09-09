package json

import "github.com/lbenedar/fresta/internal/models/db"

type Journal struct {
	Folder    any            `json:"folder"`
	Name      string         `json:"name"`
	Pages     []*JournalPage `json:"pages"`
	Sort      int            `json:"sort"`
	Ownership map[string]int `json:"ownership,omitempty"`
	ID        string         `json:"_id"`
	// JournalFlags      any               `json:"flags,omitempty"`
	// JournalStats      Stats             `json:"_stats"`
	// Categories []any          `json:"categories"`
}

func (j *Journal) ToDB(dest **db.Journal) bool {
	if dest == nil {
		return false
	}

	journal := &db.Journal{
		Name: j.Name,
		Sort: j.Sort,
		ID:   j.ID,
	}

	OwnershipToDB(&journal.Ownership, j.Ownership)
	CopySliceToDB(&journal.Pages, j.Pages)

	*dest = journal

	return true
}
