package wfs

import (
	"fmt"

	"github.com/jban332/kin-openapi/openapi3"
)

// parseSpec attempts to parse and validate the provided specification.
// if the components are missing, it patches some builtins in and then
// attempts to validate
func parseSpec(data []byte) (*openapi3.Swagger, error) {
	swag := &openapi3.Swagger{}
	if err := swag.UnmarshalJSON(data); err != nil {
		panic(err)
	}
	loader := openapi3.NewSwaggerLoader()
	if err := loader.ResolveRefsIn(swag); err != nil {
		comps := openapi3.NewComponents()
		if err := comps.UnmarshalJSON([]byte(_componentsJSON)); err != nil {
			panic(err)
		}
		fmt.Println("WARNING: patching in std components", comps.Parameters)
		swag.Components = comps
	}
	if err := loader.ResolveRefsIn(swag); err != nil {
		return nil, err
	}
	return swag, nil
}

const _componentsJSON = `
{
  "schemas" : {
    "exception" : {
      "required" : [ "code" ],
      "type" : "object",
      "properties" : {
        "code" : {
          "type" : "string",
          "xml" : {
            "name" : "code",
            "attribute" : true
          }
        },
        "description" : {
          "type" : "string",
          "xml" : {
            "name" : "description",
            "namespace" : "http://www.opengis.net/wfs/3.0",
            "prefix" : "wfs"
          }
        }
      },
      "xml" : {
        "name" : "Exception",
        "namespace" : "http://www.opengis.net/wfs/3.0",
        "prefix" : "wfs"
      }
    },
    "bbox" : {
      "required" : [ "bbox" ],
      "type" : "object",
      "properties" : {
        "crs" : {
          "type" : "string",
          "enum" : [ "http://www.opengis.net/def/crs/OGC/1.3/CRS84" ],
          "default" : "http://www.opengis.net/def/crs/OGC/1.3/CRS84"
        },
        "bbox" : {
          "type" : "array",
          "description" : "minimum longitude, minimum latitude, maximum longitude, maximum latitude",
          "example" : [ -180, -90, 180, 90 ],
          "items" : {
            "maximum" : 180,
            "minimum" : -180,
            "maxItems" : 4,
            "minItems" : 4,
            "type" : "number"
          }
        }
      }
    },
    "content" : {
      "required" : [ "collections" ],
      "type" : "object",
      "properties" : {
        "collections" : {
          "type" : "array",
          "xml" : {
            "namespace" : "http://www.opengis.net/wfs/3.0",
            "prefix" : "wfs"
          },
          "items" : {
            "$ref" : "#/components/schemas/collectionInfo"
          }
        }
      },
      "xml" : {
        "name" : "Collections",
        "namespace" : "http://www.opengis.net/wfs/3.0",
        "prefix" : "wfs"
      }
    },
    "collectionInfo" : {
      "required" : [ "links", "name" ],
      "type" : "object",
      "properties" : {
        "name" : {
          "type" : "string",
          "example" : "address"
        },
        "title" : {
          "type" : "string",
          "example" : "address"
        },
        "description" : {
          "type" : "string",
          "example" : "An address."
        },
        "links" : {
          "type" : "array",
          "items" : {
            "$ref" : "#/components/schemas/link"
          }
        },
        "extent" : {
          "$ref" : "#/components/schemas/bbox"
        },
        "crs" : {
          "type" : "array",
          "items" : {
            "type" : "string"
          }
        }
      }
    },
    "link" : {
      "required" : [ "href" ],
      "type" : "object",
      "properties" : {
        "href" : {
          "type" : "string",
          "example" : "http://data.example.com/buildings/123",
          "xml" : {
            "attribute" : true
          }
        },
        "rel" : {
          "type" : "string",
          "example" : "prev",
          "xml" : {
            "attribute" : true
          }
        },
        "type" : {
          "type" : "string",
          "example" : "application/gml+xml;version=3.2",
          "xml" : {
            "attribute" : true
          }
        },
        "hreflang" : {
          "type" : "string",
          "example" : "en",
          "xml" : {
            "attribute" : true
          }
        }
      }
    },
    "featureCollectionGeoJSON" : {
      "required" : [ "features" ],
      "type" : "object",
      "properties" : {
        "features" : {
          "type" : "array",
          "items" : {
            "$ref" : "#/components/schemas/featureGeoJSON"
          }
        }
      }
    },
    "featureGeoJSON" : {
      "required" : [ "geometry", "properties", "type" ],
      "type" : "object",
      "properties" : {
        "type" : {
          "type" : "string",
          "enum" : [ "Feature" ]
        },
        "geometry" : {
          "$ref" : "#/components/schemas/geometryGeoJSON"
        },
        "properties" : {
          "type" : "object",
          "nullable" : true
        },
        "id" : { }
      }
    },
    "geometryGeoJSON" : {
      "required" : [ "type" ],
      "type" : "object",
      "properties" : {
        "type" : {
          "type" : "string",
          "enum" : [ "Point", "MultiPoint", "LineString", "MultiLineString", "Polygon", "MultiPolygon", "GeometryCollection" ]
        }
      }
    },
    "featureCollectionGML" : {
      "type" : "object",
      "properties" : {
        "features" : {
          "type" : "array",
          "items" : {
            "xml" : {
              "name" : "featureMember",
              "namespace" : "http://www.opengis.net/wfs/3.0",
              "prefix" : "wfs"
            },
            "oneOf" : [ {
              "$ref" : "#/components/schemas/referenceXlink"
            }, {
              "$ref" : "#/components/schemas/featureGML"
            } ]
          }
        }
      },
      "xml" : {
        "name" : "FeatureCollection",
        "namespace" : "http://www.opengis.net/wfs/3.0",
        "prefix" : "wfs"
      }
    },
    "featureGML" : {
      "required" : [ "gml:id" ],
      "type" : "object",
      "properties" : {
        "gml:id" : {
          "type" : "string",
          "xml" : {
            "namespace" : "http://www.opengis.net/gml/3.2",
            "prefix" : "gml",
            "attribute" : true
          }
        }
      },
      "xml" : {
        "name" : "AbstractFeature",
        "namespace" : "http://www.opengis.net/gml/3.2",
        "prefix" : "gml"
      }
    },
    "referenceXlink" : {
      "required" : [ "href" ],
      "type" : "object",
      "properties" : {
        "href" : {
          "type" : "string",
          "xml" : {
            "namespace" : "http://www.w3.org/1999/xlink",
            "prefix" : "xlink",
            "attribute" : true
          }
        },
        "title" : {
          "type" : "string",
          "xml" : {
            "namespace" : "http://www.w3.org/1999/xlink",
            "prefix" : "xlink",
            "attribute" : true
          }
        }
      }
    }
  },
  "parameters" : {
    "f" : {
      "name" : "f",
      "in" : "query",
      "description" : "The format of the response. If no value is provided, the standard http rules apply, i.e., the accept header shall be used to determine the format.\\\nPre-defined values are \"xml\", \"json\" and \"html\". The response to other  values is determined by the server.",
      "required" : false,
      "style" : "form",
      "explode" : false,
      "schema" : {
        "type" : "string",
        "enum" : [ "json", "xml", "html" ]
      },
      "example" : "json"
    },
    "startIndex" : {
      "name" : "startIndex",
      "in" : "query",
      "description" : "The optional startIndex parameter indicates the index within the result set from which the server shall begin presenting results in the response document. The first element has an index of 0.\\\nMinimum = 0.\\\nDefault = 0.",
      "required" : false,
      "style" : "form",
      "explode" : false,
      "schema" : {
        "minimum" : 0,
        "type" : "integer",
        "default" : 0
      },
      "example" : 0
    },
    "count" : {
      "name" : "count",
      "in" : "query",
      "description" : "The optional count parameter limits the number of items that are presented in the response document.\\\nOnly items are counted that are on the first level of the collection in the response document.  Nested objects contained within the explicitly requested items shall not be counted.\\\nMinimum = 1.\\\nMaximum = 10000.\\\nDefault = 10.",
      "required" : false,
      "style" : "form",
      "explode" : false,
      "schema" : {
        "maximum" : 10000,
        "minimum" : 1,
        "type" : "integer",
        "default" : 10
      },
      "example" : 10
    },
    "bbox" : {
      "name" : "bbox",
      "in" : "query",
      "description" : "Only features that have a geometry that intersects the bounding box are selected. The bounding box is provided as four numbers:\n \n* Lower corner, coordinate axis 1 (minimum longitude) * Lower corner, coordinate axis 2 (minimum latitude) * Upper corner, coordinate axis 1 (maximum longitude) * Upper corner, coordinate axis 2 (maximum latitude)",
      "required" : false,
      "style" : "form",
      "explode" : false,
      "schema" : {
        "type" : "array",
        "items" : {
          "maximum" : 180,
          "minimum" : -180,
          "maxItems" : 4,
          "minItems" : 4,
          "type" : "number"
        }
      }
    },
    "resultType" : {
      "name" : "resultType",
      "in" : "query",
      "description" : "This service will respond to a query in one of two ways (excluding an exception response). It may either generate a complete response document containing resources that satisfy the operation or it may simply generate an empty response container that indicates the count of the total number of resources that the operation would return. Which of these two responses is generated is determined by the value of the optional resultType parameter.\\\nThe allowed values for this parameter are \"results\" and \"hits\".\\\nIf the value of the resultType parameter is set to \"results\", the server will generate a complete response document containing resources that satisfy the operation.\\\nIf the value of the resultType attribute is set to \"hits\", the server will generate an empty response document containing no resource instances.\\\nDefault = \"results\".",
      "required" : false,
      "style" : "form",
      "explode" : false,
      "schema" : {
        "type" : "string",
        "enum" : [ "hits", "results" ],
        "default" : "results"
      },
      "example" : "results"
    },
    "id" : {
      "name" : "id",
      "in" : "path",
      "description" : "The id of a feature",
      "required" : true,
      "schema" : {
        "type" : "string"
      }
    }
  }
}
`
