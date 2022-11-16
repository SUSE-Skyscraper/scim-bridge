package filters

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type filterOperator int

const (
	Eq filterOperator = iota
	InvalidOperator
)

type filterField int

const (
	Username filterField = iota
	InvalidField
)

type Filter struct {
	FilterField    filterField
	FilterOperator filterOperator
	FilterValue    string
}

func ParseFilter(filterString string) ([]Filter, error) {
	if filterString == "" {
		return []Filter{}, nil
	}

	r := regexp.MustCompile("^(?P<field>\\S+)\\s*(?P<operator>\\w+)\\s*\"(?P<value>\\S+)\"$")
	if !r.MatchString(filterString) {
		return []Filter{}, errors.New("invalid filter")
	}

	match := r.FindStringSubmatch(filterString)
	field, err := parseField(match[1])
	if err != nil {
		return []Filter{}, err
	}
	operator, err := parseOperator(match[2])
	if err != nil {
		return []Filter{}, err
	}

	filter := Filter{
		FilterField:    field,
		FilterOperator: operator,
		FilterValue:    match[3],
	}

	return []Filter{filter}, nil
}

func parseField(field string) (filterField, error) {
	switch field {
	case "userName":
		return Username, nil
	default:
		return InvalidField, errors.New("invalid field")
	}
}

func parseOperator(operator string) (filterOperator, error) {
	switch strings.ToLower(operator) {
	case "eq":
		return Eq, nil
	default:
		return InvalidOperator, errors.New("invalid operator")
	}
}
