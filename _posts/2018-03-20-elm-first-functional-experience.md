---
permalink: "/thanks-elm-i-had-my-first-taste-at-functional-programming-your-turn-now-893f2bf8f4be"
layout: post
title: "Thanks Elm, I had my first taste at functional programming, your turn now"
date: 2018-03-20 11:26:37 UTC
updated: 2018-03-20 11:26:37 UTC
comments: false
summary: "..."
---

After 18 years of programming with imperative languages, I finally had a chance
at completing a tiny tool in Elm giving me my first experience using a
functional programming language.

TL:DR: It was harder than what I expected, but as I wrote and read more and more
code, it started to feel less awkward. Also, the fine folks in #beginner [Slack
channel](https://elmlang.herokuapp.com/) were **very** supportive.

Here’s the code.

#### What were my motivations to use Elm?

I enjoy learning new things. I remember how I was excited when I built my SaaS
in Go 4 years ago.

On the UI I’ve been building my two [last](https://www.leadfuze.com/)
[startups](https://roadmap.space/) with React. I find React with TypeScript to
be a very comfortable combination to build and maintain client-side
applications.

After some time using the same stack though I start to feel a void, and learning
new things give me energy. Elm seems a pretty good language to try since I’m
trying to get started in functional programming for a long time now without ever
succeeding.

I like using the right tool for the job. When looking at the selling points of
Elm, I immediately compared with the situations I see at
[Roadmap](https://roadmap.space/).

**No runtime exceptions**. I would say that currently 90% of all front-end
errors captured by Bugsnag are related to X is null or undefined and similar
that TypeScript did not caught. To me, a language that says it can protect my
code from these types of errors would represent a huge improvement.

*Note that I’m not going to start rewriting Roadmap’s front-end code with Elm
(not yet at least). But starting to get familiar with this language could mean
that the app I’ll build with it should get some nice production stability.*

**Easier build and update process**. It might just be me, but most often than
not when I do either a TypeScript, React or Webpack update my build process is
broken for at minimum a couple of hours while I try to fix everything.

Elm has this nice tool where you simply run a command, and it builds your app
like so:

    elm-make src/file.elm

During development, I used `elm-reactor`and frankly that was the quickest way to
have an auto-reload development server. I’m not saying that these things are
unique to Elm, but from my experience so far with the front-end apps I worked
on, their tooling is simply right on. It feels a bit similar to Go, and I really
enjoy Go still after all these years.

The tooling include:

`elm-make`to build your app, either a standalone HTML file or a bundle
JavaScript file.

`elm-package`to manage packages.

`elm-reactor`as development server handling reloading.

`elm-repl`an Elm REPL. It is not working for me on Arch, I have a missing
dependency on `libtinfo.so.5`but I did not see any use for the REPL at this
moment. Maybe this will come later where I’ll need to focus on that issue.

#### How to get started with Elm

It is quick to get started. You basically [install
Elm](https://guide.elm-lang.org/install.html) tools via npm and you’re good to
go. I would highly suggest reading their [guide](https://guide.elm-lang.org/)
and writing the tutorial example code yourself.

Configure your code editor for Elm and dive into a brand new world, well for me
at least, functional programming is something I’m targeting for a long time.
Haskell is a language I want to use one day.

#### Everything was fine until I started my tiny tool.

No more guided tour, no more hand holding. You’re alone my friend. And boy ho
boy it **has been** difficult.

I started from this:

{% include push-content.html %}

```elm
import Html exposing (..)

main : 
	
	
	
main =
		Html.program
				{ init = init
				, view = view
				, update = update
				, subscriptions = subscriptions
		}

type alias Model =
		{ name : String
		}

type Msg
		= InputName String

update : 
	-> 
	-> ( 
, 
	
	)
update msg model =
		case msg of
				InputName n ->
						( { model | name = n }, Cmd.none )

view : 
	-> 
	
view model =
		div []
				[ input [ onClick InputName, value model.name ] []
				]

subscriptions : 
	-> 
	
subscriptions model =
		Sub.none

init : ( 
, 
	
	)
init =
		( { name= ""
			}
			Cmd.none
		)
```

I will not go into detail about the [Elm
architecture](https://guide.elm-lang.org/architecture/). Their getting started
guide is there for that and really useful. I just want to describe my experience
as I think other developers might feel the same way, which is; Ok, what I do
from this starting point?

My goal was to build a simple [customer research
form](https://roadmap.space/call/). That small project was a perfect one I think
since it was so small but involved using HTTP calls, JSON encoding and decoding
and multiple views.

*I think one mistake would be to start with a project that’s way to big. Because
at first you might even ask how to organize your Elm code. Having a one file
simple project for your first try will increase your chances of having a fun
experience.*

Then you start writing your first piece of Elm code from yourself. I was shocked
by how a simple task was harder to implement in Elm. For example, thinking that
a function should do one thing only is not immediately what my 18 y/o brain of
programming is used to.

Also it almost feels like there’s `[]`and `()`everywhere. And even if you’ve
read and understood that this `(do "X")`is a function call, for some reason, to
me at least, that did not came naturally at all (even now after this project
lol).

> One trick I have for new comers. If the compiler give you an error on a function
> call, you might be missing `() `around the call.

And surprisingly enough things are not bad after 1 hour, until the compiler
returns an error you cannot really figure out as a beginner Elm developer.

The Elm compiler’s error messages are extremely useful, most of the time even
when starting you should be able to decently understand the issue. But the other
15% of the time when starting some errors are not very clear.

To me, the majority of weird issues I was having was due to missing parenthesis.
Meaning that the value I was passing to a function needed to be evaluated before
calling the function. This is a good example of something I was not
understanding after 2–3 hours of writing Elm code.

```elm
div [] [List.map learningView l]
vs
div [] (List.map learningView l)
```

This will run the `learningView`function for each element of the list `l`. The
correct way of expressing is the second one. To me, that was not very clear at
all. So maybe you’ll have those kind of errors. Here’s the compiler error
message:

I selected a project that I had no deadline as well. I wanted to build this
during March, but that was it. I did 2–3 hours grouped work sessions working on
this.

I learned so many concepts that it was difficult to stay focused more than that.
Note that I *might be* getting old that could explain YMMV.

#### Things I dig about Elm.

During the time I was building the app, I did not focus at all on the UI, no CSS
classes, no styling. This is not unique to Elm, but for some reason I found
myself focussing more on the code than the UI. Maybe it had to do with the fact
that I was intensively learning, we will see the second project.

Once the app was all functioning I added the classes and added more elements
that had no impact on the functionality.

Again, elm-reactor is just dead simple and pretty powerful. The development flow
was streamlined.

I felt really good once I wrote my first functions that had nothing to do with
the Elm architecture process of model, update, view. My only advice is to go
slowly and don’t get discouraged.


