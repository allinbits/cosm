package app

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
)

// New ...
func New(opts *Options) (*genny.Generator, error) {
	fmt.Println(opts.Denom)
	g := genny.New()
	if err := g.Box(packr.New("app/templates", "./templates")); err != nil {
		return g, err
	}
	ctx := plush.NewContext()
	ctx.Set("ModulePath", opts.ModulePath)
	ctx.Set("AppName", opts.AppName)
	ctx.Set("Denom", opts.Denom)
	ctx.Set("title", func(s string) string {
		return strings.Title(s)
	})
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{appName}}", fmt.Sprintf("%s", opts.AppName)))
	return g, nil
}
