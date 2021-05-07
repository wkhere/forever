package main

import (
	"bytes"
	"io"
)

type redbufWriter struct {
	llw io.Writer
	buf *bytes.Buffer
}

func newRedbufWriter(w io.Writer) *redbufWriter {
	return &redbufWriter{llw: w, buf: new(bytes.Buffer)}
}

func (w *redbufWriter) Write(p []byte) (n int, err error) {
	return w.buf.Write(p)
}

func (w *redbufWriter) FlushInRed(red bool) (err error) {
	if !red {
		_, err = io.Copy(w.llw, w.buf)
		return err
	}

	_, err = w.llw.Write([]byte("\033[31;1m"))
	if err != nil {
		return err
	}
	_, err = io.Copy(w.llw, w.buf)
	if err != nil {
		return err
	}
	_, err = w.llw.Write([]byte("\033[0m"))
	return err
}
