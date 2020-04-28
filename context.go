package helper

import (
	"github.com/gernest/alien"
	"net/http"
	"net/url"
	"strings"
	"reflect"
	"strconv"
	"mime/multipart"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"fmt"
)

const (
	defaultMemory = 32 << 20 // 32 MB
	indexPage = "index.html"
	headerContentType = "Content-Type"
	headerContentDisposition = "Content-Disposition"
	mimeMultipartForm = "multipart/form-data"
	mimeTextPlainCharsetUTF8 = "text/plain;charset=UTF-8"
	mimeApplicationJSONCharsetUTF8 = "application/json;charset=UTF-8"
)

type Context struct {
	w  http.ResponseWriter
	r *http.Request
	p  alien.Params
	q  url.Values
}

func NewHttpContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{w:w, r:r}
}

func (c *Context) Request() *http.Request {
	return c.r
}

func (c *Context) Response() http.ResponseWriter {
	return c.w
}

func (c *Context) Param(name string) string {
	if c.p == nil {
		c.p = alien.GetParams(c.r)
	}
	if c.p == nil {
		return ""
	}
	return c.p.Get(name)
}

// read values from path, query string or form, store them in a struct specified by vals.
// param name is specified by field tag. e.g.
//
// var vals struct {
//    V1 int    `path:"v1"`   // read path param ":v1", int8,uint8, ..., int64, uint64 are acceptable
//    V2 bool   `query:"v2"`  // read query param "v2=xxx"
//    V3 string `form:"v3"`   // read form param "v3=xxx"
// }
// if status, err := c.ReadParams(&vals); err != nil {
//    c.Error(status, err.Error())
//    return
// }
func (c *Context) ReadParams(vals interface{}) (status int, err error) {
	if vals == nil {
		return http.StatusOK, nil
	}
	p := reflect.ValueOf(vals)
	if p.Kind() != reflect.Ptr {
		return http.StatusInternalServerError, fmt.Errorf("vals must be pointer")
	}
	v := p.Elem() // struct Value
	if v.Kind() != reflect.Struct {
		return http.StatusInternalServerError, fmt.Errorf("vals must be pointer of struct")
	}
	t := v.Type() // struct Type
	n := t.NumField()
	var val string
	for i:=0; i<n; i++ {
		field := t.Field(i) // StructField
		if tag, ok := field.Tag.Lookup("path"); !ok {
			if tag, ok = field.Tag.Lookup("query"); !ok {
				if tag, ok = field.Tag.Lookup("form"); !ok {
					continue
				} else {
					val = c.FormValue(tag)
				}
			} else {
				val = c.QueryParam(tag)
			}
		} else {
			val = c.Param(tag)
		}
		if len(val) == 0 {
			continue
		}

		fv := v.Field(i) // field Value
		ft := field.Type // field Type
		switch ft.Kind() {
		case reflect.String:
			fv.Set(reflect.ValueOf(val))
		case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64:
			i, err := strconv.ParseInt(val, 10, ft.Bits())
			if err != nil {
				return http.StatusBadRequest, err
			}
			fv.Set(reflect.ValueOf(i).Convert(ft))
		case reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
			i, err := strconv.ParseUint(val, 10, ft.Bits())
			if err != nil {
				return http.StatusBadRequest, err
			}
			fv.Set(reflect.ValueOf(i).Convert(ft))
		case reflect.Float64,reflect.Float32:
			f, err := strconv.ParseFloat(val, ft.Bits())
			if err != nil {
				return http.StatusBadRequest, err
			}
			fv.Set(reflect.ValueOf(f).Convert(ft))
		case reflect.Bool:
			b, err := strconv.ParseBool(val)
			if err != nil {
				return http.StatusBadRequest, err
			}
			fv.Set(reflect.ValueOf(b))
		default:
			return http.StatusNotImplemented, fmt.Errorf("value of type %s not implemented", ft.Name())
		}
	}
	return http.StatusOK, nil
}

func (c *Context) QueryParam(name string) string {
	if c.q == nil {
		c.q = c.r.URL.Query()
	}
	return c.q.Get(name)
}

func (c *Context) GetQueryArray(name string) ([]string, bool) {
	if c.q == nil {
		c.q = c.r.URL.Query()
	}
	vals, ok := c.q[name]
	return vals, ok
}

func (c *Context) GetQueryParam(name string) (string, bool) {
	if c.q == nil {
		c.q = c.r.URL.Query()
	}
	if vals, ok := c.q[name]; ok {
		return vals[0], true
	} else {
		return "", false
	}
}

