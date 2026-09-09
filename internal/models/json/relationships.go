package json

import (
	"sync"

	"github.com/lbenedar/fresta/internal/models/db"
)

type Relationships struct {
	Systems    []*RelationshipsData `json:"systems,omitempty"`
	Requires   []*RelationshipsData `json:"requires,omitempty"`
	Recommends []*RelationshipsData `json:"recommends,omitempty"`
	Conflicts  []*RelationshipsData `json:"conflicts,omitempty"`
	// RelationshipsFlags RelationshipsFlags `json:"flags"`
}

func (r *Relationships) ToDB(dest *db.Relationships) bool {
	if dest == nil {
		return false
	}

	wg := sync.WaitGroup{}
	CopySliceToDBParallel(&wg, &dest.Systems, r.Systems)
	CopySliceToDBParallel(&wg, &dest.Requires, r.Requires)
	CopySliceToDBParallel(&wg, &dest.Recommends, r.Recommends)
	CopySliceToDBParallel(&wg, &dest.Conflicts, r.Conflicts)
	wg.Wait()

	return true
}

type RelationshipsData struct {
	Id            string        `json:"id"`
	Type          string        `json:"type"`
	Manifest      string        `json:"manifest"`
	Compatibility Compatibility `json:"compatibility"`
}

func (r *RelationshipsData) ToDB(dest *db.RelationshipData) bool {
	if dest == nil {
		return false
	}

	dest.Key = r.Id
	dest.Type = r.Type
	dest.Manifest = r.Manifest

	r.Compatibility.ToDB(&dest.Compatibility)

	return true
}
