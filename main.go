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
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [-h hostname] [-n node] [-r] action\n\n", os.Args[0])
		fmt.Println("  action")
		fmt.Println("    new <name>  - Create the project from template")
		fmt.Println("    init        - Init the directory structure")
		fmt.Println("    build       - Build the project")
		fmt.Println("    clean       - Clean the project")
		fmt.Println("    run <args>  - Run the app")
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
		Errorf(" -- Error: %v\n", err)
		os.Exit(1)
	}
}
