package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/ischneider/go-wfs-client/wfs"
	flags "github.com/jessevdk/go-flags"
)

var opts = &struct {
	Encoding string `short:"e" long:"encoding" description:"specify the encoding" default:"application/json"`
	Verbose  bool   `short:"v" long:"verbose" description:"be noisier"`
}{}

func createClient() wfs.Client {
	cdir := os.Getenv("HTTP_CACHE_DIR")
	if cdir == "" {
		cdir = filepath.Join(os.TempDir(), "wfs-http-cache")
	}
	return wfs.NewClient(&http.Client{Transport: httpcache.NewTransport(diskcache.New(cdir))})
}

func connect(svc string) (wfs.Service, error) {
	cl := createClient()
	fmt.Println("connecting to", svc)
	// @todo can we sniff out new/old style or make explicit via flag
	return cl.Connect(svc, true)
}

type Info struct {
	Args struct {
		Source string
	} `positional-args:"y"`
}

func (r *Info) Execute([]string) error {
	svc, err := connect(r.Args.Source)
	if err != nil {
		return err
	}
	info := svc.Info()
	fmt.Println("Service Info:")
	fmt.Println("\tURL: ", info.URL)
	fmt.Println("\tDescription: ", info.Description)
	fmt.Println()
	ops := svc.Operations()
	fmt.Println("Operations:")
	for _, op := range ops {
		fmt.Println("\tOperation: ", op.ID, "[", op.URL(), "]")
		if opts.Verbose {
			fmt.Println("\t", op.Description)
		}
		fmt.Println("\tParameters:")
		for _, p := range op.Params {
			required := ""
			if p.Required {
				required = "*Required*"
			}
			fmt.Print("\t\t", p.Name, "["+p.Type+"]", required)
			if opts.Verbose {
				fmt.Print(p.Description)
			}
			fmt.Println()
			fmt.Println()
		}
	}
	return nil
}

type Collections struct {
	Args struct {
		Source string
	} `positional-args:"y"`
}

func (r Collections) Execute([]string) error {
	svc, err := connect(r.Args.Source)
	if err != nil {
		return err
	}
	op, err := svc.DescribeCollections()
	if err != nil {
		return err
	}
	return op.SimpleCall().Accept(opts.Encoding).ExecuteWriter(os.Stdout)
}

type Operation struct {
	Args struct {
		Source    string
		Operation string
	} `positional-args:"y"`
}

func (o Operation) Execute(args []string) error {
	svc, err := connect(o.Args.Source)
	if err != nil {
		return err
	}
	op, err := svc.GetOperation(o.Args.Operation)
	if err != nil {
		return err
	}
	params := map[string]interface{}{}
	for _, a := range args {
		parts := strings.Split(a, "=")
		if len(parts) != 2 {
			return fmt.Errorf("arg inputs require key=value, got %q", a)
		}
		params[parts[0]] = parts[1]
	}
	call, err := op.Call(params)
	if err != nil {
		return err
	}
	return call.Accept(opts.Encoding).ExecuteWriter(os.Stdout)
}

func buildParser() *flags.Parser {
	parser := flags.NewParser(opts, flags.Default)
	for _, c := range []struct {
		cmd               interface{}
		name, short, long string
	}{
		{&Info{}, "info", "Service Info", ""},
		{&Collections{}, "coll", "Collection Info", ""},
		{&Operation{}, "op", "Execute Operation", "Arguments in form of name=value"},
	} {
		_, e := parser.AddCommand(c.name, c.short, c.long, c.cmd)
		if e != nil {
			panic(e)
		}
	}
	return parser
}

func main() {
	parser := buildParser()
	_, err := parser.Parse()
	if flagerr, ok := err.(*flags.Error); ok {
		// error message will be printed by go-flags
		// if not the help error, write help out to be nice
		if flagerr.Type != flags.ErrHelp {
			fmt.Println()
			parser.WriteHelp(os.Stdout)
		}
		os.Exit(1)
	} else if err != nil {
		// error is already printed by go-flags
		os.Exit(1)
	}
}
