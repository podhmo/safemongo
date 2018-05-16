package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"path/filepath"

	"go/printer"
	"go/token"
	"go/types"

	"go/ast"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/loader"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run() error {
	c := loader.Config{}

	f, err := os.Open("../testdata/simple/main.go")
	if err != nil {
		return err
	}
	mainf, err := c.ParseFile(filepath.Base(f.Name()), f)
	if err != nil {
		return err
	}
	c.CreateFromFiles("main", mainf)
	// c.Import("github.com/me/app/database")
	p, err := c.Load()
	if err != nil {
		return err
	}

	for _, info := range p.InitialPackages() {
		for id, ob := range info.Uses {
			name := id.Name
			if strings.HasPrefix(name, "Update") || strings.HasPrefix(name, "Upsert") {
				tokenf := p.Fset.File(id.Pos())
				filepos := token.Pos(tokenf.Base())
				for _, f := range info.Files {
					if f.Pos() <= filepos && filepos <= f.End() {
						nodes, exact := astutil.PathEnclosingInterval(f, id.Pos(), id.Pos())
						_ = exact
						for _, n := range nodes {
							if stmt, _ := n.(ast.Stmt); stmt != nil {
								if t, _ := ob.Type().Underlying().(*types.Signature); t != nil {
									var b bytes.Buffer
									printer.Fprint(&b, p.Fset, n)
									fmt.Fprintf(os.Stdout, "%s: %s\n", tokenf.Position(stmt.Pos()), b.String())
								} else {
									// skip
									// printer.Fprint(os.Stderr, p.Fset, n)
								}
								break
							}
						}
					}
				}
			}
		}
	}
	return nil
}
