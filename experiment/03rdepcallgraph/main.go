package main

import (
	"fmt"
	"log"

	"github.com/podhmo/gomvpkg-light/build"
	"github.com/podhmo/gomvpkg-light/collect"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run() error {
	inpkg := "github.com/podhmo/safemongo/experiment"
	pkg := "github.com/podhmo/safemongo/experiment/safemongo"

	ctxt := build.OnePackageOnly()
	root, err := collect.TargetRoot(ctxt, inpkg)
	if err != nil {
		return err
	}
	log.Printf("get in-pkg %s", root.Path)

	pkgdirs, err := collect.GoFilesDirectories(ctxt, root)
	if err != nil {
		return err
	}
	log.Printf("collect candidate directories %d", len(pkgdirs))

	affected, err := collect.AffectedPackages(ctxt, pkg, root, pkgdirs)

	if err != nil {
		return err
	}
	log.Printf("collect affected packages %d", len(affected))
	for _, a := range affected {
		fmt.Println(a.Pkg, a.Files)
	}
	return nil
}
