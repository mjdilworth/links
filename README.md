# links

## to set up go on a Centos base box

sudo yum install git
sudo yum install golang
mkdir go

Create env file
>sudo vi /etc/profile.d/goenv.sh

>export GOROOT=/usr/lib/golang
>export GOPATH=$HOME/go
>export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

>source /etc/profile.d/goenv.sh

to test installation

create src/github.com/user_name/hello directory, where user_name is your github user

inside hello directory create a hello.go file and insert the following text

`package main

import "fmt"

func main(){
        fmt.Println("Hello world")
}`


save this file and then run 
>go install github.com/mjdilworth/hello

this will creat the binary in your GOPATH/bin

run this hello binary




