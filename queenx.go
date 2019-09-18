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
	"errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

type queenx struct {
	config     Configuration
	config_dir string
	args       []string
	host       string
	node       uint
	reinit     bool
}

func QueenxInit(args []string, host string, node uint, reinit bool) queenx {
	var qx queenx

	qx.config_dir = fmt.Sprintf("%s/.config/queenx", os.Getenv("HOME"))

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

func (qx *queenx) Run() error {
	switch qx.args[0] {
	case "new":
		{
			if len(qx.args) < 2 {
				return errors.New("You must specify template name")
			}
			//is_path_exists
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

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [-h hostname] action\n\n", os.Args[0])
		fmt.Println("  action")
		fmt.Println("    new <name>  - Create project from template")
		fmt.Println("    init        - Init the directory structure")
		fmt.Println("    build       - Build the project")
		fmt.Println("    clean       - Clean the project")
		fmt.Println("    run         - Run the app")
		fmt.Println("")
		flag.PrintDefaults()
	}

	host_ptr := flag.String("h", "", "Host name")
	node_ptr := flag.Uint("n", 0, "Run on the node")
	reinit_ptr := flag.Bool("r", false, "Reinit (the directory structure on remote host will be removed)")

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	var qx = QueenxInit(flag.Args(), *host_ptr, *node_ptr, *reinit_ptr)

	var err = qx.Run()

	if err != nil {
		Errorf("Error: %v\n", err)
		os.Exit(1)
	}
}
