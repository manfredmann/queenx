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
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Project struct {
	config      *ProjectConfiguration
	config_qx   *QueenxConfiguration
	remote_host string
	remote_path string
}

func ProjectInit(config *ProjectConfiguration, config_qx *QueenxConfiguration) *Project {
	project := new(Project)

	project.config = config
	project.config_qx = config_qx
	project.remote_host = fmt.Sprintf("root@%s", config.Remote.Host)
	project.remote_path = filepath.Join(config.Remote.Proejcts_path, config.Local.Project_name)

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
	cmd := exec.Command(bin_ssh, prj.remote_host, "rm -rf", path)

	return cmd.Run()
}

func (prj *Project) remote_transfer(local_path string, remote_path string) error {
	var rsync_args = prj.config_qx.Tools.Rsync_args

	rsync_args = append(rsync_args, local_path)
	rsync_args = append(rsync_args, remote_path)

	cmd := exec.Command(bin_rsync, rsync_args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (prj *Project) Init(reinit bool) error {
	if reinit {
		log.Printf("Removing the directory structure on remote host... ")
		if prj.remote_check_dir(prj.remote_path) == true {
			err := prj.remote_remove_dir(prj.remote_path)
			if err == nil {
				log.Println("OK")
			} else {
				return err
			}
		} else {
			log.Println("Nothing to do")
		}
	}

	log.Println("Checking the directory structure on remote host...")

	log.Printf("[%s]: ", prj.remote_path)

	if prj.remote_check_dir(prj.remote_path) == true {
		log.Println("OK")
	} else {
		log.Printf("Creating... ")
		err := prj.remote_create_dir(prj.remote_path)

		if err != nil {
			return errors.New(fmt.Sprintf("May be the project directory doesn't exists: %v", err))
		}

		log.Println("OK")
	}

	for _, dir := range prj.config.Local.Project_dirs {
		var path = filepath.Join(prj.remote_path, dir)

		log.Printf("[%s]: ", path)

		if prj.remote_check_dir(path) == true {
			log.Println("OK")
		} else {
			log.Printf("Creating... ")
			err := prj.remote_create_dir(path)

			if err != nil {
				return err
			}

			log.Println("OK")
		}
	}
	return nil
}

func (prj *Project) Build() error {
	log.Println("Transferring files to remote host...")

	for _, path := range prj.config.Local.Project_dirs {
		log.Printf("[%s --> %s]: ", path, filepath.Join(prj.remote_path, path))

		var path_remote = fmt.Sprintf("%s:%s", prj.remote_host, prj.remote_path)
		var path_local = path

		if is_path_exists(path_local) == false {
			log.Warningln("Skip")
			continue
		}

		log.Println("")

		err := prj.remote_transfer(path_local, path_remote)

		if err != nil {
			return err
		}
	}

	for _, file := range prj.config.Local.Project_files {
		log.Printf("[%s --> %s]: ", file, filepath.Join(prj.remote_path, file))

		var path_remote = fmt.Sprintf("%s:%s/", prj.remote_host, prj.remote_path)
		var path_local = file

		if is_path_exists(path_local) == false {
			log.Warningln("Skip")
			continue
		}

		log.Println("")

		err := prj.remote_transfer(path_local, path_remote)

		if err != nil {
			return err
		}
	}

	var remote_cmd string
	var cmd *exec.Cmd
	var err error

	log.Println("Prebuild...")
	if len(prj.config.Build.Cmd_pre) != 0 {
		remote_cmd = fmt.Sprintf("cd \"%s\" && %s", prj.remote_path, prj.config.Build.Cmd_pre)

		var ssh_args = prj.config_qx.Tools.SSH_Build_args
		ssh_args = append(ssh_args, prj.remote_host)
		ssh_args = append(ssh_args, remote_cmd)

		cmd = exec.Command(bin_ssh, ssh_args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()

		if err != nil {
			return err
		}
	} else {
		log.Println("Nothing to do")
	}

	log.Println("Build...")
	if len(prj.config.Build.Cmd_build) != 0 {
		remote_cmd = fmt.Sprintf("cd \"%s\" && %s", prj.remote_path, prj.config.Build.Cmd_build)

		var ssh_args = prj.config_qx.Tools.SSH_Build_args
		ssh_args = append(ssh_args, prj.remote_host)
		ssh_args = append(ssh_args, remote_cmd)

		cmd = exec.Command(bin_ssh, ssh_args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()

		if err != nil {
			return err
		}
	} else {
		log.Println("Nothing to do")
	}

	log.Println("Postbuild...")
	if len(prj.config.Build.Cmd_post) != 0 {
		remote_cmd = fmt.Sprintf("cd \"%s\" && %s", prj.remote_path, prj.config.Build.Cmd_post)

		var ssh_args = prj.config_qx.Tools.SSH_Build_args
		ssh_args = append(ssh_args, prj.remote_host)
		ssh_args = append(ssh_args, remote_cmd)

		cmd = exec.Command(bin_ssh, ssh_args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()

		if err != nil {
			return err
		}
	} else {
		log.Println("Nothing to do")
	}

	return nil
}

func (prj *Project) Clean() error {
	var remote_cmd string
	var cmd *exec.Cmd
	var err error

	log.Println("Clean...")

	if len(prj.config.Build.Cmd_clean) != 0 {
		remote_cmd = fmt.Sprintf("cd %s && %s", prj.remote_path, prj.config.Build.Cmd_clean)

		var ssh_args = prj.config_qx.Tools.SSH_Build_args
		ssh_args = append(ssh_args, prj.remote_host)
		ssh_args = append(ssh_args, remote_cmd)

		cmd = exec.Command(bin_ssh, ssh_args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()

		if err != nil {
			return err
		}

	} else {
		log.Println("Nothing to do")
	}

	return nil
}

func (prj *Project) Run(args []string, node uint) {
	var prj_cmd string
	var cmd *exec.Cmd

	var bin_path string
	var bin_name string

	if len(prj.config.Run.Bin_name) == 0 {
		log.Warningln("Couldn't find the binary name in configuration. Using defaults")
		bin_name = prj.config.Local.Project_name
		bin_path = "bin"
	} else {
		bin_path = prj.config.Run.Bin_path
		bin_name = prj.config.Run.Bin_name
	}

	bin_path = filepath.Join(prj.remote_path, bin_path)

	log.Printf("Binary path: %s\n", bin_path)
	log.Printf("Binary name: %s\n", bin_name)

	if node == 0 {
		prj_cmd = fmt.Sprintf("cd \"%s\" && \"./%s\"", bin_path, prj.config.Local.Project_name)
	} else {
		prj_cmd = fmt.Sprintf("cd \"%s\" && on -n%d \"./%s\"", bin_path, node, bin_name)
	}

	var ssh_args = prj.config_qx.Tools.SSH_Run_args

	ssh_args = append(ssh_args, prj.remote_host)
	ssh_args = append(ssh_args, prj_cmd)

	for arg_i, arg := range args {
		args[arg_i] = fmt.Sprintf("\"%s\"", arg)
	}

	ssh_args = append(ssh_args, args...)

	cmd = exec.Command(bin_ssh, ssh_args...)

	if prj.config.Run.Log_output {
		stdout_writer, err := newRunWriter(fmt.Sprintf("%s.stdout.log", bin_name), os.Stdout)

		if err != nil {
			log.Warningf("Couldn't create \"%s.stdout.log\" file: %v. Output to stdout only\n", bin_name, err)

			cmd.Stdout = os.Stdout
		} else {
			cmd.Stdout = stdout_writer
		}

		stderr_writer, err := newRunWriter(fmt.Sprintf("%s.stderr.log", bin_name), os.Stderr)

		if err != nil {
			log.Warningf("Couldn't create \"%s.stderr.log\" file: %v. Output to stderr only\n", bin_name, err)

			cmd.Stderr = os.Stderr
		} else {
			cmd.Stderr = stderr_writer

		}
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	cmd.Stdin = os.Stdin

	cmd.Run()
}
