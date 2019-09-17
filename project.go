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
	"fmt"
	"os"
	"os/exec"
)

type Project struct {
	config      *Configuration
	remote_host string
	remote_path string
}

func ProjectInit(config *Configuration) *Project {
	project := new(Project)

	project.config = config
	project.remote_host = fmt.Sprintf("root@%s", config.Remote.Host)
	project.remote_path = fmt.Sprintf("%s/%s", config.Remote.Proejcts_path, config.Local.Project_name)

	return project
}

func (prj *Project) remote_check_dir(path string) bool {
	check_cmd := fmt.Sprintf("[ -d \"%s\" ];", path)
	cmd := exec.Command(bin_ssh, prj.remote_host, check_cmd)

	err := cmd.Run()

	if err != nil {
		return false
	} else {
		return true
	}
}

func (prj *Project) remote_create_dir(path string) error {
	cmd := exec.Command(bin_ssh, prj.remote_host, "mkdir", path)

	return cmd.Run()
}

func (prj *Project) remote_remove_dir(path string) error {
	cmd := exec.Command(bin_ssh, prj.remote_host, "rm -r", path)

	return cmd.Run()
}

func (prj *Project) remote_transfer(local_path string, remote_path string) error {
	cmd := exec.Command(bin_scp, "-r", local_path, remote_path)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (prj *Project) Init(reinit bool) error {
	if reinit {
		fmt.Printf("\033[1;37m -- Removing the directory structure on remote host...\033[0m")
		if prj.remote_check_dir(prj.remote_path) == true {
			err := prj.remote_remove_dir(prj.remote_path)
			if err == nil {
				fmt.Println("OK\033[0m")
			} else {
				fmt.Printf("Error %v\033[0m", err)
				return err
			}
		} else {
			fmt.Println("Nothing to do\033[0m")
		}
	}

	fmt.Println("\033[1;37m -- Checking the project dirs on remote host...\033[0m")

	fmt.Printf("\033[1;37m -- [%s]: ", prj.remote_path)

	if prj.remote_check_dir(prj.remote_path) == true {
		fmt.Println("OK\033[0m")
	} else {
		fmt.Printf("Creating... ")
		err := prj.remote_create_dir(prj.remote_path)

		if err != nil {
			fmt.Printf("Error %v\033[0m", err)
			return err
		}

		fmt.Println("OK\033[0m")
	}

	for _, dir := range prj.config.Local.Project_dirs {
		var path = fmt.Sprintf("%s/%s", prj.remote_path, dir)

		fmt.Printf("\033[1;37m -- [%s]: ", path)

		if prj.remote_check_dir(path) == true {
			fmt.Println("OK\033[0m")
		} else {
			fmt.Printf("Creating... ")
			err := prj.remote_create_dir(path)

			if err != nil {
				fmt.Printf("Error %v\033[0m", err)
				return err
			}

			fmt.Println("OK\033[0m")
		}
	}
	return nil
}

func (prj *Project) Build() error {
	fmt.Println("\033[1;37m -- Transferring files to remote host...\033[0m")

	for _, path := range prj.config.Local.Project_dirs {
		fmt.Printf("\033[1;37m -- [./%s --> %s/%s]: ", path, prj.remote_path, path)

		var path_remote = fmt.Sprintf("%s:%s", prj.remote_host, prj.remote_path)
		var path_local = fmt.Sprintf("./%s", path)

		if _, err := os.Stat(path_local); os.IsNotExist(err) {
			fmt.Println("Skip\033[0m")
			continue
		}

		fmt.Println("\033[0m")

		err := prj.remote_transfer(path_local, path_remote)

		if err != nil {
			return err
		}
	}

	for _, file := range prj.config.Local.Project_files {
		fmt.Printf("\033[1;37m -- [./%s --> %s/%s]: ", file, prj.remote_path, file)

		var path_remote = fmt.Sprintf("%s:%s/", prj.remote_host, prj.remote_path)
		var path_local = fmt.Sprintf("./%s", file)

		if _, err := os.Stat(path_local); os.IsNotExist(err) {
			fmt.Println("Skip\033[0m")
			continue
		}

		fmt.Println("\033[0m")

		err := prj.remote_transfer(path_local, path_remote)

		if err != nil {
			return err
		}
	}

	var remote_cmd string
	var cmd *exec.Cmd
	var err error

	if len(prj.config.Build.Cmd_pre) != 0 {
		fmt.Println("\033[1;37m -- Prebuild...\033[0m")
		remote_cmd = fmt.Sprintf("cd %s && %s", prj.remote_path, prj.config.Build.Cmd_pre)

		cmd = exec.Command(bin_ssh, "-t", "-o LogLevel=QUIET", prj.remote_host, remote_cmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()

		if err != nil {
			return err
		}
	}

	if len(prj.config.Build.Cmd_build) != 0 {
		fmt.Println("\033[1;37m -- Build...\033[0m")
		remote_cmd = fmt.Sprintf("cd %s && %s", prj.remote_path, prj.config.Build.Cmd_build)

		cmd = exec.Command(bin_ssh, "-t", "-o LogLevel=QUIET", prj.remote_host, remote_cmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()

		if err != nil {
			return err
		}
	}

	if len(prj.config.Build.Cmd_post) != 0 {
		fmt.Println("\033[1;37m -- Postbuild...\033[0m")
		remote_cmd = fmt.Sprintf("cd %s && %s", prj.remote_path, prj.config.Build.Cmd_post)

		cmd = exec.Command(bin_ssh, "-t", "-o LogLevel=QUIET", prj.remote_host, remote_cmd)
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

func (prj *Project) Run(args []string, node uint) {
	var prj_cmd string
	var cmd *exec.Cmd
	var args_str string

	for _, arg := range args {
		args_str += fmt.Sprintf("\"%s\" ", arg)
	}

	if node == 0 {
		prj_cmd = fmt.Sprintf("cd %s/bin && ./%s", prj.remote_path, prj.config.Local.Project_name)
	} else {
		prj_cmd = fmt.Sprintf("cd %s/bin && on -n%d ./%s", prj.remote_path, node, prj.config.Local.Project_name)
	}

	cmd = exec.Command(bin_ssh, "-t", "-t", prj.remote_host, prj_cmd, args_str)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	cmd.Run()
}
