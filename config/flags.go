package config

import (
	"errors"
	"net/url"
	"strings"
)

var (
	errZeroLengthLabelKey = errors.New("zero-length label key")
)

//URL represents a URI reference
type URL struct{ url.URL }

//Set sets a url
func (u *URL) Set(value string) error {
	parsed, err := url.Parse(value)
	if err != nil {
		return err
	}
	u.URL = *parsed
	return nil
}

// Label Key, value pair used to store free form user-data.
type Label struct {
	Key   string
	Value *string
}

//Labels is a list of label
type Labels []Label

//Set adds a label to the list
func (labels *Labels) Set(value string) error {
	set := func(k, v string) {
		var val *string
		if v != "" {
			val = &v
		}
		*labels = append(*labels, Label{
			Key:   k,
			Value: val,
		})
	}
	e := strings.IndexRune(value, '=')
	c := strings.IndexRune(value, ':')
	if e != -1 && e < c {
		if e == 0 {
			return errZeroLengthLabelKey
		}
		set(value[:e], value[e+1:])
	} else if c != -1 && c < e {
		if c == 0 {
			return errZeroLengthLabelKey
		}
		set(value[:c], value[c+1:])
	} else if e != -1 {
		if e == 0 {
			return errZeroLengthLabelKey
		}
		set(value[:e], value[e+1:])
	} else if c != -1 {
		if c == 0 {
			return errZeroLengthLabelKey
		}
		set(value[:c], value[c+1:])
	} else if value != "" {
		set(value, "")
	}
	return nil
}

func (labels Labels) String() string {
	// super inefficient, but it's only for occassional debugging
	s := ""
	valueString := func(v *string) string {
		if v == nil {
			return ""
		}
		return ":" + *v
	}
	for _, x := range labels {
		if s != "" {
			s += ","
		}
		s += x.Key + valueString(x.Value)
	}
	return s
}
