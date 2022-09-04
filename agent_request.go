package fastreq

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"mime/multipart"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2/utils"
)

// BodyString sets request body.
func (a *Agent) BodyString(bodyString string) *Agent {
	a.req.SetBodyString(bodyString)

	return a
}

// Body sets request body.
func (a *Agent) Body(body []byte) *Agent {
	a.req.SetBody(body)

	return a
}

func (a *Agent) BodyStream(bodyStream io.Reader, bodySize int) *Agent {
	a.req.SetBodyStream(bodyStream, bodySize)

	return a
}

// JSON sends a JSON request.
func (a *Agent) JSON(v interface{}) *Agent {
	if a.jsonEncoder == nil {
		a.jsonEncoder = json.Marshal
	}

	a.req.Header.SetContentType(MIMEApplicationJSON)

	if body, err := a.jsonEncoder(v); err != nil {
		a.errs = append(a.errs, err)
	} else {
		a.req.SetBody(body)
	}

	return a
}

// XML sends an XML request.
func (a *Agent) XML(v interface{}) *Agent {
	a.req.Header.SetContentType(MIMEApplicationXML)

	if body, err := xml.Marshal(v); err != nil {
		a.errs = append(a.errs, err)
	} else {
		a.req.SetBody(body)
	}

	return a
}

// Form sends form request with body if args is non-nil.
//
// It is recommended obtaining args via AcquireArgs and release it
// manually in performance-critical code.
func (a *Agent) Form(args *Args) *Agent {
	a.req.Header.SetContentType(MIMEApplicationForm)

	if args != nil {
		a.req.SetBody(args.QueryString())
	}

	return a
}

// FormFile represents multipart form file
type FormFile struct {
	// Fieldname is form file's field name
	Fieldname string
	// Name is form file's name
	Name string
	// Content is form file's content
	Content []byte
	// autoRelease indicates if returns the object
	// acquired via AcquireFormFile to the pool.
	autoRelease bool
}

// FileData appends files for multipart form request.
//
// It is recommended obtaining formFile via AcquireFormFile and release it
// manually in performance-critical code.
func (a *Agent) FileData(formFiles ...*FormFile) *Agent {
	a.formFiles = append(a.formFiles, formFiles...)

	return a
}

// SendFile reads file and appends it to multipart form request.
func (a *Agent) SendFile(filename string, fieldname ...string) *Agent {
	content, err := ioutil.ReadFile(filepath.Clean(filename))
	if err != nil {
		a.errs = append(a.errs, err)
		return a
	}

	ff := AcquireFormFile()
	if len(fieldname) > 0 && fieldname[0] != "" {
		ff.Fieldname = fieldname[0]
	} else {
		ff.Fieldname = "file" + strconv.Itoa(len(a.formFiles)+1)
	}
	ff.Name = filepath.Base(filename)
	ff.Content = append(ff.Content, content...)
	ff.autoRelease = true

	a.formFiles = append(a.formFiles, ff)

	return a
}

// SendFiles reads files and appends them to multipart form request.
//
// Examples:
//
//	SendFile("/path/to/file1", "fieldname1", "/path/to/file2")
func (a *Agent) SendFiles(filenamesAndFieldnames ...string) *Agent {
	pairs := len(filenamesAndFieldnames)
	if pairs&1 == 1 {
		filenamesAndFieldnames = append(filenamesAndFieldnames, "")
	}

	for i := 0; i < pairs; i += 2 {
		a.SendFile(filenamesAndFieldnames[i], filenamesAndFieldnames[i+1])
	}

	return a
}

// Boundary sets boundary for multipart form request.
func (a *Agent) Boundary(boundary string) *Agent {
	a.boundary = boundary

	return a
}

// MultipartForm sends multipart form request with k-v and files.
//
// It is recommended obtaining args via AcquireArgs and release it
// manually in performance-critical code.
func (a *Agent) MultipartForm(args *Args) *Agent {
	if a.mw == nil {
		a.mw = multipart.NewWriter(a.req.BodyWriter())
	}

	if a.boundary != "" {
		if err := a.mw.SetBoundary(a.boundary); err != nil {
			a.errs = append(a.errs, err)
			return a
		}
	}

	a.req.Header.SetMultipartFormBoundary(a.mw.Boundary())

	if args != nil {
		args.VisitAll(func(key, value []byte) {
			if err := a.mw.WriteField(utils.UnsafeString(key), utils.UnsafeString(value)); err != nil {
				a.errs = append(a.errs, err)
			}
		})
	}

	for _, ff := range a.formFiles {
		w, err := a.mw.CreateFormFile(ff.Fieldname, ff.Name)
		if err != nil {
			a.errs = append(a.errs, err)
			continue
		}
		if _, err = w.Write(ff.Content); err != nil {
			a.errs = append(a.errs, err)
		}
	}

	if err := a.mw.Close(); err != nil {
		a.errs = append(a.errs, err)
	}

	return a
}
