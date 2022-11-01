<h1 align="center">
    KittyStager
</h1>

<p align="center">
  <a href="" rel="noopener">
 <img width=150px height=150px src="./img/chat.png"> </a>
</p>


KittyStager is a simple stage 0 C2. It is made of a web server to host the shellcode and an implant, called kitten.
The purpose of this project is to be able to have a web server and some implant for various usage and be 
able to use it with any shellcode.

## Features
- [x] A simple cli to interact with the implant
- [x] Different implants for various usage
- [x] User agent whitelist to prevent unwanted connections
- [x] A basic recon to get some information about the target (hostname, ip, program files...)
- [x] A Chacha20 encryption to encrypt all the traffic with a none hardcoded key
- [x] Supports hex and bin shellcode
- [x] Shellcode Reflective DLL Injection


***I would not use this project in red team, at least not now.***



## Installation
```
git clone https://github.com/Enelg52/KittyStager.git
cd KittyStager
go build -o KittyStager
```
### Build the implants
#### BasicKitten
```
cd /kitten/basicKitten
go build -o basicKitten.exe
```
#### BasicKitten
```
cd /kitten/bananaKitten
go build -o basicKitten.exe
```
#### DllKitten
```
go build -buildmode=c-shared -o dllKitten.dll
```
### On linux
```
env GOOS=windows GOARCH=amd64 go build -o basicKitten.exe
```

## Usage

```
./KittyStager.exe -h
Usage of kittyStager.exe:
  -p string
        Path to the config file (default "cmd/config/conf.yml")
```
### DllKitten
```
runDll32.exe dllKitten.dll,DllMain
```

## Example
#### Implant check's in
![](/img/KittyCheckin.png)
#### Show target's information
![](/img/KittyTarget.png)
#### Interact with the implant
![](/img/KittyInteract.png)
#### Change sleep time
![](/img/KittySleep.png)
#### Inject shellcode or dll
![img.png](img/KittyShellcode.png)
#### Show recon
![img.png](img/KittyRecon.png)

## Contributing

Pull requests are welcome. Feel free to open an issue if you want to add other features.

## Credits
- https://github.com/C-Sto/BananaPhone
- https://github.com/timwhitez/Doge-Gabh
- https://github.com/c-bata/go-prompt
- https://gist.github.com/leoloobeek/c726719d25d7e7953d4121bd93dd2ed3
- https://github.com/BishopFox/sliver/
- https://github.com/alinz/crypto.go/blob/main/chacha20.go
- ... and many others