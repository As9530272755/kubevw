**Tool introduction:**
Implement automatic container creation and external container command execution tool through go code writing

**工具介绍:**
通过 go 编写代码实现自动创建容器及外部执行容器命令工具

**Requirements:**


As the company provides PAAS platform operation for other business departments, but some colleagues in business departments are not used to using web pages, so we have to provide them with background kubectl terminal to manage their business ns.


The leader's solution is to manage the ns of the business by generating individual users through the UA on the k8s, and provide a container to manage them. Therefore, each container operation and maintenance colleague needs to manually create UA, docker and other operations on the command line, which seems to create a wheel. Then I want to directly write a command tool to replace these manual operations





**Kubectl automatic container creation tool design:**


1. the kubeconfig file can be automatically obtained in the command line tool

2. enter container name

3. input and use image

4. replace the config file in the container

5. implement k8s context switching in the container

6. implement external command operation on the inside of the container

**需求：**

由于公司专门为其他业务部门提供了 paas 平台操作，但是有的业务部门的同事不太习惯使用 web 页面，所以我们不得不为他们提供后台的 kubectl 终端对他们的业务 NS 进行管理。

而领导的解决方式是通过在 K8S 上通过 UA 生成单独的用户来实现对业务的 NS 进行管理，并提供一个容器来实现对其的管理，所以每次容器运维同事都需要通过手动的方式在命令行创建 UA、docker 等操，这就显得比较造轮子，那么我就想通过直接编写一个命令工具来替代这些手动操作的问题



**kubectl 自动创建容器工具设计：**

1. 在该命令行工具中能够自动获取 kubeconfig 文件
2. 输入容器 name
3. 输入使用 image
4. 在容器中替换 config 文件
5. 实现对容器中做 K8S 上下文切换
6. 实现外部实现对容器内部的命令操作

**Example:**


1 build code is binary
```bash

$ go build -o kubevw

```




2 view help


```bash

root@consul -3:~/go/src/kube# ./ kubevw -h

The tool provides the following functions:

Automatically create containers and automatically obtain k8s certificates.

Enables a single business administrator to access and operate ns.

The tool provides the following functions:

Automatically create containers and automatically obtain k8s certificates.

Enables a single business administrator to access and operate ns.


Usage:

kubevw [command]


Available Commands:

CMD execute the corresponding command operation in the container

completion Generate the autocompletion script for the specified shell

Create parameter to create a container

help Help about any command

linkkube Parameters for obtaining k8s certificate implementation link


Flags:

--config string config file (default is $HOME/.kube.yaml)

-h, --help help for kubevw

-t, --toggle Help message for toggle


Use "kubevw [command] --help" for more information about a command.

```




3 view the create parameter help


```bash

root@consul -3:~/go/src/kube# ./ kubevw create -h

example: kubevw create Name image HostPort

Example: create container name container uses mirror host mapping port

-f : Copy the host /host/kubeconfig directory as the /root/config file of the container when creating the container

(copy the host /host/kubeconfig directory as the /root/config file of the container when creating the container)

example: kubecmd create Name image HostPort -f /host/kubeconfig


Usage:

kubevw create [flags]


Flags:

-f, --file Copy config file to container

-h, --help help for create


Global Flags:

--config string config file (default is $HOME/.kube.yaml)

```




3 view CMD parameter help


```bash

root@consul -3:~/go/src/kube# ./ kubevw cmd -h

example: kubevw cmd ContainerName CMD


Currently, the supported commands are:

udd (create user named after container)

kdir (prerequisite: you must first execute the UDD command to create the.Kube directory in the home directory)

mv (Prerequisite: the UDD command must be executed first. By default, the config file is moved to the business manager directory/ In Kube)

da (prerequisite: the UDD command must be executed first to authorize the business administrator file)

kcu (precondition: UDD, kdir, Da commands must be executed first, and k8s context switching is required). mmand must be executed,Create in home directory Kube directory)


Currently supported commands are:

UDD (create user named after container)

Kdir (prerequisite: you must first execute the UDD command to create the.Kube directory in the home directory)

MV (prerequisite: the UDD command must be executed first. By default, the config file is moved to the business manager directory. /kube)

Da (prerequisite: the UDD command must be executed first to authorize the business administrator file)

KCU (precondition: UDD, kdir, Da commands must be executed first, and k8s context switching is required).


Usage:

kubevw cmd [flags]


Flags:

-h, --help help for cmd


Global Flags:

--config string config file (default is $HOME/.kube.yaml)

```




4 create containers and copy files to containers


```bash

#Create

root@consul -3:~/go/src/kube# ./ kubevw create test kubectl-cmd:v3 1111 -f /root/test. config

test Container Created Successfully!


#View container created successfully

root@consul -3:~/go/src/kube# docker ps

CONTAINER ID IMAGE COMMAND CREATED STATUS PORTS NAMES

d19c50c4c99b kubectl-cmd:v3 "/bin/sh -c '/usr/sb…" 11 seconds ago Up 10 seconds 0.0.0.0:1111->22/tcp test


#Enter the container to check whether the file is copied in

root@consul -3:~/go/src/kube# docker exec -it test /bin/bash

root@d19c50c4c99b :/# ll /root/config

-rw------- 1 root root 4152 Mar 31 09:20 /root/config

```




5 exit the container to execute CMD parameter function verification


```bash

root@d19c50c4c99b :/# exit


#Execute CMD parameter to use default command

root@consul -3:~/go/src/kube# ./ kubevw cmd test udd

root@consul -3:~/go/src/kube# ./ kubevw cmd test kdir

root@consul -3:~/go/src/kube# ./ kubevw cmd test mv

root@consul -3:~/go/src/kube# ./ kubevw cmd test da

root@consul -3:~/go/src/kube# ./ kubevw cmd test kcu

```




6 enter the container to check whether the corresponding default parameters are successfully executed


```bash

#Enter container

root@consul -3:~/go/src/kube# docker exec -it test /bin/bash


#View test ID

root@d19c50c4c99b :/# id test

uid=1000(test) gid=1000(test) groups=1000(test)


#Check that the test directory has been created Kube file

root@d19c50c4c99b :/# ll /home/test/

total 24

drwxr-xr-x 3 test test 4096 Jun 6 14:56 ./

drwxr-xr-x 1 root root 4096 Jun 6 14:56 ../

-rw-r--r-- 1 test test 220 Apr 4 2018 . bash_ logout

-rw-r--r-- 1 test test 3771 Apr 4 2018 . bashrc

drwxr-xr-x 2 test test 4096 Jun 6 14:56 . kube/

-rw-r--r-- 1 test test 807 Apr 4 2018 . profile


#View config moved

root@d19c50c4c99b :/# ll /home/test/. kube/

total 16

drwxr-xr-x 2 test test 4096 Jun 6 14:56 ./

drwxr-xr-x 3 test test 4096 Jun 6 14:56 ../

-rw------- 1 test test 4152 Mar 31 09:20 config

```
