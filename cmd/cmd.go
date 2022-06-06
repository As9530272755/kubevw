// Copyright © 2022 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
)

// adduserCmd represents the adduser command
var cmdCmd = &cobra.Command{
	Use:   "cmd",
	Short: "Execute the corresponding command operation in the container (在容器中执行对应命令操作)",
	Long: `example: kubevw cmd ContainerName CMD
	
	Currently, the supported commands are:
	udd (create user named after container)
	kdir (prerequisite: you must first execute the UDD command to create the.Kube directory in the home directory)
	mv (Prerequisite: the UDD command must be executed first. By default, the config file is moved to the business manager directory/ In Kube)
	da (prerequisite: the UDD command must be executed first to authorize the business administrator file)
	kcu (precondition: UDD, kdir, Da commands must be executed first, and k8s context switching is required).mmand must be executed,Create in home directory Kube directory) 
	
	目前支持的命令有：
	udd(创建以容器命名的用户) 
	kdir(前提条件：必须先执行 udd 命令，在家目录中创建 .kube 目录) 
	mv (前提条件：必须先执行 udd 命令,默认将 config 文件移动到业务管理家目录 ./kube 中)
	da(前提条件：必须先执行 udd 命令，授权业务管理员文件) 
	kcu(前提条件：必须先执行 udd、kdir、da 命令，k8s上下文切换)。
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println(`Command syntax error:
			kubevw cmd -h view help`)
		} else {
			CodExec(args)
		}

	},
}

func init() {
	RootCmd.AddCommand(cmdCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// adduserCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// adduserCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// 获取容器 ID
func CodID(name string) string {

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal("cli:", err)
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	cname := fmt.Sprintf("/" + name)
	cID := ""

	// 获取需要执行命令的容器
	for _, v := range containers {
		for _, Cname := range v.Names {
			if Cname == cname {
				cID = v.ID
			}
		}
	}
	return cID
}

// 容器内执行命令
func CodExec(agrs []string) {
	cname := agrs[0]

	cmd := agrs[1]
	if cmd == "udd" {
		cmd = fmt.Sprintf("useradd -s /bin/bash -m %s", cname)
	}

	if cmd == "kdir" {
		cmd = fmt.Sprintf("mkdir /home/%s/.kube/", cname)
	}

	if cmd == "da" {
		cmd = fmt.Sprintf("chown -R %s.%s /home/%s/", cname, cname, cname)
	}

	if cmd == "mv" {
		cmd = fmt.Sprintf("mv /root/config /home/%s/.kube/config", cname)
	}

	if cmd == "kcu" {
		cmd = fmt.Sprintf("kubectl config use-context kubernetes --kubeconfig=/home/%s/.kube/config ", cname)
	}

	cl, err := docker.NewClient("unix:///run/docker.sock")
	if err != nil {
		log.Fatal("cl", err)
	}

	command := []string{"bash", "-c", cmd}

	exec, err := cl.CreateExec(docker.CreateExecOptions{
		AttachStderr: true,
		AttachStdin:  false,
		AttachStdout: true,
		Tty:          false,
		Cmd:          command,
		Container:    CodID(cname),
	})

	startOpts := docker.StartExecOptions{
		Tty:          true,
		RawTerminal:  true,
		Detach:       false,
		ErrorStream:  os.Stderr,
		InputStream:  os.Stdin,
		OutputStream: os.Stdout,
	}

	err = cl.StartExec(exec.ID, startOpts)
	if err != nil {
		log.Fatal("startExec err", err)
	}
}
