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

var bin_scp = "/usr/bin/scp"
var bin_ssh = "/usr/bin/ssh"
