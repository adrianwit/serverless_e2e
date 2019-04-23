package tail

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

type Route struct {
	Expr    string
	regExpr *regexp.Regexp
	Dest    *Destination
}

type Routes []*Route

func (r *Route) RegExpr() (*regexp.Regexp, error) {
	var err error
	if r.regExpr != nil {
		return r.regExpr, err
	}
	r.regExpr, err = regexp.Compile(r.Expr)
	return r.regExpr, err
}

func (r Routes) Match(URL string) (*Route, error) {
	if !strings.Contains(URL, ".json") {
		return nil, nil
	}
	parsedURL, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	for _, candidate := range r {
		compiledExpr, err := candidate.RegExpr()
		if err != nil {
			return nil, fmt.Errorf("failed to compile expr %v for dest: %v, %v", candidate.Expr, candidate.Dest.Table, err)
		}
		matched := compiledExpr.Match([]byte(parsedURL.Path))
		if matched {
			return candidate, nil
		}
	}
	return nil, nil
}
