package xmlgen

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Element struct {
	Name     string
	Attrs    map[string]interface{}
	Contents []interface{}
}

type Elementifiable interface {
	ToElement() *Element
}

func E(name string, attrs map[string]interface{}, contents ...interface{}) *Element {
	return &Element{
		Name:     name,
		Attrs:    attrs,
		Contents: contents,
	}
}

func (e *Element) Marshal(w io.Writer) error {
	err, pathOfError := e.doMarshal(w, []string{})
	if err != nil {
		return errors.New(err.Error() + " (Path: " + strings.Join(pathOfError, " > ") + ")")
	}
	return nil
}

func NoAttrs() map[string]interface{} {
	return map[string]interface{}{}
}

func (e *Element) doMarshal(w io.Writer, curPath []string) (err error, pathOfError []string) {
	// We want the path to bubble up from the recursive call to doMarshal, so
	// we let it and simply empty it if returning with no error.
	defer func() {
		if err == nil {
			pathOfError = []string{}
		}
	}()

	curPath = append(curPath, e.Name)
	pathOfError = curPath

	if !validName(e.Name) {
		err = errors.New("Invalid name for tag: " + e.Name)
		return
	}

	if _, err = w.Write([]byte("<" + e.Name)); err != nil {
		return
	}

	for k, v := range e.Attrs {
		if !validName(k) {
			err = errors.New("Invalid name for attribute: " + k)
			return
		}
		if _, err = w.Write([]byte(" ")); err != nil {
			return
		}
		if err = xml.EscapeText(w, []byte(k)); err != nil {
			return
		}
		if _, err = w.Write([]byte("=\"")); err != nil {
			return
		}
		if err = writeEscaped(w, v); err != nil {
			return
		}
		if _, err = w.Write([]byte("\"")); err != nil {
			return
		}
	}

	if _, err = w.Write([]byte(">")); err != nil {
		return
	}

	if e.Contents != nil {
		for _, c := range e.Contents {
			switch c.(type) {
			case *Element:
				if err, pathOfError = c.(*Element).doMarshal(w, curPath); err != nil {
					return
				}
			default:
				el, ok := c.(Elementifiable)
				if ok {
					if err, pathOfError = el.ToElement().doMarshal(w, curPath); err != nil {
						return
					}
				}
				if err = writeEscaped(w, c); err != nil {
					return
				}
			}
		}
	}

	_, err = w.Write([]byte("</" + e.Name + ">"))
	return
}

func writeEscaped(w io.Writer, v interface{}) error {
	switch v.(type) {
	case bool:
		if v.(bool) {
			if _, err := w.Write([]byte("true")); err != nil {
				return err
			}
		} else {
			if _, err := w.Write([]byte("false")); err != nil {
				return err
			}
		}
	case string:
		return xml.EscapeText(w, []byte(v.(string)))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		_, err := w.Write([]byte(fmt.Sprintf("%d", v)))
		return err
	case float32, float64:
		_, err := w.Write([]byte(fmt.Sprintf("%f", v)))
		return err
	default:
		// By default, defer to xml.Marshal. Bad idea?
		bs, err := xml.Marshal(v)
		if err != nil {
			return errors.New(fmt.Sprintf("Unable to write: %#v", v))
		}
		if _, err := w.Write(bs); err != nil {
			return err
		}
	}
	return nil
}

var validNameRe = regexp.MustCompile("^[:A-Z_a-z\u00C0\u00D6\u00D8-\u00F6\u00F8-\u02ff\u0370-\u037d\u037f-\u1fff\u200c\u200d\u2070-\u218f\u2c00-\u2fef\u3001-\ud7ff\uf900-\ufdcf\ufdf0-\ufffd\u10000-\uEFFFF][-.:A-Z_a-z\u00C0\u00D6\u00D8-\u00F6\u00F8-\u02ff\u0370-\u037d\u037f-\u1fff\u200c\u200d\u2070-\u218f\u2c00-\u2fef\u3001-\ud7ff\uf900-\ufdcf\ufdf0-\ufffd0-9\u00b7\u0300-\u036f\u203f-\u2040]*$")

func validName(s string) bool {
	return validNameRe.MatchString(s)
}
