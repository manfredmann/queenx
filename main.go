/*
* QueenX
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
}

var scp_path = "/usr/bin/scp"
var ssh_path = "/usr/bin/ssh"

func load_configuration() (Configuration, error) {
	var config Configuration

	source, err := ioutil.ReadFile("project.queenx.yml")

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
	fmt.Println("Checking project dir's on remote host...")

	var path_prj = fmt.Sprintf("%s/%s", config.Remote.Proejcts_path, config.Local.Project_name)

	fmt.Printf("[%s]: ", path_prj)

	if remote_check_dir(path_prj, config.Remote.Host) == true {
		fmt.Println("OK")
	} else {
		fmt.Printf("Creating... ")
		err := remote_create_dir(path_prj, config.Remote.Host)

		if err != nil {
			fmt.Printf("Error %v", err)
			return err
		}

		fmt.Println("OK")
	}

	for _, dir := range config.Local.Project_dirs {
		var path = fmt.Sprintf("%s/%s", path_prj, dir)

		fmt.Printf("[%s]: ", path)

		if remote_check_dir(path, config.Remote.Host) == true {
			fmt.Println("OK")
		} else {
			fmt.Printf("Creating... ")
			err := remote_create_dir(path, config.Remote.Host)

			if err != nil {
				fmt.Printf("Error %v", err)
				return err
			}

			fmt.Println("OK")
		}
	}
	return nil
}

func project_build(config Configuration) error {
	var path_prj = fmt.Sprintf("%s/%s", config.Remote.Proejcts_path, config.Local.Project_name)

	fmt.Println("Transferring files to remote host...")

	for _, path := range config.Local.Project_dirs {
		fmt.Printf("[%s]: ", path)

		var path_remote = fmt.Sprintf("root@%s:%s", config.Remote.Host, path_prj)
		var path_local = fmt.Sprintf("./%s", path)

		if _, err := os.Stat(path_local); os.IsNotExist(err) {
			fmt.Println("Skip")
			continue
		}

		fmt.Println("")

		err := remote_transfer(path_local, path_remote, config.Remote.Host)

		if err != nil {
			return err
		}
	}

	for _, file := range config.Local.Project_files {
		fmt.Printf("[%s]: ", file)

		var path_remote = fmt.Sprintf("root@%s:%s/", config.Remote.Host, path_prj)
		var path_local = fmt.Sprintf("./%s", file)

		if _, err := os.Stat(path_local); os.IsNotExist(err) {
			fmt.Println("Skip")
			continue
		}

		fmt.Println("")

		err := remote_transfer(path_local, path_remote, config.Remote.Host)

		if err != nil {
			return err
		}
	}

	var ssh_host = fmt.Sprintf("root@%s", config.Remote.Host)
	var ssh_cmd = fmt.Sprintf("cd %s && make clean && make", path_prj)

	cmd := exec.Command(ssh_path, ssh_host, ssh_cmd)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func main() {
	config, err := load_configuration()

	if err != nil {
		fmt.Printf("Couldn't open project configuration file: %v", err)
		os.Exit(1)
	}

	host_ptr := flag.String("h", config.Remote.Host, "Host name")

	flag.Parse()

	if flag.NArg() == 0 {
		os.Exit(0)
	}

	config.Remote.Host = *host_ptr

	switch flag.Arg(0) {
	case "build":
		{
			project_build(config)
		}
	case "init":
		{
			project_init(config)
		}
	}
}
