---
permalink: "/getting-started-with-go-guide"
layout: post
title: "Getting started with Go guide"
date: 2020-03-18 07:55:40 UTC
updated: 2020-03-18 07:55:40 UTC
comments: false
summary: "..."
---

After 20 years in software development, you forget what it is to be a beginner. 
Reality hit hard, and I'm currently trying to get started on Erlang, and this 
experience inspired me to write this getting started guide for programmers that 
want to jump into Go.

### Who you are

Maybe you're curious about Go, it's getting a lot of traction these days and is 
on a high trajectory. 

<script type="text/javascript" src="https://ssl.gstatic.com/trends_nrtr/2152_RC02/embed_loader.js"></script> <script type="text/javascript"> trends.embed.renderExploreWidget("TIMESERIES", {"comparisonItem":[{"keyword":"golang","geo":"","time":"2009-02-18 2020-03-18"}],"category":0,"property":""}, {"exploreQuery":"date=2009-02-18%202020-03-18&q=golang","guestPath":"https://trends.google.com:443/trends/embed/"}); </script> 

Maybe you wrote some JavaScript either on the frontend or 
backend and want to learn a new language.

Or you're like me, an enterprise programmer that mainly used C# and Java and at 
some point needed to get some excitement.

No matter this guide will help you pass this initial "what should I do next" 
phase when learning a new language.

### What you probably already heard

If you asked around how to get started in Go, you most certainly were pointed 
to the following resources, **which you should check** before anything else. 

