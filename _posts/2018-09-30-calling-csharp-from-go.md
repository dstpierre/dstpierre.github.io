---
layout: post
title: "Calling .NET 4.5 C# methods from Go, a WIP"
date: 2018-09-30 11:26:37 UTC
updated: 2018-09-30 11:26:37 UTC
comments: false
summary: "..."
---

> TL:DR Before you get too excited keep in mind that you’ll need to write lots of
> Go structs and some C# trickery to do this. There’s no easy way out. But a
working prototype which is not concurrent friendly yet. Your Go program will
also needs to run on Windows (for now).

![](https://cdn-images-1.medium.com/max/1000/1*dObZ2-MHnfmtUWiqS1vOXA.png)

**My situation**: I’m currently doing some consulting at a 25 years old FinTech
company. While [FinTech](https://en.wikipedia.org/wiki/Financial_technology) is
a buzzword nowadays, I can guarantee it was not 18 years ago, I know because I
was there in that same company in 2000 as a junior developer.

I started writing C# applications in that same company back in the days. And
today I have to *pay* my technical debt, imagine how satiric this is for an
engineer to have to support their code from the past.

Long story short, I’ve been hired mainly to support a lift and shift.

Ho yes! I forgot to mention. The company is hosting their entire infrastructure
on-site, hmm. Anyway, lift and shift the infrastructure to the cloud than start
migrating to cloud-native to repay the technical debt. That’s my mandate.

It means that I will most certainly use lots of Go code. Migrating ~30
applications/processes/background jobs at the same time is impossible. I’ll need
to **keep the current data access layer code**, read lots of C# libraries, as
long as possible, and hopefully be able to call the methods from Go.

Test-driven development and unit tests were not hyper-popular in those days
specifically for small companies / small team like this one. I’m putting my
trust in the C# code that runs for ten years without any update.

### The prototype

I wanted to wait until I had a complete version of this and make the solution
open source. I’ll do that, but I wanted to start talking about this as soon as
possible. Maybe you’re in that situation, and this could give you some head
start on how to do this.

It will be a two-part series, in this part I talk about my first approach.

#### Using stdin and stdout

You’ve read that correctly. By creating a wrapper around the C# libraries that
you need to expose to Go. Using the standard in and standard out, it is possible
to call native, not DotNet-Core, but DotNet 4.5 methods from Go and get the
value returned by the method.

> Shameless plug: Please check out my book on building Go web applications / SaaS.

This is all the code for the prototype, I’ll talk about it after:

Lets start with `Program.cs`

This is a C# .NET console application that continuously read the standard input
for method to execute. The communication between Go and C# uses JSON and Go
needs to send the following object to execute a method:

```json
{
	"typename": "Your.Library.Here.ClassName",
	"functionname": "DoSomething",
	"parameters": ["abc", 1234, bool]
}
```

You need to fully qualify the class that contains the method you want to call.
From there the main C# program tries to load this assembly and will cache it,
see in `DLLLoader.cs`file. You need to add all the libraries as reference that
you’ll use to this C# console application.

The function name is straightforward. And the parameters are simply an array of
object in the order that the function wants to receive the parameters. Code in
`InputParser.cs` tries to handle some data-type parsed via `Newtosoft` JSON
parser mismatch with int64 vs. int. For instance, imagine we have this C#
function inside the Your.Library.Here assembly:

```csharp
namespace Your.Library.Here
{
	public class ClassName
	{
		public List<string> DoSomething(string s, int num, bool flag)
		{
			// do something and return a List<string>
		}
	}
}
```

If the function can be found, it will be executed with the parameters and the
result will be sent back to the standard out via JSON.

From there you can decode the JSON back to a Go struct. This is where you would
need to write lots of structs if you need to send and receive lots of different
data type. Also another important thing to keep in mind, if your functions
accept C# class you’ll need to add code in the `InputParser.cs 's parser to map
incoming object to specific C# class.`

The Go program is dead simple here and only demonstrate the communication
between the C# wrapper and Go.

### It’s far from being production ready!

Yes, of course. It’s a work-in-progress. Obviously, the major issue with this
approach is concurrency. My goal is to call the .NET data access methods from a
Go web application. This cannot work from a simultaneous multi-callers like a
web app in the current form.

I’m working on a better approach that uses UDP connections between the Go
application and the C# wrapper instead of the standard in/out way. It would
handle the concept of “request id” so the returning JSON can be caught by the
original caller via this id. The execution of multiple C# methods simultaneously
would be possible.

Nonetheless, I’m still using this technique in background processes that were
written in C# and accessed those data access DLLs. If you do not need concurrent
calls, it works great. No need to rewrite those old business logic/data access
libraries (yet).

What do you think? Is it worth continuing or the amount of code required is too
considerable, I’m mostly thinking for cases where the C# function accept classes
instead of primitive data types.

* [Programming](https://dominicstpierre.com/tagged/programming?source=post)
* [Golang](https://dominicstpierre.com/tagged/golang?source=post)
* [Csharp](https://dominicstpierre.com/tagged/csharp?source=post)
* [Dotnet](https://dominicstpierre.com/tagged/dotnet?source=post)