func (c *Context) QueryParams() url.Values {
	if c.q == nil {
		c.q = c.r.URL.Query()
	}
	return c.q
}

func (c *Context) QueryString() string {
	return c.r.URL.RawQuery
}

func (c *Context) FormValue(name string) string {
	return c.r.FormValue(name)
}

func (c *Context) FormParams() (url.Values, error) {
	if strings.HasPrefix(c.r.Header.Get(headerContentType), mimeMultipartForm) {
		if err := c.r.ParseMultipartForm(defaultMemory); err != nil {
			return nil, err
		}
	} else {
		if err := c.r.ParseForm(); err != nil {
			return nil, err
		}
	}
	return c.r.Form, nil
}

func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	_, fh, err := c.r.FormFile(name)
	return fh, err
}

func (c *Context) MultipartForm() (*multipart.Form, error) {
	err := c.r.ParseMultipartForm(defaultMemory)
	return c.r.MultipartForm, err
}

func (c *Context) Cookie(name string) (*http.Cookie, error) {
	return c.r.Cookie(name)
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.w, cookie)
}

func (c *Context) Cookies() []*http.Cookie {
	return c.r.Cookies()
}

func (c *Context) Header(name string) string {
	return c.r.Header.Get(name)
}

func (c *Context) SetHeader(key, value string) {
	c.w.Header().Set(key, value)
}

func (c *Context) AddHeader(key, value string) {
	c.w.Header().Add(key, value)
}

func (c *Context) writeContentType(contentType string) {
	c.w.Header().Set(headerContentType, contentType)
}

func (c *Context) Blob(code int, contentType string, b []byte) (err error) {
	c.writeContentType(contentType)
	c.w.WriteHeader(code)
	_, err = c.w.Write(b)
	return
}

func (c *Context) String(code int, s string) error {
	return c.Blob(code, mimeTextPlainCharsetUTF8, []byte(s))
}

func (c *Context) json(code int, i interface{}, indent string) error {
	enc := json.NewEncoder(c.w)
	if indent != "" {
		enc.SetIndent("", indent)
	}
	c.writeContentType(mimeApplicationJSONCharsetUTF8)
	c.w.WriteHeader(code)
	return enc.Encode(i)
}

func (c *Context) JSON(code int, i interface{}) error {
	return c.json(code, i, "")
}

func (c *Context) JSONBlob(code int, b []byte) (err error) {
	return c.Blob(code, mimeApplicationJSONCharsetUTF8, b)
}

func (c *Context) JSONPretty(code int, i interface{}, indent string) (err error) {
	return c.json(code, i, indent)
}

func (c *Context) Stream(code int, contentType string, r io.Reader) (err error) {
	c.writeContentType(contentType)
	c.w.WriteHeader(code)
	_, err = io.Copy(c.w, r)
	return
}

func NotFoundHandler(c *Context) (err error) {
	return c.Error(http.StatusNotFound, "File not found")
}

func (c *Context) Error(code int, msg string) (err error) {
	return c.JSON(code, map[string]interface{}{"code":code, "msg":msg})
}

func (c *Context) File(file string) (err error) {
	f, err := os.Open(file)
	if err != nil {
		return NotFoundHandler(c)
	}
	defer f.Close()

	fi, _ := f.Stat()
	if fi.IsDir() {
		file = filepath.Join(file, indexPage)
		f, err = os.Open(file)
		if err != nil {
			return NotFoundHandler(c)
		}
		defer f.Close()
		if fi, err = f.Stat(); err != nil {
			return
		}
	}
	http.ServeContent(c.w, c.r, fi.Name(), fi.ModTime(), f)
	return
}

func (c *Context) contentDisposition(file, name, dispositionType string) error {
	c.w.Header().Set(headerContentDisposition, fmt.Sprintf("%s;filename=%s", dispositionType, name))
	return c.File(file)
}

func (c *Context) Attachment(file, name string) error {
	return c.contentDisposition(file, name, "attachment")
}

func (c *Context) Inline(file, name string) error {
	return c.contentDisposition(file, name, "inline")
}

func (c *Context) Redirect(code int, url string) error {
	if code < 300 || code > 308 {
		return fmt.Errorf("ErrInvalidRedirectCode %d", code)
	}
	http.Redirect(c.w, c.r, url, code)
	return nil
}

func (c *Context) ReadJSON(res interface{}) (code int, err error) {
	if c.r.Body == nil {
		return http.StatusBadRequest, fmt.Errorf("bad request")
	}
	defer c.r.Body.Close()

	dec := json.NewDecoder(c.r.Body)
	if err = dec.Decode(res); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
