/*

  Determine S2 cells involved in geometries.

  inspired by
        "github.com/akhenakh/oureadb/index/geodata"
	"github.com/akhenakh/oureadb/store"

  With S2 CillID's we can find which items are contained in given
  filter geometry.

*/

package main

import (
	"fmt"
	"github.com/go-spatial/geom"
	"github.com/go-spatial/geom/encoding/wkt"
	"github.com/golang/geo/s2"
	"strings"
	"sync"
)

var minLevel int
var maxLevel int
var maxCells int

var s2Lock = sync.RWMutex{}
var geoIndex s2.CellIndex

type s2Cells map[int]s2.Cell

var S2CELLS s2Cells

func init() {
	minLevel = 7
	maxLevel = 20
	maxCells = 50

	//not used for now.
	geoIndex = s2.CellIndex{}
	S2CELLS = make(s2Cells)
}

func BuildGeoIndex() {
	for i, v := range ITEMS {
		v.GeoIndex(i)
	}
	//geoIndex.Build()
}

//GeoIndex for each items determine S2Cell and store it.
func (i Item) GeoIndex(idx int) error {
	if i.GetGeometry() == "" {
		return fmt.Errorf("missing wkt geometry")
	}
	sreader := strings.NewReader(i.GetGeometry())
	g, err := wkt.Decode(sreader)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(i.GetGeometry())
		fmt.Println(i.Ekey)
		return fmt.Errorf("wkt error encountered with %s", i.Point)
	}

	p, err := geom.GetCoordinates(g)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(i.Ekey)
		fmt.Println(i.GetGeometry())
		fmt.Printf("geom error encountered with %s", i.Point)
		return fmt.Errorf("geom error")
	}
	s2Lock.Lock()
	defer s2Lock.Unlock()

	y := p[0][0]
	x := p[0][1]
	center := s2.PointFromLatLng(s2.LatLngFromDegrees(x, y))
	cell := s2.CellFromPoint(center)

	S2CELLS[idx] = cell

	return nil

}

//CalculateCover calculate S2 covering from given user polygon.
func CalculateCover(geom string) {

}

// Simple search algo
func SearchOverlapItems(items Items, cu s2.CellUnion) Items {

	newItems := make(Items, 0)

	for idx, i := range items {
		if cu.ContainsCellID(S2CELLS[idx].ID()) {
			newItems = append(newItems, i)
		}
	}

	return newItems
}

// GeoIdsAtCells returns all GeoData keys contained in the cells, without duplicates
/*
func (idx *s2.S2FlatIdx) GeoIdsAtCells(cells []s2.CellID) ([]GeoID, error) {
	m := make(map[string]struct{})

	for _, c := range cells {
		ids, err := idx.GeoIdsAtCell(c)
		if err != nil {
			return nil, errors.Wrap(err, "fetching geo ids from cells failed")
		}
		for _, id := range ids {
			m[string(id)] = struct{}{}
		}
	}

	res := make([]GeoID, len(m))
	var i int
	for k := range m {
		res[i] = []byte(k)
		i++
	}

	return res, nil
}
*/
