package util

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

func DecodeToUtf8(content string, encoding string) (string, error) {
	encoding = strings.ToLower(encoding)
	if encoding == "utf8" || encoding == "utf-8" {
		return content, nil
	}

	charset.Lookup(encoding)
	e, _ := charset.Lookup(encoding)
	if e == nil {
		return "", errors.New("encoding=" + encoding + " don't fond")
	}
	s, err := transformString(e.NewDecoder(), content)
	if err != nil {
		return "", errors.New(fmt.Sprintf("%s: decode %q: %v", encoding, content, err))
	}
	return s, nil
}

func EncodeUtf8(content string, encoding string) (string, error) {
	charset.Lookup(encoding)
	e, _ := charset.Lookup(encoding)
	if e == nil {
		return "", errors.New("encoding=" + encoding + " don't fond")
	}
	s, err := transformString(e.NewEncoder(), content)
	if err != nil {
		return "", errors.New(fmt.Sprintf("%s: decode %q: %v", encoding, content, err))
	}
	return s, nil
}

func transformString(t transform.Transformer, s string) (string, error) {
	r := transform.NewReader(strings.NewReader(s), t)
	b, err := ioutil.ReadAll(r)
	return string(b), err
}