* [A tour of Go](https://tour.golang.org/welcome/1)
* [Go by example](https://gobyexample.com/)

Those will teach you the basics and concepts of Go and get familiar with its 
syntax. My guide assumes you took those because I'll focus on blockers that 
happen when you decide to create your first project without hand-holding, and 
you're not sure on how to proceed.

### No npm init or rails new

The first aspect I immediately loved about Go was the simplicity of files and 
directory structure when creating a new package. Coming from 15 years of C# and 
3 years of Node, it was utterly refreshing to get cleaner directories with Go.

Packages are just directory, and this cannot be simpler. We'll create our first 
web server. For me, the first thing I try to create when learning a new language 
is some web API.

Since there's no tooling needed to create project and no specific directory 
structure enforce by the language, you **cannot break anything**. It's a 
considerable burden to remove for newcomers of a language. Don't worry about 
creating files at the right place and what not, you can change it later and create 
sub-package later.

We'll first create our `main` package. As you probably learn by now the `main` 
package can be executed as an OS process.

```shell
$> mkdir demo && cd demo
```

You can create this package anywhere on your file system. With the Go module, 
we don't need to bother with having a GOPATH environment variable anymore. The 
dependencies of your program are also using specific versions describe in the 
`go.mod` file.

Let's create our main file.

```go
package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", home)
	log.Println(http.ListenAndServe(":8080", nil))
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}
```

We can now init our Go module files by running:

```shell
$> go mod init github.com/dstpierre/demo
```

Where `github.com/dstpierre/demo` is the module name. This step can be delayed 
but it makes things easier when it comes to dependencies management.

We've got a lot to talk about from this code sample. As someone who's learning 
Go, how could you write the code above? You can't unless you read the document 
or follow a tutorial from the internet.

The Go standard library has a lot to offer. Most often than not, you'll find a 
package in there that will do what you need. Here are some standard library 
packages that you'll probably need on week one, and you should check their 
documentation.

* [net/http](https://golang.org/pkg/net/http/): if you're going to write 
anything web related.
* [encoding](https://golang.org/pkg/encoding/): if you're going to manipulate 
data in [JSON](https://golang.org/pkg/encoding/json/), 
[XML](https://golang.org/pkg/encoding/xml/), 
[CSV](https://golang.org/pkg/encoding/csv/) etc.
* [io/ioutil](https://golang.org/pkg/io/ioutil/): for easier IO helpers.
* [strings](https://golang.org/pkg/strings/): for strings manipulation.
* [time](https://golang.org/pkg/time/): for time manipulation.

It's just a quick list that can get you very far, depending on what you're 
looking to build. What I'm trying to convey here is that Go is probably the 
most straightforward language available at the moment.

### It's OK to not have sub-packages at first

Don't stress about having great packages on your first 2-3 projects. Here's 
something that could work for you if you're building or porting a line-of-business 
application.

Separate in a package your data access functionalities say a `data` package. Your 
main package can be a web server or a command-line program. At least you 
would have your data package reusable on other Go packages if you need to share 
the data access layer.

Eventually, you will most certainly create packages based on scope and domain. 
But if it's not clear at first or you're not comfortable, it's OK to have it 
~wrong at first. It **will not ruin your project** I promise.

As you write more and more Go, you'll get used to the scope of packages and 
what they should expose and what not. At first, there are no good reasons to 
be blocked by this aspect. Refactoring the Go code is dead simple. Don't worry 
about getting it wrong or not having a package at first.

Packages and interfaces will come naturally after a couple of months writing Go 
code that does the work needed.

### Write tests and document your code

Another beautiful aspect of Go is how fast the tests run. I'll encourage you to 
write tests as much as you can. Here are the packages that you should know 
about to get you started:

* [testing](https://golang.org/pkg/testing/)
* [http/httptest](https://golang.org/pkg/net/http/httptest/)

The `httptest` package is one of my favorite packages. Since I'm mainly building 
web applications, it's one of the most useful package for me. It allows you to 
test your web handlers. Here's an example testing our `home` function we wrote 
previously.

```go
import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHome(t *testing.T) {
	// this is the request that the browser would do
	req := httptest.NewRequest("GET", "http://localhost:8080/", nil)

	// this is the ResponseWriter our handler will write to
	rec := httptest.NewRecorder()

	// We execute our handler
	home(rec, req)

	// we get the body of the response
	resp := rec.Result()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// assert that we have the correct response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	} else if string(body) != "hello world" {
		t.Errorf("expected body to be hello world, got %s", string(body))
	}
}
```

We can run our tests:

```shell
$> go test
PASS
ok      _/home/dstpierre/projects/osp/demo   
```

One quick and nice option you can use for your test is the `-cover` argument 
which display the code coverage percentage by your tests.

```shell
$> go test -cover
PASS
coverage: 33.3% of statements
ok      _/home/dstpierre/projects/osp/demo 
```

It takes discipline to have a good set of tests. I'm personally struggling with 
this every day. I mainly write code either for my own SaaS or for small 
companies that hired me as a consultant. Often the tests are skipped. They 
should not, and if you're able to take this habit of writing useful tests, your 
programs will be of higher quality and easier to maintain.

Another aspect that will help you with the long-term maintenance of your 
software is the built-in documentation of your packages. You may document your 
code as follow:

```go
// DoSomething returns exactly the same value as the x
// parameter. You use it like this:
//
//  newX := DoSomething(5)
func DoSomething(x int) int {
  return x
}
```

You can look at doocumentation via `go doc` for instance:

```shell
$> go doc DoSomething
func DoSomething(x int) int
    DoSomething returns exactly the same value as the x parameter. You use it
    like this:

        newX := DoSomething(5)
```

Refer to [this documentaiton](https://blog.golang.org/godoc-documenting-go-code) 
to know more about how to document your code.

### Typical development flow

The development flow for writing Go programs is straightforward. Ideally, you 
would write the tests and write your code at the same time, so running the test 
is how you'd make sure your program works.

So it will be like this:

1. Write tests
2. Write the code needed to have the tests pass
3. Run the tests with `go test -cover`

You may skip the tests and run your program with:

```shell
$> go run main.go
$> go run *.go
```

You may build you program if it's not a library like this:

```
$> go build
$> ./executable-name
```

Where `executable-name` is the name of the binary built by `go build`.

Personally I prefer using `make` and `Makefile`.

```
build:
  go build

test:
  go test -cover
```

And I run:

```shell
$> make build
$> make test
```

I'm using `Makefile` because sooner or later there's more to add to the 
different steps. But that's not important at this stage.

What you want to do is have a program up and running and iterate from there—no 
need to install an auto-reload helper for now. Stop the execution of your program, 
modify the code, and re-run it. The goal is to add more complex workflows as 
you go gradually. There's enough learning curve to get used to advance stuff 
in Go like concurrency and interface. Leave the distraction aside and focus on 
becoming used to the language and concept first.

After your first week writing Go code you are ready to read 
[Effective Go](https://golang.org/doc/effective_go.html).

Reading this will ensure you're writing idiomatic Go code and will level-up 
your understanding of the language.

You need to write and read Go code, there's the only way to get better.

* Pick an open source project you like and try to contribute.
* Start a side project and get better as you go.
* Do not hesitate to read Go's standard library when you want to know host they 
did something.

Another great resource if you want to start using SQL databases is the 
[Go database/sql tutorial](http://go-database-sql.org/).

From there, you'll have all the resources you should need to get you started. 
This guide does not show how to write Go code, and there are already too many 
resources for that. Instead, it tries to remove the road bump you might encounter 
during your first contact with writing Go code.

### First deployment

As everything, deployment is a breeze in Go, you build a binary, and you push 
it to your server(s). You can even cross-compile from your OS to other supported 
OSes.

As you get more experience, you might want to use CI/CD to automate the 
deployment. This guide focus on the first month of getting started with Go. You 
don't need to add superfluous processes to deploy your first prototype. Keep it 
simple; it just works.

```shell
$> GOOS=linux GOARCH=amd64 go build
```

This build a binary for Linux 64-bit from either Mac or Windows. if you're 
already on Linux, you just `go build`.

I'll leave you explore the different options for deployment from here. This 
is a topic worth its own article I think.

On that note, I understand that leaving a comfortable language and stack to 
learn something new can be scary. In the last five years, I've learned Elm, Go, 
and now trying to get started with Erlang.

There's also a difference between lurking and playing a new language. For 
instance, I do have Rust and Cargo installed on my Arch Linux. And from time to 
time, I do write the occasional functions. But deciding to write the next side 
project in a new language is challenging, I think Go is the easiest of all 
languages available to date.