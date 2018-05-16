package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"strings"

	"sort"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/loader"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run() error {
	c := loader.Config{
		ParserMode: parser.ParseComments,
	}

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

	updated := map[*token.File]*ast.File{}
	for _, info := range p.InitialPackages() {
		for id, ob := range info.Uses {
			if t, _ := ob.Type().Underlying().(*types.Signature); t == nil {
				continue
			}

			name := id.Name
			if strings.HasPrefix(name, "Update") || strings.HasPrefix(name, "Upsert") {
				tokenf := p.Fset.File(id.Pos())
				filepos := token.Pos(tokenf.Base())
				for _, f := range info.Files {
					if f.Pos() <= filepos && filepos <= f.End() {
						f := f
						id.Name = "Unsafe" + name
						updated[tokenf] = f

						nodes, exact := astutil.PathEnclosingInterval(f, id.Pos(), id.Pos())
						_ = exact
						for _, n := range nodes {
							if stmt, _ := n.(ast.Stmt); stmt != nil {
								f.Comments = append(f.Comments, &ast.CommentGroup{
									List: []*ast.Comment{
										&ast.Comment{Text: "// xxx: *update*", Slash: stmt.End()},
									},
								})
								break
							}
						}
						break
					}
				}
			}
		}
	}
	for _, f := range updated {
		sort.Slice(f.Comments, func(i, j int) bool { return f.Comments[i].Pos() < f.Comments[j].Pos() })
		printer.Fprint(os.Stdout, p.Fset, f)
	}
	return nil
}
