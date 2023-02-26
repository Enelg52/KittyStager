<h1 align="center">
    KittyStager
</h1>

<p align="center">
  <a href="" rel="noopener">
 <img width=150px height=150px src="./img/chat.png"> </a>
</p>


KittyStager is a stage 0 C2. It is made of an API, a client and a malware. The purpose of the API is to deliver some 
basic tasks and shellcode for the malware to inject in memory. The client also connects to the API and is used to interact 
with the malware by creating the tasks for the API. The purpose of KittyStager is to drop a simple stage 0 on the target
and to get enough information to adapt the payload to the target. Because the stage 0 is ment to be as stealth as possible,
I will not add some shell exec capabilities. 

**This project is made for educational purpose only. I am not responsible for any damage caused by this project.**


## Features
- A simple cli to interact with the implant
- API :
    - [x] A web server to host your kittens
    - [ ] User agent whitelist to prevent unwanted connections
- Reconnaissance :
    - [x] Hostname, domain, pid, ip...
    - [x] AV or EDR solution
    - [x] Process list
- Encryption :
    - [x] Key exchange with Opaque
    - [x] Chacha20 encryption
- Malware capabilities :
    - [x] Standard injection
    - [ ] ETW patching
- Sandbox :
    - [ ] Check ram
    - [ ] Check a none existing website
- Payload :
    - [x] Raw shellcode
    - [x] Hex shellcode
    - [ ] Dll
    - [ ] PE

Some settings can be changed in the [config](config.yaml) file

## Architecture
The projet is divided in 3 parts : 
- The [client](client)
- The [server](server)
- The [kitten](kitten)

The crypto part, kitten, config, tasks struct are in [internal](internal).  

## Contributing

Pull requests are welcome. Feel free to open an issue if you want to add other features.

## Contact
Enelg#9993 on discord

## Credits
- https://github.com/C-Sto/BananaPhone
- https://github.com/timwhitez/Doge-Gabh
- https://github.com/c-bata/go-prompt
- https://gist.github.com/leoloobeek/c726719d25d7e7953d4121bd93dd2ed3
- https://github.com/BishopFox/sliver/
- https://github.com/alinz/crypto.go/blob/main/chacha20.go
- https://github.com/frekui/opaque
- ... and many others