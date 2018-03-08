package wfs

import (
	"net/url"
	"testing"
)

func TestPather(t *testing.T) {
	u, err := url.Parse("http://server.domain/path/")
	if err != nil {
		panic(err)
	}
	old := pather{u, oldStylePaths}
	if p := old.spec(); p != "http://server.domain/path/api/" {
		t.Errorf("old spec %s", p)
	}
	if p := old.collectionInfo(); p != "http://server.domain/path/" {
		t.Errorf("old collection info %s", p)
	}
	if p := old.collectionItems("c"); p != "http://server.domain/path/c/" {
		t.Errorf("old collection items %s", p)
	}
	if p := old.collectionItem("c", "f"); p != "http://server.domain/path/c/f/" {
		t.Errorf("old collection item %s", p)
	}

}
