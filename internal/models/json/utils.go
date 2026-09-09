package json

import (
	"errors"
	"sync"

	"github.com/lbenedar/fresta/internal/models/db"
)

var (
	ErrorSetupMoreThanOne = errors.New("More than one setup data")
)

func OwnershipToDB(dest *[]db.OwnershipString, src map[string]int) {
	*dest = make([]db.OwnershipString, len(src))

	i := 0
	for k, v := range src {
		(*dest)[i].Key = k
		(*dest)[i].Value = v
		i++
	}
}

func HotbarToDB(dest *db.UserHotbars, src map[int]string) {
	*dest = make(db.UserHotbars, len(src))

	i := 0
	for k, v := range src {
		(*dest)[i].Key = k
		(*dest)[i].Value = v
		i++
	}
}

func PackageWarningsToDB(dest *db.PackageWarnings, src map[string]PackageWarningsData) {
	*dest = make([]*db.PackageWarning, len(src))

	i := 0
	for k, v := range src {
		(*dest)[i] = &db.PackageWarning{Key: k}
		v.ToDB(&(*dest)[i].Value)
		i++
	}
}

type CastableToDB[destType any] interface {
	ToDB(*destType) bool
}

func CopySliceToDB[destType any, destSlice ~[]destType, srcType CastableToDB[destType]](
	dest *destSlice,
	src []srcType,
) bool {
	if dest == nil {
		return false
	}

	*dest = make(destSlice, len(src))
	for i := range src {
		src[i].ToDB(&(*dest)[i])
	}

	return true
}

func CopySliceToDBParallel[destType any, destSlice ~[]destType, srcType CastableToDB[destType]](
	wg *sync.WaitGroup,
	dest *destSlice,
	src []srcType,
) {
	wg.Go(func() {
		CopySliceToDB(dest, src)
	})
}
