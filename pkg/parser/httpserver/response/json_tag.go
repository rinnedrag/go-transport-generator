package response

import (
	"errors"
	"strings"

	"github.com/rinnedrag/go-transport-generator/pkg/api"
)

var (
	errHTTPJsonTagDidNotSet = "http json tag did not set"
)

type jsonTag struct {
	prefix string
	suffix string
	next   Parser
}

func (t *jsonTag) Parse(info *api.HTTPMethod, firstTag string, tags ...string) (err error) {
	if strings.HasPrefix(firstTag, t.prefix) && strings.HasSuffix(firstTag, t.suffix) {
		if len(tags) == 2 {
			if info.ResponseJSONTags == nil {
				info.ResponseJSONTags = make(map[string]string)
			}
			info.ResponseJSONTags[tags[0]] = tags[1]
			return
		}
		return errors.New(errHTTPJsonTagDidNotSet)
	}
	return t.next.Parse(info, firstTag, tags...)
}

// NewJSONTag ...
func NewJSONTag(prefix string, suffix string, next Parser) Parser {
	return &jsonTag{
		prefix: prefix,
		suffix: suffix,
		next:   next,
	}
}
