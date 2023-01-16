package db

import (
	"os"
	"path"
	"runtime"
	"strings"
	"text/template"

	"github.com/wildberries-ru/go-transport-generator/pkg/api"
)

const cacheTpl = `// Package {{.PkgName}} ...
package {{.PkgName}}

import (
	"context"
	"strconv"
	"time"

)

// middleware wraps Service
type middleware struct {
	svc         {{ .Iface.Name }}
}

{{$methods := .HTTPMethods}}
{{range .Iface.Methods -}}
{{$method := index $methods .Name}}
// {{.Name}} ...
func (m *middleware) {{.Name}}({{joinFullVariables .Args ","}}) ({{joinFullVariables .Results ","}}) {
	return m.svc.{{.Name}}({{joinVariableNamesWithEllipsis .Args ","}})
}
{{end}}

// NewMiddleware ...
func NewMiddleware(
	svc {{ .Iface.Name }},
) {{ .Iface.Name }} {
	return &middleware{
		svc:         svc,
	}
}
`

// Base ...
type Base struct {
	*template.Template
	filePath []string
	imports  imports
}

// Generate ...
func (s *Base) Generate(info api.Interface) (err error) {
	if runtime.GOOS == "windows" {
		info.AbsOutputPath = strings.Replace(info.AbsOutputPath, `\`, "/", -1)
	}
	info.PkgName = path.Base(info.AbsOutputPath)
	info.AbsOutputPath = strings.Join(append(strings.Split(info.AbsOutputPath, "/"), s.filePath...), "/")
	dir, _ := path.Split(info.AbsOutputPath)
	err = os.MkdirAll(dir, 0750)
	if err != nil {
		return
	}
	serverFile, err := os.Create(info.AbsOutputPath)
	defer func() {
		_ = serverFile.Close()
	}()
	t := template.Must(s.Parse(cacheTpl))
	if err = t.Execute(serverFile, info); err != nil {
		return
	}
	err = s.imports.GoImports(info.AbsOutputPath)
	return
}

// NewBase ...
func NewBase(template *template.Template, filePath []string, imports imports) *Base {
	return &Base{
		Template: template,
		filePath: filePath,
		imports:  imports,
	}
}
