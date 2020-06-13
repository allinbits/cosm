package typed

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
	g := genny.New()
	g.RunFn(handlerModify(opts))
	g.RunFn(aliasModify(opts))
	g.RunFn(keyModify(opts))
	g.RunFn((codecModify(opts)))
	g.RunFn((cliTxModify(opts)))
	g.RunFn((cliQueryModify(opts)))
	g.RunFn((querierModify(opts)))
	g.RunFn((keeperQuerierModify(opts)))
	g.RunFn((clientRestRestModify(opts)))
	if err := g.Box(packr.New("new/templates", "./templates")); err != nil {
		return g, err
	}
	ctx := plush.NewContext()
	ctx.Set("AppName", opts.AppName)
	ctx.Set("TypeName", opts.TypeName)
	ctx.Set("ModulePath", opts.ModulePath)
	ctx.Set("Fields", opts.Fields)
	ctx.Set("title", func(s string) string {
		return strings.Title(s)
	})
	ctx.Set("strconv", func() bool {
		strconv := false
		for _, field := range opts.Fields {
			if field.Datatype != "string" {
				strconv = true
			}
		}
		return strconv
	})
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("{{appName}}", fmt.Sprintf("%s", opts.AppName)))
	g.Transformer(genny.Replace("{{typeName}}", fmt.Sprintf("%s", opts.TypeName)))
	g.Transformer(genny.Replace("{{TypeName}}", fmt.Sprintf("%s", strings.Title(opts.TypeName))))
	return g, nil
}

func handlerModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/handler.go", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		replaceContent := fmt.Sprintf(`case MsgCreate%[1]v:
			return handleMsgCreate%[1]v(ctx, k, msg)
		default:`, strings.Title(opts.TypeName))
		content := strings.Replace(f.String(), "default:", replaceContent, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func aliasModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/alias.go", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		content := f.String() + fmt.Sprintf(`
var (
	NewMsgCreate%[1]v = types.NewMsgCreate%[1]v
)

type (
	MsgCreate%[1]v = types.MsgCreate%[1]v
)
		`, strings.Title(opts.TypeName))
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func keyModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/key.go", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		content := f.String() + fmt.Sprintf(`
const (
	%[2]vPrefix = "%[1]v-"
)
		`, opts.TypeName, strings.Title(opts.TypeName))
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func codecModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/codec.go", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		replaceContent := fmt.Sprintf(`func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreate%[1]v{}, "%[2]v/Create%[1]v", nil)
`, strings.Title(opts.TypeName), opts.AppName)
		content := strings.Replace(f.String(), "func RegisterCodec(cdc *codec.Codec) {", replaceContent, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func cliTxModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/client/cli/tx.go", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		replaceContent := fmt.Sprintf(`TxCmd.AddCommand(flags.PostCommands(
	  GetCmdCreate%[1]v(cdc),`, strings.Title(opts.TypeName), opts.AppName)
		content := strings.Replace(f.String(), "TxCmd.AddCommand(flags.PostCommands(", replaceContent, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func cliQueryModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/client/cli/query.go", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		replaceContent := fmt.Sprintf(`flags.GetCommands(
			GetCmdList%[1]v(queryRoute, cdc),`, strings.Title(opts.TypeName))
		content := strings.Replace(f.String(), "flags.GetCommands(", replaceContent, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func querierModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/types/querier.go", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		content := f.String() + fmt.Sprintf(`
const (QueryList%[2]v = "list-%[1]v")
`, opts.TypeName, strings.Title(opts.TypeName))
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func keeperQuerierModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/keeper/querier.go", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		replaceContentImport := fmt.Sprintf(`import (
	"%[1]v/x/%[2]v/types"
`, opts.ModulePath, opts.AppName)
		replaceContentDefault := fmt.Sprintf(`
		case types.QueryList%[1]v:
			return list%[1]v(ctx, k)
		default:`, strings.Title(opts.TypeName))
		content := strings.Replace(f.String(), "import (", replaceContentImport, 1)
		content = strings.Replace(content, "default:", replaceContentDefault, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}

func clientRestRestModify(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		path := fmt.Sprintf("x/%s/client/rest/rest.go", opts.AppName)
		f, err := r.Disk.Find(path)
		if err != nil {
			return err
		}
		replaceContent := fmt.Sprintf(`func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/%[1]v/%[3]v", list%[2]vHandler(cliCtx, "%[1]v")).Methods("GET")`, opts.AppName, strings.Title(opts.TypeName), opts.TypeName)
		content := strings.Replace(f.String(), "func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {", replaceContent, 1)
		newFile := genny.NewFileS(path, content)
		return r.File(newFile)
	}
}
