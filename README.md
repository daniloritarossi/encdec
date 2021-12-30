# ENC and DEC STRING (es:PASSWORD)
## _Developed in GOLANG_
[![Danilo Ritarossi](https://media-exp1.licdn.com/dms/image/C5116AQHnrgF1Z-9Wyg/profile-displaybackgroundimage-shrink_200_800/0/1516649190076?e=1644451200&v=beta&t=uejYUnxpt_2lERCRXybdRFr4cRf8mGSMx2Y27EkVNsw)](https://www.linkedin.com/in/daniloritarossi/)


[![License: GPL v3](https://img.shields.io/badge/license-MIT-green)](https://github.com/daniloritarossi/encanddendc/blob/main/LICENSE)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/gomods/athens)
![PowerShell Gallery](https://img.shields.io/powershellgallery/p/DNS.1.1.1.1)
![GO](https://img.shields.io/static/v1?label=Version&message=1.0&color=<COLOR>)






Stand alone software developed in GO Language, allows you to uniquely encrypt the string of the application based on keywords:

- Machine ID
- OS

***It is generally very useful for encrypting passwords***

## Features

- Portability on Windows, MAC and UNIX / LINUX and much more operating systems 
- Crypt and Decrypt string (password) from a shell_exec
- Robust string encryption system in GOLANG
- No further installations are required
- It has no dependencies on the operating system

## Tech

We use a number of open source projects to work properly:

- [GO] - Developed


And of course Danilo itself is open source with a [public repository][daniloritarossi] on GitHub.

## Installation

ENC and DEC STRING requires [public repository][go] v1.17+ installed to run.

If you want compile the code you need to install the dependency with "github.com/denisbrodbeck/machineid"
move to directory installation es: c:\Program Files\Go and put the code

```sh
go get github.com/denisbrodbeck/machineid
```

#### Easy to use

Normally you can use directly the executable **encdec** file already generated and ready to use for LINUX distribution in /src folder of your project  (please visit [Linux Executable v1.0]  ) once loaded into the system you can generate the encrypted password as I describe it

ENCRYPT--> You can run the program from the bash command by running the statement:

```sh
go run main.go ENC PasswordToEncrypt
```

The result is something like:

```sh
encrypted : 26a8d6a84be3c8ae4ef446f32cccc7affcf2d9612a7bc73a19f7821264a8b0f64f
```

DECRYPT--> You can run the program from the bash command by running the statement:

```sh
go run main.go DEC 26a8d6a84be3c8ae4ef446f32cccc7affcf2d9612a7bc73a19f7821264a8b0f64f
```

The result is something like:

```sh
PasswordToEncrypt
```
If you want you can generate an executable code, see section "Building for source"

#### Building for source

For production release move to directory program es: c:\Users\MYNAME\Go\src and past the code:

```sh
env GOOS=target-OS GOARCH=target-architecture go build package-import-path
```

Example is for LINUX release:
```sh
env GOOS=linux GOARCH=amd64 go build main.go
```
You are generating a build for linux distribution called "main" in the same directory,
if you want you can create a package with other name, use this code:

```sh
go build -o <your desired name>
```
The example before is like:

```sh
env GOOS=linux GOARCH=amd64 go build -o encdec main.go 
```

you can run the executable from a SHELL_EXEC in the same way as reported in the "Easy to use" section

you can find all the possible export configurations in the [Table of config contents] section

## Development

Want to contribute? Great!

Make a change in your file and instantaneously see your updates!

## Table of config contents 

The following table shows the possible combinations of GOOS and GOARCH you can use:

| GOOS - Target Operating System       | GOARCH - Target Platform  |
| ------------------------------------ |:-------------------------:| 
|android	|arm|
|darwin	|386|
|darwin	|amd64|
|darwin	|arm|
|darwin	|arm64|
|dragonfly	|amd64|
|freebsd	|386|
|freebsd	|amd64|
|freebsd	|arm|
|linux	|386|
|linux	|amd64|
|linux	|arm|
|linux	|arm64|
|linux	|ppc64|
|linux	|ppc64le|
|linux	|mips|
|linux	|mipsle|
|linux	|mips64|
|linux	|mips64le|
|netbsd	|386|
|netbsd	|amd64|
|netbsd	|arm|
|openbsd	|386|
|openbsd	|amd64|
|openbsd	|arm|
|plan9	|386|
|plan9	|amd64|
|solaris	|amd64|
|windows	|386|
|windows	|amd64|	



## License

GNU General Public License v3.0
The GNU General Public License v3.0 (GPL) — Danilo Ritarossi. Please have a look at the [public repository][LICENSE.md] for more details.


**Free Software, Hell Yeah!**

   [daniloritarossi]: <https://github.com/daniloritarossi>
   [go]: <https://go.dev>   
   [LICENSE.md]: <https://github.com/daniloritarossi/encdec/blob/main/LICENSE>
   [Linux Executable v1.0]: <https://github.com/daniloritarossi/encanddendc/releases/tag/v1.0>
   [Table of config contents]: <https://github.com/daniloritarossi/encdec#table-of-config-contents>