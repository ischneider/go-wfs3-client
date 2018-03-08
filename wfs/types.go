package wfs

// MediaType represents supported encodings.
type MediaType struct {
	p     interface{}
	Short string
	Full  string
}

type mediaTypes []MediaType

// Lookup returns a MediaType by short or full name. If not found, a zero-value
// will be returned.
func (m mediaTypes) Lookup(v string) MediaType {
	for _, t := range m {
		if v == t.Short || v == t.Full {
			return t
		}
	}
	return MediaType{}
}

// LookupShort works as per Lookup but only checks the Short specifier.
func (m mediaTypes) LookupShort(v string) MediaType {
	for _, t := range m {
		if v == t.Short {
			return t
		}
	}
	return MediaType{}
}

// MediaTypes is the collection of supported MediaType.
var MediaTypes = mediaTypes{
	MediaType{nil, "json", "application/json"},
	MediaType{nil, "geojson", "application/geo+json"},
	MediaType{nil, "html", "text/html"},
	MediaType{nil, "xml", "application/xml"},
	MediaType{nil, "ldjson", "application/ld+json"},
}
