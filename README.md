# Gentee script programming language

Gentee is a free open source script programming language. The Gentee programming language is designed to create scripts to automate repetitive actions and processes on your computer. If you use or plan to use .bat files, bash, PowerShell or special programs to automate actions, then try doing the same thing with Gentee. 

At this moment, Gentee is **under construction** but all current tests are successful.

## Documentation

Browse the [online documentation here](https://gentee.github.io). It describes features that have already been realized.

## How to run Gentee scripts

* Build the *gentee* executable file from *cli/gentee.go* using go compiler.
```
cd ...github.com/gentee/cli/gentee
go build
```
* Specify the script file when running *gentee*. The script file can have any extension.
```
Linux: ./gentee scipts/myscript.g 
Wndows: gentee.exe scipts/myscript.g
```
* Also, you can associate the *gentee* program with script files in your operating system.

### Gentee compiler

```gentee [-t] <scriptname> (<scriptname>)```

By default, the program prints the output of the script to the console and returns 0 if successful.

#### Command line parameters

* **scriptname** - full or relative path to the script file. If several scripts are specified, then they will be executed sequentially.
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
