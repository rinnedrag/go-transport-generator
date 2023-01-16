package db

import (
	"os"
	"path"
	"runtime"
	"strings"
	"text/template"

	"github.com/rinnedrag/go-transport-generator/pkg/api"
)

const workerMetricsTpl = `// Package {{.PkgName}} ...
// CODE GENERATED AUTOMATICALLY
// DO NOT EDIT
package {{.PkgName}}

import (
	"context"
	"strconv"
	"time"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/go-kit/kit/metrics"
)

// instrumentingMiddleware wraps Service and enables request metrics
type instrumentingMiddleware struct {
	reqCount       metrics.Counter
	processedBatch metrics.Counter
	reqDuration    metrics.Histogram
	svc            {{ .Iface.Name }}
}

{{$methods := .HTTPMethods}}
{{range .Iface.Methods -}}
{{$method := index $methods .Name}}
{{$metricsPlaceholders := $method.AdditionalMetricsLabels}}
// {{.Name}} ...
func (s *instrumentingMiddleware) {{.Name}}({{joinFullVariables .Args ","}}) ({{joinFullVariables .Results ","}}) {
	defer func(startTime time.Time) {
{{range $from, $to := $metricsPlaceholders}}
	{{if in $method.AdditionalMetricsLabels $to.Name}}
		{{if $to.IsPointer}}
					var _{{$to.Name}} string
					if {{$to.Name}} != nil {
						{{if $to.IsString}}
							_{{$to.Name}} = *{{$to.Name}}
						{{end}}
						{{if $to.IsInt}}
							_{{$to.Name}} = strconv.Itoa(int(*{{$to.Name}}))
						{{end}}
					} else {
							_{{$to.Name}} = "empty"
					}
		{{end}}
	{{end}}
{{end}}
		labels := []string{
			"method", "{{.Name}}",
			"error", strconv.FormatBool(err != nil),
            	{{range $from, $to := $metricsPlaceholders}}
					{{if in $method.AdditionalMetricsLabels $to.Name}}
						{{if $to.IsPointer}}
							"{{$to.Name}}", _{{$to.Name}},
						{{else}}
							{{if $to.IsString}}
								"{{$to.Name}}", {{$to.Name}},
							{{end}}
							{{if $to.IsInt}}
								"{{$to.Name}}", strconv.Itoa(int({{$to.Name}})),
							{{end}}
						{{end}}
					{{end}}
				{{end}}
		}
		s.reqCount.With(labels...).Add(1)
		s.processedBatch.With(labels...).Add(float64(cnt))
		s.reqDuration.With(labels...).Observe(time.Since(startTime).Seconds())
	}(time.Now())
	return s.svc.{{.Name}}({{joinVariableNamesWithEllipsis .Args ","}})
}
{{end}}

// NewInstrumentingMiddleware ...
func NewInstrumentingMiddleware(
	reqCount       metrics.Counter,
	processedBatch metrics.Counter,
	reqDuration    metrics.Histogram,
	svc            {{ .Iface.Name }},
) {{ .Iface.Name }} {
	return &instrumentingMiddleware{
		reqCount:       reqCount,
		reqDuration:    reqDuration,
		processedBatch: processedBatch,
		svc:            svc,
	}
}
`

// WorkerMetrics ...
type WorkerMetrics struct {
	*template.Template
	filePath []string
	imports  imports
}

// Generate ...
func (s *WorkerMetrics) Generate(info api.Interface) (err error) {
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
	t := template.Must(s.Parse(workerMetricsTpl))
	if err = t.Execute(serverFile, info); err != nil {
		return
	}
	err = s.imports.GoImports(info.AbsOutputPath)
	return
}

// NewWorkerMetrics ...
func NewWorkerMetrics(template *template.Template, filePath []string, imports imports) *WorkerMetrics {
	return &WorkerMetrics{
		Template: template,
		filePath: filePath,
		imports:  imports,
	}
}
