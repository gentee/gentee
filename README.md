# Gentee script programming language

[![Build Status](https://travis-ci.org/gentee/gentee.png)](https://travis-ci.org/gentee/gentee)
[![Go Report Card](https://goreportcard.com/badge/github.com/gentee/gentee)](https://goreportcard.com/report/github.com/gentee/gentee)
[![GoDoc](https://godoc.org/github.com/gentee/gentee?status.svg)](https://godoc.org/github.com/gentee/gentee)

Gentee is a free open source script programming language. The Gentee programming language is designed to create scripts to automate repetitive actions and processes on your computer. If you use or plan to use .bat files, bash, PowerShell or special programs to automate actions, then try doing the same thing with Gentee. 

## Documentation

- [Gentee programming language (English)](https://docs.gentee.org/)
- [Язык программирования Gentee (Russian)](https://ru.gentee.org/)

All documentation is available on [GitHub pages](https://github.com/gentee/gentee.github.io). 

## Download

- [Linux amd64](https://github.com/gentee/gentee/releases/download/v1.4.0/gentee-1.4.0-linux-amd64.zip)
- [Windows amd64](https://github.com/gentee/gentee/releases/download/v1.4.0/gentee-1.4.0-windows-amd64.zip)
- [macOS amd64](https://github.com/gentee/gentee/releases/download/v1.4.0/gentee-1.4.0-darwin-amd64.zip)

You can download other binary distributions for Linux, macOS, Windows [here](https://github.com/gentee/gentee/releases).


## How to run Gentee scripts

* [Download the binary version](https://github.com/gentee/gentee/releases) of Gentee compiler for your operating system or build the *gentee* executable file from *cli/gentee.go* using [go compiler](https://golang.org/dl/).
```
$ go get -u github.com/gentee/gentee
$ cd gentee/gentee/cli
$ go build
```
* Specify the script file when running *gentee*. The script file can have any extension.
```
Linux: ./gentee myscript.g 
Wndows: gentee.exe myscript.g
```
* Also, you can associate the *gentee* program with script files in your operating system.

### Gentee compiler/interpreter

```gentee [-ver] [-t] <scriptname> [command-line parameters for script]```

By default, the program prints the output of the script to the console and returns 0 if successful.

#### Command line parameters

* **scriptname** - full or relative path to the script file. You can specify the command line parameters for the script after the script file name.
* **-ver** - show the current version of Gentee language.
* **-t** - test the script. When using this parameter, the script must have the **result** parameter in the header with the expected value ([example](https://github.com/gentee/gentee/blob/master/test/scripts/ok.g)). In this mode, the program does not output the result of 
the script execution to the console. If the result does not match, an error message is displayed and an error code 4 is returned.

#### Error code

Code | Description
-----|----------
1 | The script file was not found.
2 | Compilation error.
3 | Runtime Error.
4 | The result is erroneous at start with the **-t** parameter.

## Support

If you have any questions and suggestions or would like to help in the development, [add your issue here](https://github.com/gentee/gentee/issues).

## License

[MIT](http://opensource.org/licenses/MIT)

Copyright (c) 2018-present, Alexey Krivonogov
