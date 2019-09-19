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
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type queenx struct {
	config        ProjectConfiguration
	config_dir    string
	templates_dir string
	args          []string
	host          string
	node          uint
	reinit        bool
}

func QueenxInit(args []string, host string, node uint, reinit bool) queenx {
	var qx queenx

	qx.config_dir = filepath.Join(os.Getenv("HOME"), ".config", "queenx")
	qx.templates_dir = filepath.Join(qx.config_dir, "templates")

	qx.args = args
	qx.node = node
	qx.reinit = reinit

	return qx
}

func (qx *queenx) load_configuration_file(fname string) error {
	source, err := ioutil.ReadFile(fname)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(source), &qx.config)

	if err != nil {
		return err
	}

	return nil
}

func (qx *queenx) load_configuration() error {
	var err error

	err = qx.load_configuration_file("queenx.yml")

	if err != nil {
		return errors.New(fmt.Sprintf("Couldn't open the project configuration file: %v", err))
	}

	if len(qx.host) != 0 {
		qx.config.Remote.Host = qx.host
	}

	qx.config.Local.Project_name = strings.TrimSpace(qx.config.Local.Project_name)

	qx.config.Remote.Proejcts_path = strings.TrimSpace(qx.config.Remote.Proejcts_path)
	qx.config.Remote.Host = strings.TrimSpace(qx.config.Remote.Host)

	qx.config.Build.Cmd_build = strings.TrimSpace(qx.config.Build.Cmd_build)
	qx.config.Build.Cmd_clean = strings.TrimSpace(qx.config.Build.Cmd_clean)
	qx.config.Build.Cmd_post = strings.TrimSpace(qx.config.Build.Cmd_post)
	qx.config.Build.Cmd_pre = strings.TrimSpace(qx.config.Build.Cmd_pre)

	if len(qx.config.Local.Project_name) == 0 {
		return errors.New("You must specify the project name")
	}

	if len(qx.config.Local.Project_dirs) == 0 && len(qx.config.Local.Project_files) == 0 {
		return errors.New("You must specify the project dirs or/and files")
	}

	if len(qx.config.Remote.Proejcts_path) == 0 {
		return errors.New("You must specify the remote projects path")
	}

	if len(qx.config.Remote.Host) == 0 {
		return errors.New("You must specify the remote host")
	}

	//Удалим слеш в начале и конце
	for path_key, path := range qx.config.Local.Project_dirs {
		qx.config.Local.Project_dirs[path_key] = strings.Trim(path, "/")
	}

	for path_key, path := range qx.config.Local.Project_files {
		qx.config.Local.Project_files[path_key] = strings.Trim(path, "/")
	}

	return nil
}

func (qx *queenx) template_unpack(template string, prj_name string) error {
	var gz bool
	var fname = filepath.Join(qx.templates_dir, fmt.Sprintf("%s.tar", template))

	Printf(" -- Trying to open \"%s\"\n", fname)

	file, err := os.Open(fname)

	if err != nil {
		fname = fmt.Sprintf("%s.gz", fname)
		Printf(" -- Trying to open \"%s\"\n", fname)

		file, err = os.Open(fname)

		if err != nil {
			return errors.New(fmt.Sprintf("Couldn't open template file: %v", err))
		}

		gz = true
	} else {
		gz = false
	}

	var tr *tar.Reader
	var dst = "./"

	dst = filepath.Join(dst, prj_name)

	err = os.Mkdir(dst, 0755)

	if err != nil {
		return errors.New(fmt.Sprintf("Couldn't create the project directory: %v", err))
	}

	if gz {
		gzr, err := gzip.NewReader(file)

		if err != nil {
			return errors.New(fmt.Sprintf("Couldn't open template file: %v", err))
		}
		defer gzr.Close()

		tr = tar.NewReader(gzr)
	} else {
		tr = tar.NewReader(file)
	}

	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			Println(" -- OK")
			return nil

		case err != nil:
			return err

		case header == nil:
			continue
		}

		target := filepath.Join(dst, header.Name)

		Printf(" -- [%s]: ", target)

		switch header.Typeflag {
		case tar.TypeDir:
			{
				if _, err := os.Stat(target); err != nil {
					if err := os.MkdirAll(target, 0755); err != nil {
						return err
					}
				}

				Println("OK")

			}
		case tar.TypeReg:
			{
				f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
				if err != nil {
					return err
				}

				if _, err := io.Copy(f, tr); err != nil {
					return err
				}

				f.Close()

				Println("OK")
			}
		}
	}

	return nil
}

func (qx *queenx) Run() error {
	if is_path_exists(qx.config_dir) == false {
		err := os.MkdirAll(qx.config_dir, 0755)

		if err != nil {
			return errors.New(fmt.Sprintf("Couldn't create configuration diretory: %v", err))
		}

		err = os.Mkdir(qx.templates_dir, 0755)

		if err != nil {
			return errors.New(fmt.Sprintf("Couldn't create templates diretory: %v", err))
		}
	}

	switch qx.args[0] {
	case "new":
		{
			if len(qx.args) < 2 {
				return errors.New("You must specify the template name")
			}

			if len(qx.args) < 3 {
				return errors.New("You must specify the project name")
			}

			return qx.template_unpack(qx.args[1], qx.args[2])
		}
	default:
		{
			var err error

			err = qx.load_configuration()

			if err != nil {
				return err
			}

			var prj = ProjectInit(&qx.config)

			switch qx.args[0] {
			case "build":
				{
					err = prj.Init(qx.reinit)

					if err != nil {
						return err
					}

					err = prj.Build()
				}
			case "clean":
				{
					err = prj.Clean()
				}
			case "init":
				{
					err = prj.Init(qx.reinit)
				}
			case "run":
				{
					prj.Run(qx.args[1:], qx.node)
				}
			}

			if err != nil {
				return err
			}
		}
	}

	return nil
}
