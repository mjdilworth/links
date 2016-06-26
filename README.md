# links
A small go app to crawl a web site and print out pages, links off site and resources.


### First to set up go on a Centos base box

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

```
package main

import "fmt"

func main(){
        fmt.Println("Hello world")
}
```

save this file and then run 
>go install github.com/user_name/hello

this will creat the binary in your GOPATH/bin

run this hello binary

### clone this repo

from /go/src/github.com/mjdilworth

type

>git clone https://github.com/mjdilworth/links.git

### get third party packages - good example of how to reuse a lot of code - many solved problems.
#### used to validate the seed URL enteredon the command line
>go get "github.com/asaskevich/govalidator"
#### used to parse HTML pages
>go get "golang.org/x/net/html"

### the source code is found in links.go
To run the code
>go run links.go

and to build the binary in the src dir
>go build links.go

and to install the binary into the project bin directory. Fron the prohect src directory
>go install github.com/user_name/links 

and to run the binary
links -url http://wiprodigital.com




