package main

import (
	"os"
)

type runWriter struct {
	logfile *os.File
	stdfile *os.File
}

func newRunWriter(file string, out *os.File) (*runWriter, error) {
	logfile, err := os.Create(file)

	if err != nil {
		return nil, err
	}

	return &runWriter{logfile, out}, nil
}

func (w *runWriter) Write(p []byte) (int, error) {
	n, err := w.stdfile.Write(p)

	if err != nil {
		return n, err
	}

	n, err = w.logfile.Write(p)

	return n, err
}

func (w *runWriter) Close() error {
	return nil
}
