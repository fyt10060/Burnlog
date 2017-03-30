# Burnlog
## Introduction
A beginner's blog server trial by golang and docker.
Learn server programming from zero

## Update History

### 0.0.1

完成了CI的工作，CI为Continuous Integration即持续集成，很多初学者会选择在本地测试自己的代码，但是如果仅仅需要一个git push就可以部署到自己的服务器上，然后开始使用那不也可以直接在服务器上测试么？当然，要说速度还是代码运行完直接监听localhost的一个端口测试来得快些。但是这个步骤等到日后有了上线版本之后就变得很有必要了，难道每次更新都还要把代码或者文件夹粘贴到服务器上么？

由于使用了go语言以及docker作为开发环境，加之本项目的开源特性，我使用了[Semaphore](https://semaphoreci.com)作为CI的工具。

当然我这里不会讲如何使用Semaphore，而是希望说一下我通过集成CI环境这个过程对go语言和docker原理的一点理解。

以下的观点只是适合和我一样是后端完全新手，都不知道后端程序如何工作的读者看的。

首先说一下CI的流程:
1. Semaphore会监测我们选定的代码仓库(这里选择的是github)的一个分支
2. 我们在本地写完代码并运行git push指令
3. Semaphore监测到了push的动作，会主动请求把代码拿过来，然后按照我们设定好编译代码先行编译（如果代码结构不复杂例如文件夹下就是main.go则可以直接使用默认代码，如果代码结构稍微多一点，例如我们这个项目，则就需要修改go build命令的参数来让命令可以正确找到main.go）
4. 编译通过，开始部署（如果选择了自动部署的话）
5. 按照部署指令，再次编译代码，然后运行docker build打包docker
6. 将打包好的docker image push到之前选择的docker仓库中
7. 进入服务器，运行已经写好的.sh脚本文件，拉取新镜像并运行docker run命令。
8. 部署完成。

然后来说下我认为的知识点：
1. 首先是go语言，go get获取编译所需的package，然后go build编译出go程序，这是最基本的允许go程序的方法。
2. 其次是docker，首先docker只是一种容器（container）技术，尽管很多人一提容器就会想到docker，但并非docker发明了容器，只是将其简单化了。
3. 容器的基本概念：容器可以类比为一个超小型的虚拟机(VM)，它拥有着能够让后端程序跑起来所需要的所有环境但同时其并不是一个完整的系统，在docker的容器理念中，一个容器仅仅会运行一个程序。至于为什么要学习容器技术，请自行Google。
4. 我们依赖Dockerfile运行的docker build命令，并不是生成了容器，而是生成了一个镜像(image)，通过后续的docker run命令才会生成一个容器。
5. 为什么拉取的新镜像直接run就可以运行程序？因为我们在生成容器的Dockerfile中已经在生成镜像的同时将我们的程序复制到了镜像中同时指定了程序的工作环境，当然我们在docker run命令中可以直接复制程序到容器中。
6. 第五点是依赖Semaphore的预编译环境支持docker命令来做到的，如果选择的CI平台不支持docker环境呢，那么按照我的设想步骤，在部署步骤中可以直接跳过打包docker的步骤，直接进入服务器，而在.sh脚本中加入从git pull代码并编译的步骤，然后加入docker build命令或者认为不需要修改编译环境则直接用docker run命令复制新程序到新容器中并删除旧容器完成部署。
7. 至于Dockerfile中指令的含义，请自行搜索学习。如果有对部署指令集或者服务器端update.sh的需求，我也会考虑后续补上，但其实其中也只是一些最基本的build、run命令而已。