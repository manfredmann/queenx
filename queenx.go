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
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

func load_configuration() (Configuration, error) {
	var config Configuration

	source, err := ioutil.ReadFile("queenx.yml")

	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal([]byte(source), &config)

	if err != nil {
		return config, err
	}

	return config, nil
}

func main() {
	config, err := load_configuration()

	if err != nil {
		fmt.Printf("Couldn't open the project configuration file: %v\n", err)
		os.Exit(1)
	}

	flag.Usage = func() {
		fmt.Printf("Usage: %s [-h hostname] action\n\n", os.Args[0])
		fmt.Println("  action")
		fmt.Println("    init   - Init the directory structure")
		fmt.Println("    build  - Build the project")
		fmt.Println("    run    - Run the app")
		fmt.Println("")
		flag.PrintDefaults()
	}

	host_ptr := flag.String("h", config.Remote.Host, "Host name")
	node_ptr := flag.Uint("n", 0, "Run on the node")
	reinit_ptr := flag.Bool("r", false, "Reinit (the directory structure on remote host will be removed)")

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(0)
	}

	config.Remote.Host = *host_ptr

	config.Local.Project_name = strings.TrimSpace(config.Local.Project_name)
	config.Remote.Proejcts_path = strings.TrimSpace(config.Remote.Proejcts_path)
	config.Remote.Host = strings.TrimSpace(config.Remote.Host)

	if len(config.Local.Project_name) == 0 {
		fmt.Println("You must specify the project name")
		os.Exit(0)
	}

	if len(config.Local.Project_dirs) == 0 && len(config.Local.Project_files) == 0 {
		fmt.Println("You must specify the project dirs or/and files")
		os.Exit(0)
	}

	if len(config.Remote.Proejcts_path) == 0 {
		fmt.Println("You must specify the remote projects path")
		os.Exit(0)
	}

	if len(config.Remote.Host) == 0 {
		fmt.Println("You must specify the remote host")
		os.Exit(0)
	}

	var prj = ProjectInit(&config)

	switch flag.Arg(0) {
	case "build":
		{
			prj.Init(*reinit_ptr)
			prj.Build()
		}
	case "clean":
		{
			prj.Clean()
		}
	case "init":
		{
			prj.Init(*reinit_ptr)
		}
	case "run":
		{
			prj.Run(flag.Args()[1:], *node_ptr)
		}
	}
}
