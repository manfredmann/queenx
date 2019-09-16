/*
* queenx - CLI tool for building projects for the QNX4 on target machine
* Copyright (C) 2019  Roman Serov <roman@serov.co>
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Configuration struct {
	Local struct {
		Project_name  string   `yaml:"project_name"`
		Project_dirs  []string `yaml:"project_dirs"`
		Project_files []string `yaml:"project_files"`
	} `yaml:"local"`
	Remote struct {
		Host          string `yaml:"host"`
		Proejcts_path string `yaml:"projects_path"`
	} `yaml:"remote"`
	Build struct {
		Cmd_pre   string `yaml:"cmd_pre"`
		Cmd_build string `yaml:"cmd_build"`
		Cmd_post  string `yaml:"cmd_post"`
	} `yaml:"build"`
}

var scp_path = "/usr/bin/scp"
var ssh_path = "/usr/bin/ssh"

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

func remote_check_dir(path string, host string) bool {
	ssh_host := fmt.Sprintf("root@%s", host)

	ssh_dir := fmt.Sprintf("[ -d \"%s\" ];", path)
	cmd := exec.Command(ssh_path, ssh_host, ssh_dir)

	err := cmd.Run()

	if err != nil {
		return false
	} else {
		return true
	}

}

func remote_create_dir(path string, host string) error {
	ssh_host := fmt.Sprintf("root@%s", host)

	cmd := exec.Command(ssh_path, ssh_host, "mkdir", path)

	return cmd.Run()
}

func remote_transfer(local_path string, remote_path string, host string) error {
	cmd := exec.Command(scp_path, "-r", local_path, remote_path)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	return err
}

func project_init(config Configuration) error {
	fmt.Println("\033[1;37m -- Checking the project dirs on remote host...\033[0m")

	var path_prj = fmt.Sprintf("%s/%s", config.Remote.Proejcts_path, config.Local.Project_name)

	fmt.Printf("\033[1;37m -- [%s]: ", path_prj)

	if remote_check_dir(path_prj, config.Remote.Host) == true {
		fmt.Println("OK\033[0m")
	} else {
		fmt.Printf("Creating... ")
		err := remote_create_dir(path_prj, config.Remote.Host)

		if err != nil {
			fmt.Printf("Error %v\033[0m", err)
			return err
		}

		fmt.Println("OK\033[0m")
	}

	for _, dir := range config.Local.Project_dirs {
		var path = fmt.Sprintf("%s/%s", path_prj, dir)

		fmt.Printf("\033[1;37m -- [%s]: ", path)

		if remote_check_dir(path, config.Remote.Host) == true {
			fmt.Println("OK\033[0m")
		} else {
			fmt.Printf("Creating... ")
			err := remote_create_dir(path, config.Remote.Host)

			if err != nil {
				fmt.Printf("Error %v\033[0m", err)
				return err
			}

			fmt.Println("OK\033[0m")
		}
	}
	return nil
}

func project_build(config Configuration) error {
	var path_prj = fmt.Sprintf("%s/%s", config.Remote.Proejcts_path, config.Local.Project_name)

	fmt.Println("\033[1;37m -- Transferring files to remote host...\033[0m")

	for _, path := range config.Local.Project_dirs {
		fmt.Printf("\033[1;37m -- [./%s --> %s/%s]: ", path, path_prj, path)

		var path_remote = fmt.Sprintf("root@%s:%s", config.Remote.Host, path_prj)
		var path_local = fmt.Sprintf("./%s", path)

		if _, err := os.Stat(path_local); os.IsNotExist(err) {
			fmt.Println("Skip\033[0m")
			continue
		}

		fmt.Println("\033[0m")

		err := remote_transfer(path_local, path_remote, config.Remote.Host)

		if err != nil {
			return err
		}
	}

	for _, file := range config.Local.Project_files {
		fmt.Printf("\033[1;37m -- [./%s --> %s/%s]: ", file, path_prj, file)

		var path_remote = fmt.Sprintf("root@%s:%s/", config.Remote.Host, path_prj)
		var path_local = fmt.Sprintf("./%s", file)

		if _, err := os.Stat(path_local); os.IsNotExist(err) {
			fmt.Println("Skip\033[0m")
			continue
		}

		fmt.Println("\033[0m")

		err := remote_transfer(path_local, path_remote, config.Remote.Host)

		if err != nil {
			return err
		}
	}

	var ssh_cmd string
	var cmd *exec.Cmd
	var err error
	var ssh_host = fmt.Sprintf("root@%s", config.Remote.Host)

	if len(config.Build.Cmd_pre) != 0 {
		fmt.Println("\033[1;37m -- Prebuild...\033[0m")
		ssh_cmd = fmt.Sprintf("cd %s && %s", path_prj, config.Build.Cmd_pre)

		cmd = exec.Command(ssh_path, "-t", "-o LogLevel=QUIET", ssh_host, ssh_cmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()

		if err != nil {
			return err
		}
	}

	if len(config.Build.Cmd_build) != 0 {
		fmt.Println("\033[1;37m -- Build...\033[0m")
		ssh_cmd = fmt.Sprintf("cd %s && %s", path_prj, config.Build.Cmd_build)

		cmd = exec.Command(ssh_path, "-t", "-o LogLevel=QUIET", ssh_host, ssh_cmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()

		if err != nil {
			return err
		}
	}

	if len(config.Build.Cmd_post) != 0 {
		fmt.Println("\033[1;37m -- Postbuild...\033[0m")
		ssh_cmd = fmt.Sprintf("cd %s && %s", path_prj, config.Build.Cmd_post)

		cmd = exec.Command(ssh_path, "-t", "-o LogLevel=QUIET", ssh_host, ssh_cmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()

		if err != nil {
			return err
		}
	}

	return nil
}

func project_run(config Configuration, args []string) {
	var ssh_host = fmt.Sprintf("root@%s", config.Remote.Host)
	var path_prj = fmt.Sprintf("%s/%s", config.Remote.Proejcts_path, config.Local.Project_name)
	var prj_cmd = fmt.Sprintf("cd %s/bin && ./%s", path_prj, config.Local.Project_name)

	var args_str string

	for _, arg := range args {
		args_str += fmt.Sprintf("\"%s\" ", arg)
	}

	cmd := exec.Command(ssh_path, "-t", "-t", ssh_host, prj_cmd, args_str)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Run()
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

	switch flag.Arg(0) {
	case "build":
		{
			project_build(config)
		}
	case "init":
		{
			project_init(config)
		}
	case "run":
		{
			project_run(config, flag.Args()[1:])
		}
	}
}
