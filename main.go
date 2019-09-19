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
