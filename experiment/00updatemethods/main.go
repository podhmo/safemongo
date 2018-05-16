package main

import (
	"fmt"
	"go/types"
	"log"
	"strings"

	"golang.org/x/tools/go/loader"
)

/*
Update*, Upsert*
Apply
*/

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run() error {
	c := loader.Config{}
	c.Import("gopkg.in/mgo.v2")
	p, err := c.Load()
	if err != nil {
		return err
	}

	//var fns []*types.Signature
	var obs []types.Object

	for _, info := range p.InitialPackages() {
		for _, name := range info.Pkg.Scope().Names() {
			ob := info.Pkg.Scope().Lookup(name)
			if !ob.Exported() {
				continue
			}

			if t, _ := ob.Type().Underlying().(*types.Signature); t != nil {
				continue
			}
			if t, _ := ob.Type().(*types.Named); t != nil {
				for i := 0; i < t.NumMethods(); i++ {
					m := t.Method(i)
					if strings.HasPrefix(m.Name(), "Update") || strings.HasPrefix(m.Name(), "Upsert") {
						fmt.Println(ob.Name(), m.FullName())
						obs = append(obs, ob)
						continue
					}
					// t := m.Type().Underlying().Underlying().(*types.Signature) // xxx
					// for i := 0; i < t.Params().Len(); i++ {
					// 	v := t.Params().At(i)
					// 	if v.Name() == "update" {
					// 		fmt.Println(v, v.Name, m.FullName())
					// 		obs = append(obs, ob)
					// 	}
					// }
				}
				continue
			}
		}
	}
	for _, ob := range obs {
		fmt.Println(ob)
	}
	return nil
}
