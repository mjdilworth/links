# links
A small go app to crawl a web site and print out the links to the pages it finds. these are sperated into off site links, resources and pages.


### First we need to setup go.  On a Linux box we can do the following:

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

run this "hello" binary

For more info please see https://golang.org/doc/install

## To get Links running
### clone this repo

from your /go/src/github.com/user_name directory

type

>git clone https://github.com/mjdilworth/links.git

### get a couple of third party packages used in links - good example of how to reuse code - many solved problems.
#### used to validate the seed URL entered on the command line
>go get "github.com/asaskevich/govalidator"
#### used to parse HTML pages
>go get "golang.org/x/net/html"

### the source code can be found in links.go
To run the code
>go run links.go

and to build the binary in the src dir
>go build links.go

and to install the binary into the project bin directory. Fron the prohect src directory
>go install github.com/user_name/links 

and to run the binary (-url switch is optional and can be used to change the seeding page)

links -url http://wiprodigital.com

I really like go and especially how its channels and goroutines make it fairly easy to implement concurrency. this is really useful when you want to simultaneously crawl many web pages.  Using threads in C would be more complicated as there is nothing "out of the box".

Go also has great packages for testing and benchmarking. Unfortunately I didnt get the time to use any of these features.

## Final comments
I spent a lot more than 2 hours on this exercise and the result is pretty.. er basic.. and not really production quality.  I am not really happy with it. I have not written any code for ages and certainly never in go. It took a long time to set up the environment and decide on editor, to do real basic stuff  etc. Its a great language and pretty simple, but to be quick at stuff you do need to use it fairly frequently, or at least i do.

anyway, I like go and will endeavour to become more accomplished with it.
