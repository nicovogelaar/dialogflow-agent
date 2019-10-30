package dialogflow

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Source interface {
	io.ReadCloser
}

type fileSource struct {
	filename string
	buffer   *bytes.Buffer
}

func NewFileSource(filename string) Source {
	return &fileSource{
		filename: filename,
		buffer:   bytes.NewBuffer(nil),
	}
}

func (source *fileSource) Read(p []byte) (n int, err error) {
	if source.buffer.Len() == 0 {
		file, err := os.Open(source.filename)
		if err != nil {
			return 0, fmt.Errorf("open file: %v", err)
		}
		defer func() {
			if closeErr := file.Close(); err == nil {
				err = closeErr
			}
		}()
		if _, err = io.Copy(source.buffer, file); err != nil {
			return 0, fmt.Errorf("copy file: %v", err)
		}
	}

	n, err = source.buffer.Read(p)
	if err != nil {
		return n, err
	}

	if source.buffer.Len() == 0 {
		source.buffer.Reset()
		return n, io.EOF
	}

	return n, err
}

func (source *fileSource) Close() error {
	source.buffer.Reset()
	return nil
}

type urlSource struct {
	url    string
	buffer *bytes.Buffer
}

func NewURLSource(url string) Source {
	return &urlSource{
		url:    url,
		buffer: bytes.NewBuffer(nil),
	}
}

func (source *urlSource) Read(p []byte) (n int, err error) {
	if source.buffer.Len() == 0 {
		var resp *http.Response
		resp, err = http.Get(source.url)
		if err != nil {
			return 0, fmt.Errorf("http get: %v", err)
		}
		defer func() {
			if closeErr := resp.Body.Close(); err == nil {
				err = closeErr
			}
		}()
		if _, err = io.Copy(source.buffer, resp.Body); err != nil {
			if err != nil {
				return 0, fmt.Errorf("copy body: %v", err)
			}
		}
	}

	n, err = source.buffer.Read(p)
	if err != nil {
		return n, err
	}

	if source.buffer.Len() == 0 {
		source.buffer.Reset()
		return n, io.EOF
	}

	return n, err
}

func (source *urlSource) Close() error {
	source.buffer.Reset()
	return nil
}
