**Requirements:**


As the company provides PAAS platform operation for other business departments, but some colleagues in business departments are not used to using web pages, so we have to provide them with background kubectl terminal to manage their business ns.


The leader's solution is to manage the ns of the business by generating individual users through the UA on the k8s, and provide a container to manage them. Therefore, each container operation and maintenance colleague needs to manually create UA, docker and other operations on the command line, which seems to create a wheel. Then I want to directly write a command tool to replace these manual operations




**Kubectl automatic container creation tool design:**


1. write a command line tool

2. the config file can be automatically transferred in the command line tool

3. enter container name

4. input and use image

5. replace the config file in the container

6. realize k8s context switching in the container


**需求：**

由于公司专门为其他业务部门提供了 paas 平台操作，但是有的业务部门的同事不太习惯使用 web 页面，所以我们不得不为他们提供后台的 kubectl 终端对他们的业务 NS 进行管理。

而领导的解决方式是通过在 K8S 上通过 UA 生成单独的用户来实现对业务的 NS 进行管理，并提供一个容器来实现对其的管理，所以每次容器运维同事都需要通过手动的方式在命令行创建 UA、docker 等操，这就显得比较造轮子，那么我就想通过直接编写一个命令工具来替代这些手动操作的问题



**kubectl 自动创建容器工具设计：**

1. 编写一个命令行工具
2. 在该命令行工具中能够自动传递 config 文件
3. 输入容器 name
4. 输入使用 image
5. 在容器中替换 config 文件
6. 实现对容器中做 K8S 上下文切换
