package main

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strings"

	"gabe565.com/trmnl-nightscout/internal/config"
)

type Render struct {
	// Access token. Required if the `ACCESS_TOKEN` env is set.
	Token string `env:"TOKEN"`

	config.Render
}

func main() {
	var output bytes.Buffer
	output.WriteString("# Query Parameters\n\n")
	output.WriteString(
		"The following GET parameters can be passed to the JSON or image endpoint. These values will override the env config defined in [`envs.md`](envs.md).\n\n",
	)

	if err := renderFile("internal/generate/docs/main.go", &output); err != nil {
		panic(err)
	}
	if err := renderFile("internal/config/render.go", &output); err != nil {
		panic(err)
	}

	if err := os.MkdirAll("docs", 0o755); err != nil {
		panic(err)
	}
	if err := os.WriteFile("docs/query-params.md", output.Bytes(), 0o644); err != nil {
		panic(err)
	}
}

func renderFile(path string, buf *bytes.Buffer) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}

		for _, spec := range gen.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok || ts.Name == nil || ts.Name.Name != "Render" {
				continue
			}

			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}

			for _, field := range st.Fields.List {
				if field.Tag == nil || len(field.Names) == 0 {
					continue
				}

				tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
				env, _, _ := strings.Cut(tag.Get("env"), ",")
				if env == "" || env == "-" {
					continue
				}
				envDefault := tag.Get("envDefault")
				query := strings.ToLower(env)

				comment := strings.TrimSpace(field.Doc.Text())
				if comment == "" {
					continue
				}

				buf.WriteString("- `")
				buf.WriteString(query)
				buf.WriteString("` ")
				if envDefault != "" {
					buf.WriteString("(default `")
					buf.WriteString(envDefault)
					buf.WriteString("`) ")
				}
				buf.WriteString("- ")
				buf.WriteString(comment)
				buf.WriteByte('\n')
			}
		}
	}

	return nil
}
