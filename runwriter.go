/*
* queenx - CLI tool for building projects for the QNX4 on target machine
* Copyright (C) 2019  Roman Serov <roman@serov.co>
*
* This file is part of queenx.
*
* queenx is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* queenx is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with queenx. If not, see <http://www.gnu.org/licenses/>.
 */

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
