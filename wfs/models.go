package wfs

// CollectionInfo is a partial model of the WFS3 concept.
type CollectionInfo struct {
	Name        string `json:"name"`
	Title       string `json:"name"`
	Description string `json:"name"`
	Extent      BBox   `json:"extent"`
}

// BBox describes an extent in the form of lx, ly, ux, uy.
type BBox [4]float64
