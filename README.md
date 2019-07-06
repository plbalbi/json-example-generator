# Random JSON example generator

### Installation steps
This guide is assuming the OS used is linux-based, and build-essentials (gcc, etc.) is installed.
1- Install the go-environment. Here is a tutorial for reference, download links, and support links. https://golang.org/doc/install
2- Once installed, check in a terminal-emulator if go was correctly setuped:
    - Run the command ```go version```, you should see your go installation version.
    - Run echo ```echo $GOPATH```, you should see sth. like ```~/some_username/go```.
3- Clone this repository, by using the ```go get github.com/plbalbi/json-example-generator``` command.
4- Go into the cloned-project directory. This should do it ```cd $HOME/go/src/github.com/plbalbi/json-example-generator```.
5- Run the installation script ```./install.sh```. This will download the project go dependencies, and compile the parser.
6- Run the tests to check evrything is working: ```go test ./...```.
7- Compile the binary, and use it: ```go build```.

Tool that generates random JSON examples, with a structure defined in a subset of GoLang struct syntax.

## Tooling
So far, the tooling that will be used is:
- GoYacc: https://about.sourcegraph.com/go/gophercon-2018-how-to-write-a-parser-in-go
- A build by hand Lexer, inspired in https://www.youtube.com/watch?v=HxaD_trXwRE
