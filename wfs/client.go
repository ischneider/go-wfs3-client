package wfs

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/jban332/kin-openapi/openapi3"
)

// Service represents a single WFS3 service.
type Service struct {
	cl    Client
	spec  *openapi3.Swagger
	paths pather
}

// ServiceInfo is a high-level summary of the service.
type ServiceInfo struct {
	URL         string
	Description string
}

// Info returns the ServiceInfo.
func (s Service) Info() ServiceInfo {
	url := "undefined"
	if len(s.spec.Servers) > 0 {
		url = s.spec.Servers[0].URL
	}
	return ServiceInfo{
		url,
		s.spec.Info.Description,
	}
}

// DescribeCollections returns the Operation that retrieves the total set of
// feature collection metadata.
func (s Service) DescribeCollections() (Operation, error) {
	for _, op := range s.Operations() {
		if op.ID == "describeCollections" {
			return op, nil
		}
	}
	return Operation{}, fmt.Errorf("No describeCollections operation defined")
}

// Operations returns all supported Operations.
func (s Service) Operations() []Operation {
	paths := s.spec.Paths
	ops := []Operation{}
	for k, path := range paths {
		for _, op := range path.Operations() {
			ops = append(ops, s.operationFromSwagger(k, op))
		}
	}
	return ops
}

func (s Service) operationFromSwagger(path string, op *openapi3.Operation) Operation {
	params := []Parameter{}
	for _, p := range op.Parameters {
		pv := p.Value
		params = append(params, Parameter{
			p:           pv,
			Description: pv.Description,
			Name:        pv.Name,
			Required:    pv.Required,
			Type:        pv.Schema.Value.Type,
		})
	}
	return Operation{
		svc:         s,
		Description: op.Description,
		ID:          op.OperationID,
		Path:        path,
		Params:      params,
	}
}

// GetOperation returns an Operation by ID. An error is returned
// if not found.
func (s Service) GetOperation(id string) (Operation, error) {
	paths := s.spec.Paths
	for k, path := range paths {
		for _, op := range path.Operations() {
			if op.OperationID == id {
				return s.operationFromSwagger(k, op), nil
			}
		}
	}
	return Operation{}, fmt.Errorf("no operation %q", id)
}

// Operation represents an operation defined by the Service
type Operation struct {
	svc         Service
	Description string
	ID          string
	Path        string
	Params      []Parameter
}

func findParameter(params []Parameter, name string) (Parameter, bool) {
	var f Parameter
	for _, p := range params {
		if p.Name == name {
			return p, true
		}
	}
	return f, false
}

// URL returns the full URL of the Operation.
func (o Operation) URL() string {
	return o.svc.Info().URL + o.Path
}

// SimpleCall returns a Call that has no parameters.
func (o Operation) SimpleCall() Call {
	return Call{o, nil, "json"}
}

// Call returns a Call with parameters as provided by the given map. An
// error is returned if any parameter key is not found.
// The default response format is configured as "json".
func (o Operation) Call(params map[string]interface{}) (Call, error) {
	pv := []parameterValue{}
	for k, v := range params {
		p, ok := findParameter(o.Params, k)
		if !ok {
			return Call{}, fmt.Errorf("no parameter named %q", k)
		}
		pv = append(pv, parameterValue{p, v})
	}
	return Call{o, pv, "json"}, nil
}

// Call represents a pending invocation of an Operation.
type Call struct {
	op        Operation
	params    []parameterValue
	mediaType string
}

// Accept returns a Call that will use the provided media type.
func (c Call) Accept(mediaType string) Call {
	c.mediaType = mediaType
	return c
}

// ExecuteWriter will invoke the Call operation writing the response to the
// provided io.Writer.
func (c Call) ExecuteWriter(w io.Writer) error {
	req, err := c.buildRequest()
	if err != nil {
		return err
	}
	mt := MediaTypes.Lookup(c.mediaType)
	if mt.Full == "" {
		return fmt.Errorf("No Media Type: %s", c.mediaType)
	}
	if c.mediaType != "" {
		req.Header.Set("Accept", mt.Full)
	}
	return c.op.svc.cl.doWriter(req, w)
}

func (c Call) buildRequest() (*http.Request, error) {
	url := c.op.URL()
	for _, pv := range c.params {
		if pv.Def.p.In == "path" {
			url = strings.Replace(url, fmt.Sprintf("{%s}", pv.Def.p.Name), fmt.Sprint(pv.Value), -1)
		} else {
			return nil, fmt.Errorf("paramter in %s not supported for %s", pv.Def.p.In, pv.Def.Name)
		}
	}
	return http.NewRequest("GET", url, nil)
}

// Parameter represents an optional or mandatory argument to an Operation.
type Parameter struct {
	p           *openapi3.Parameter
	Description string
	Name        string
	Required    bool
	Type        string
}

type parameterValue struct {
	Def   Parameter
	Value interface{}
}

// Client provides a WFS3 client.
type Client struct {
	client *http.Client
}

// NewClient creates a Client that will use the provided http.Client.
func NewClient(cl *http.Client) Client {
	return Client{cl}
}

func (c Client) do(r *http.Request) ([]byte, error) {
	// @todo config
	r.Header.Set("Cache-Control", "max-age=300")
	resp, err := c.client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error calling %s : %s", r.URL, err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http error calling %s : %d - %s", r.URL, resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (c Client) doWriter(r *http.Request, w io.Writer) error {
	resp, err := c.client.Do(r)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("http error calling %s : %d - %s", r.URL, resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()
	_, err = io.Copy(w, resp.Body)
	return err
}

// Connect will request the spec from the provided service as defined by the
// urlRoot. oldStyle exists as a temporary toggle to switch between path
// conventions as the spec evolves.
func (c Client) Connect(urlRoot string, oldStyle bool) (Service, error) {
	st := oldStylePaths
	if !oldStyle {
		st = newStylePaths
	}
	if !strings.HasSuffix(urlRoot, "/") {
		urlRoot = urlRoot + "/"
	}
	u, err := url.Parse(urlRoot)
	if err != nil {
		return Service{}, err
	}
	paths := pather{u, st}
	specURL := paths.spec()
	req, err := http.NewRequest("GET", specURL, nil)
	if err != nil {
		return Service{}, fmt.Errorf("invalid spec path: %s", specURL)
	}
	req.Header.Set("Accept", MediaTypes.LookupShort("json").Full)
	bytes, err := c.do(req)
	if err != nil {
		return Service{}, err
	}
	spec, err := parseSpec(bytes)
	if err != nil {
		return Service{}, err
	}
	return Service{c, spec, paths}, nil
}
