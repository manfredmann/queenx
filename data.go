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

type ProjectConfiguration struct {
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
		Cmd_clean string `yaml:"cmd_clean"`
	} `yaml:"build"`
	Run struct {
		Bin_path string `yaml:"bin_path"`
		Bin_name string `yaml:"bin_name"`
	} `yaml:"run"`
}

type QueenxConfiguration struct {
	Tools struct {
		Rsync_args     []string `yaml:"rsync_args"`
		SSH_Build_args []string `yaml:"ssh_build_args"`
		SSH_Run_args   []string `yaml:"ssh_run_args"`
	} `yaml:"tools"`
}

var bin_ssh = "ssh"
var bin_rsync = "rsync"
