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
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Parameter to create a container (创建容器参数)",
	Long: `example: kubevw create Name image HostPort
示例: create  容器名 容器使用镜像 宿主机映射端口 
	-f : Copy the host /host/kubeconfig directory as the /root/config file of the container when creating the container
	(创建容器时将主机 /host/kubeconfig 目录拷贝为容器的 /root/config 文件)
	example: kubecmd create Name image HostPort -f /host/kubeconfig 
	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetBool("file")
		if file {
			if len(args) == 0 || len(args) > 5 {
				fmt.Println(`Command syntax error:
			kubevw create -h view help`)
			}
			createCOD_copyfile(args)
		} else {
			if len(args) == 0 || len(args) > 3 {
				fmt.Println(`Command syntax error:
			kubevw create -h view help`)
			}
			CreateCOD(args)
		}

	},
}

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.Flags().BoolP("file", "f", false, "Copy config file to container	")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func cod(args []string) {
	NAME := args[0]
	image := args[1]
	port := args[2]

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	body, err := cli.ContainerCreate(context.TODO(), &container.Config{
		Tty:       true,
		OpenStdin: true,
		Image:     image,
	}, &container.HostConfig{
		PortBindings: nat.PortMap{nat.Port("22/tcp"): []nat.PortBinding{{"0.0.0.0", port}}},
	}, nil, nil, NAME)
	if err != nil {
		panic(err)
	}

	containerID := body.ID
	err = cli.ContainerStart(context.TODO(), containerID, types.ContainerStartOptions{})
	fmt.Printf("%s Container Created Successfully!\n", NAME)
}

// 单独创建容器不 copy config 文件到
func CreateCOD(args []string) {
	cod(args)
}

// 拷贝 config 文件到容器中
func createCOD_copyfile(args []string) {
	cod(args)
	NAME := args[0]
	hostPath := args[3]

	dockerCP := fmt.Sprintf("docker cp %s %s:/root/config", hostPath, NAME)
	// fmt.Println(dockerCP)
	// dockerCP := fmt.Sprintf("docker cp /root/config  test:/root")

	cmd := exec.Command("/bin/bash", "-c", dockerCP)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return
	}
}
