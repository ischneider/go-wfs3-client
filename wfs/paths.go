package wfs

import (
	"net/url"
	"path"
)

type pathStyle string

var (
	oldStylePaths = pathStyle("oldStyle")
	newStylePaths = pathStyle("newStyle")
)

// pather exists to abstract path conventions as they migrate from the
// current draft spec to the newer one
// Both styles assume a root at / and build from that
//
// The older style is:
// / -> collectionInfo
// /api -> openapi spec
// /<collection> -> get collection
// /<collection>/<fid> - get feature by id
//
// The proposed, newer style is:
// / -> api links (pointers to /api, etc.)
// /api -> openapi spec
// /collections -> collectionInfo (aggregate)
// /ceollection/<collection>/ -> single collectionInfo
// /collections/<collection>/items -> get collection items
// /collections/<collection>/items/<fid> -> get feature by id
type pather struct {
	root  *url.URL
	style pathStyle
}

func (p pather) url(paths ...string) string {
	if len(paths) == 0 {
		return p.root.String()
	}
	rel, err := url.Parse(path.Join(paths...) + "/")
	if err != nil {
		panic(err)
	}
	return p.root.ResolveReference(rel).String()
}

func (p pather) collectionInfo() string {
	switch p.style {
	case oldStylePaths:
		return p.url()
	case newStylePaths:
		return p.url("collections")
	}
	panic("path style")
}

func (p pather) spec() string {
	return p.url("api")
}

func (p pather) collectionItems(cid string) string {
	switch p.style {
	case oldStylePaths:
		return p.url(cid)
	case newStylePaths:
		return p.url("collections", cid, "items")
	}
	panic("path style")
}

func (p pather) collectionItem(cid, fid string) string {
	switch p.style {
	case oldStylePaths:
		return p.url(cid, fid)
	case newStylePaths:
		return p.url("collections", cid, "items", fid)
	}
	panic("path style")
}
