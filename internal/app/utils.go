package app

import (
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/dashotv/fae"
)

// QueryString retrieves a string param from the gin request querystring
func QueryString(c echo.Context, name string) string {
	return c.QueryParam(name)
}
func QueryDefaultString(c echo.Context, name, def string) string {
	v := c.QueryParam(name)
	if v == "" {
		return def
	}
	return v
}

// QueryInt retrieves an integer param from the gin request querystring
func QueryInt(c echo.Context, name string) int {
	v := c.QueryParam(name)
	i, _ := strconv.Atoi(v)
	return i
}

// QueryDefaultInt retrieves an integer param from the gin request querystring
// defaults to def argument if not found
func QueryDefaultInteger(c echo.Context, name string, def int) (int, error) {
	v := c.QueryParam(name)
	if v == "" {
		return def, nil
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return def, err
	}

	if n < 0 {
		return def, fae.Errorf("less than zero")
	}

	return n, nil
}

// QueryBool retrieves a boolean param from the gin request querystring
func QueryBool(c echo.Context, name string) bool {
	return c.QueryParam(name) == "true"
}

// QueryParamString retrieves a string param from the gin request querystring
func QueryParamString(c echo.Context, name string) string {
	return c.QueryParam(name)
}

// QueryParamInt retrieves an integer param from the gin request querystring
func QueryParamInt(c echo.Context, name string) int {
	v := c.QueryParam(name)
	i, _ := strconv.Atoi(v)
	return i
}

// QueryParamBool retrieves a boolean param from the gin request querystring
func QueryParamBool(c echo.Context, name string) bool {
	return c.QueryParam(name) == "true"
}

func QueryParamIntDefault(c echo.Context, name string, def string) int {
	param := c.QueryParam(name)
	result, err := strconv.Atoi(param)
	if err == nil && result > 0 {
		return result
	}

	d, err := strconv.Atoi(def)
	if err == nil && d > 0 {
		return d
	}

	return 0
}

func QueryParamFloatDefault(c echo.Context, name string, def string) float64 {
	param := c.QueryParam(name)
	result, err := strconv.ParseFloat(param, 64)
	if err == nil && result > 0 {
		return result
	}

	if def != "" {
		d, err := strconv.ParseFloat(param, 64)
		if err == nil && d > 0 {
			return d
		}
	}

	return 0
}

func QueryParamBoolDefault(c echo.Context, name string, def string) bool {
	param := c.QueryParam(name)
	if param != "" {
		return param == "true"
	}

	if def != "" {
		return def == "true"
	}

	return false
}

func QueryParamStringDefault(c echo.Context, name string, def string) string {
	param := c.QueryParam(name)
	if param != "" {
		return param
	}
	return def
}

// stolen from gin gonic
// H is a shortcut for map[string]any
type H map[string]any

// MarshalXML allows type H to be used with xml.Marshal.
func (h H) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: "",
		Local: "map",
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range h {
		elem := xml.StartElement{
			Name: xml.Name{Space: "", Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

func exists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}

func sumFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", errors.Wrap(err, "failed to open file")
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", errors.Wrap(err, "failed to hash file")
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func sumFiles(source, destination string) (bool, error) {
	type response struct {
		Hash string
		Err  error
	}
	ch_source := make(chan response)
	ch_dest := make(chan response)

	go func() {
		sum, err := sumFile(source)
		ch_source <- response{sum, err}
	}()
	go func() {
		sum, err := sumFile(destination)
		ch_dest <- response{sum, err}
	}()

	r_source := <-ch_source
	r_dest := <-ch_dest

	if r_source.Err != nil {
		return false, errors.Wrap(r_source.Err, "failed to sum source")
	}
	if r_dest.Err != nil {
		return false, errors.Wrap(r_dest.Err, "failed to sum destination")
	}

	return r_source.Hash == r_dest.Hash, nil
}

// https://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file
// Copy copies the contents of the file at srcpath to a regular file at dstpath.
// If dstpath already exists and is not a directory, the function truncates it.
// The function does not copy file modes or file attributes.
func FileCopy(srcpath, dstpath string) (err error) {
	r, err := os.Open(srcpath)
	if err != nil {
		return errors.Wrap(err, "source")
	}
	defer r.Close() // ignore error: file was opened read-only.

	w, err := os.Create(dstpath)
	if err != nil {
		return errors.Wrap(err, "destination")
	}

	defer func() {
		// Report the error from Close, if any,
		// but do so only if there isn't already
		// an outgoing error.
		if c := w.Close(); c != nil && err == nil {
			err = c
		}
	}()

	_, err = io.Copy(w, r)
	return err
}

// FileLink creates hard link, if destination exists and force is true, it will remove the destination before linking
func FileLink(srcpath, dstpath string, force bool) error {
	if err := FileDir(dstpath); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	if exists(dstpath) {
		if !force {
			return fmt.Errorf("destination exists, force false")
		}
		if err := os.Remove(dstpath); err != nil {
			return fmt.Errorf("failed to remove destination: %w", err)
		}
	}
	return os.Link(srcpath, dstpath)
}

func FileDir(file string) error {
	path := filepath.Dir(file)
	return os.MkdirAll(path, 0755)
}

// WithTimeout runs a delegate function with a timeout,
//
// Example: Wait for a channel
//
//	if value, ok := WithTimeout(func()interface{}{return <- inbox}, time.Second); ok {
//	    // returned
//	} else {
//	    // didn't return
//	}
//
// Example: To send to a channel
//
//	_, ok := WithTimeout(func()interface{}{outbox <- myValue; return nil}, time.Second)
//	if !ok {
//	    // didn't send
//	}
func WithTimeout(delegate func() interface{}, timeout time.Duration) (ret interface{}, ok bool) {
	ch := make(chan interface{}, 1) // buffered
	go func() { ch <- delegate() }()
	select {
	case ret = <-ch:
		return ret, true
	case <-time.After(timeout):
	}
	return nil, false
}
