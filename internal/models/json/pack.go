package json

import (
	"sync"

	"github.com/lbenedar/fresta/internal/models/db"
)

type Pack struct {
	Name        string        `json:"name"`
	Label       string        `json:"label"`
	Banner      string        `json:"banner"`
	Path        string        `json:"path"`
	Type        string        `json:"type"`
	System      string        `json:"system"`
	Ownership   Ownership     `json:"ownership"`
	PackageType string        `json:"packageType,omitempty"`
	PackageName string        `json:"packageName,omitempty"`
	Id          string        `json:"id,omitempty"`
	Index       []*Index      `json:"index"`
	Folders     []*PackFolder `json:"folders"`
	// SystemPacksFlags SystemPacksFlags `json:"flags"`
}

func (p *Pack) ToDB(dest **db.Pack) bool {
	if dest == nil {
		return false
	}

	pack := &db.Pack{
		Name:        p.Name,
		Label:       p.Label,
		Banner:      p.Banner,
		Path:        p.Path,
		Type:        p.Type,
		System:      p.System,
		PackageType: p.PackageType,
		PackageName: p.PackageName,
		ID:          p.Id,
	}

	p.Ownership.ToDB(&pack.Ownership)

	wg := sync.WaitGroup{}
	CopySliceToDBParallel(&wg, &pack.Index, p.Index)
	CopySliceToDBParallel(&wg, &pack.Folders, p.Folders)
	wg.Wait()

	*dest = pack

	return true
}

type PackFolder struct {
	ID          string `json:"_id"`
	Color       any    `json:"color"`
	Description string `json:"description"`
	Folder      any    `json:"folder"`
	Name        string `json:"name"`
	Sort        int    `json:"sort"`
	Sorting     string `json:"sorting"`
	Type        string `json:"type"`
	// Packs0FoldersFlags any `json:"flags"`
}

func (p *PackFolder) ToDB(dest **db.PackFolder) bool {
	if dest == nil {
		return false
	}

	packFolder := &db.PackFolder{
		ID:          p.ID,
		Description: p.Description,
		Name:        p.Name,
		Sort:        p.Sort,
		Sorting:     p.Sorting,
		Type:        p.Type,
	}

	*dest = packFolder

	return true
}
