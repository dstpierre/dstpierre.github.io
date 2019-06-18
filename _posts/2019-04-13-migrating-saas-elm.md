---
permalink: "/elm-for-go-developers"
layout: post
title: "As Go developers had you tried Elm yet?"
date: 2019-04-13 09:09:45 UTC
updated: 2019-04-13 09:09:45 UTC
comments: false
summary: "..."
---

I've written small prototypes recently in a multitude of languages to determine 
which would be the best fit for a massive cloud migration I'm participating via 
consulting work.

I'm already sold to Go that's not a secret, having built my last two startup 
products with Go and React and [publishing a book](https://buildsaasappingo.com/) 
in last September. Erlang is following and is very close. I would encourage you 
to look at Erlang in your spare time. It's extremely enjoyable.

The short version of this story is that we've tried .NET Core, Clojure, Erlang, 
Python and Go. And the reasons I'm talking about this is simple.

I've been playing with Elm since its version 0.17, and I finally decided to use 
it as my main frontend language for all the product I'm building. 

The tooling and the compiler of Elm to me, as a Go developer is making me happy. 
And at the point where I'm in my career developer happiness is now #1 on my list. 
I started programming a long time ago because I had fun writing code and it's 
still true today.

### Slowly introducing Elm

You probably noticed If you read article or watched conferences on Elm that the 
recommended way to getting started is by introducing Elm in an existing 
application by replacing a tiny piece of functionality. 

You may probably do that if you're an employee at a company and want to see if 
your team will embrace Elm or not. I'm a consultant (for now), but I'm more an 
entrepreneur. And one unique aspect of building your projects is that you have 
total control over what technology you use.

### Go straight to the point.

Here's the equivalent of introducing Elm slowly into an existing code base for 
those that are starting their project and want to see if Elm is a good fit for 
them.

Start by creating the page that your users will have the most value from using 
your web application. For example, I'm currently migrating and relaunching a 
[legacy SaaS](https://www.getosmosis.com) I've built in 2012 in Elm and Go. The 
original stack was ASP.NET MVC and some jQuery.

Don't worry about SPA, session management, authentication, data access, JSON 
decoder/encoder, etc. Those are not tricky aspects of Elm and will not block 
you if you decide to go with Elm in the long run. Well maybe the JSON 
decoder/encoder, but we will see that in a later post I promise.

Just create one Main.elm file implementing the flow of the 
[Elm architecture](https://guide.elm-lang.org/architecture/). Don't worry about 
anything else, but make sure you build the most useful piece that users of your 
app will get the most value, it's usually the most complex page anyway.

The product I'm relaunching is [Osmosis](https://www.getosmosis.com). It's a 
proposal and sales documents management web app. The page I used for the Elm 
prototype was the document editor, or at least a good enough prototype to let me 
feel confident that I wanted to continue using Elm for the rest of the application.

### Start small

Your first task is to read the [Elm getting started guide](https://guide.elm-lang.org). 
Then you may proceed. Here's what my Main.elm file looked like at first.

```elm
module Main exposing (main)

import Browser
import Html exposing (Attribute, Html, a, div, h1, h4, i, input, label, li, ol, option, p, select, span, text, textarea, ul)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, onInput)



-- MAIN


main =
    Browser.document
        { init = init
        , subscriptions = subscriptions
        , update = update
        , view = view
        }



-- MODEL


type alias Page =
    { name : String
    , order : Int
    }


type alias Model =
    { name : String
    , pages : List Page
    }



-- INIT


init : () -> ( Model, Cmd Msg )
init flags =
    ( Model "test" [], Cmd.none )



-- UPDATE


type Msg
    = Todo


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Todo ->
            ( model, Cmd.none )



-- VIEW


view : Model -> Browser.Document Msg
view model =
    { title = "Page title"
    , body =
        [ h1 [] [ text "todo" ]
        ]
    }


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none


```

As oppose to Go, in Elm I like to start with the view and climb to the update 
function one step at a time. For instance, I needed a sidebar and a fixed 
top-right toolbar and the main content, this is what I prototyped.

Let's implement the add page button click. This should be enough of a demo to 
get you started. I'll go into common issues I faced during my first two weeks 
of using Elm full-time at the end of this post.

{% include push-content.html %}

```elm
viewToolbar : Model -> Html Msg
viewToolbar model =
    div [ id "toolbar" ]
        [ div [ class "buttons has-background-dark" ]
            [ viewToolButton "plus" NewSection
            , viewToolButton "plus" NoOp
            , viewToolButton "check" NoOp
            , viewToolButton "times" NoOp
            ]
        ]


viewToolButton : String -> Msg -> Html Msg
viewToolButton icon msg =
    a [ class "button is-black", onClick msg ]
        [ span [ class "icon" ] [ i [ class ("fas fa-" ++ icon) ] [] ]
        ]
```

I've created a function to display a toolbar button and I'm passing the desired 
`Msg` for the `onClick` event.

I like to isolate reusable piece of the UI inside smaller function. The NoOp `Msg` 
is used as placeholder for now, I find it useful to focus on one flow at a time. 
Since Elm require you implement all possible path, having a `NoOp` is helpful 
when prototyping.

From there the next step is to add the `Msg` named `NewSection`:

```elm
type Msg
    = NoOp
    | NewSection
```

From here I like to update my `Model` before implementing the behavior in the 
`update` function for this new `Msg`.

We will need a `List Section` in our model to append a new item to the list. 
`Section` is defined elsewhere but lets define it inside this module for this 
demo prototype.

```elm
type aias Section =
    { title : String
    , order : Int
    }

type alias Model =
    { sections : List Section
    }
```

And now we can implement the case for `NewSection` in the `update` function:

```elm
update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        NoOp ->
            ( model, Cmd.none )

        NewSection ->
            let
                ns =
                    Section "no title" 0
            in
            ( { model | sections = List.append model.sections [ ns ] }, Cmd.none )
```

To test this interaction we could quickly display the `sections` list inside our 
main view:

```elm
view : Model -> Browser.Document Msg
view model =
    { title = "Page title"
    , body =
        [ viewToolbar model
        , List.map viewSection model.sections
            |> ui []
        ]
    }

viewSection : Section -> Html Msg
viewSection s =
    li [] [ text s.title ]
```

It displays the each sections inside an ul/li HTML elements.

### Some learning I thought might help when starting real project in Elm

**Difference between Msg and msg**: This one is actually pretty simple but not 
extremely clear, at least for me, at first. The capital `Msg` means it's a type 
where the lower-case `msg` represent a generic unbound type.

For example, let's say you would want a function to output a button. You might 
not want to use the `button` directly and would prefer to have `myButton` instead.

You could define this in a module where you could have other common elements 
wrapped in helper functions:

```elm
module Element exposing(myButton)

import Html exposing(Html, button)
import Html.Attributes exposing(class)
import Html.Events exposing(onClick)

myButton : String -> msg -> Html msg
myButton caption action =
    button [ class "button", onClick action ] [ text caption ]
```

And you would use it like this in other module:

```elm
module PageX exposing(Model, Msg, init, view, update)

import Element exposing(myButton)

type Msg
    = GotClicked

...

view : Model -> Html Msg
view model =
    myButton "Click me" GotClicked
```

The declaration of the `myButton` function in the module `Element` accept a 
generic type. The `msg` could be renamed to anything else like so:

```elm
myButton : String -> a -> Html a
```

**Use Maybe not default value**: As Go developer we're used to deal with default 
value and only `nil` when we have pointer. It pays in the long terms to use `Maybe` 
in Elm. The compiler will help making sure all possible code path are handled.

```elm
type alias Model =
    { id : Maybe Int
    }

init flags =
    ( Model Nothing, Cmd.none )
```

instead of:

```elm
type alias Model
    { id : Int
    }

init flags =
    ( Model 0, Cmd.none )
```

If you know that `id` might not have a value, using `Maybe` is recommended over 
trying to replicate Go's default value.

**What's the <| operator**: This prevent having to use the () to evaluate 
functions:

```elm
doSomething (String.fromInt 3)
```

versus

```elm
doSomething <| String.fromInt 3
```

I personally find the first with () way clearer. Maybe this has to do with the 
fact that I have a ~18 years of imperative programming behind me, YMMV.

**Don't worry about not having it "the correct way"**: It is actually very easy 
to refactor Elm code. It took me two more weeks to publish this article, so at 
the start of this writting I had 2 weeks of Elm full-time, but now I'm at 4 weeks.

And of course I've made lots of mistakes and I've not used optimal way of doing X 
or writting Y. I've learned new things along the way and I changed what I thought 
could be better.

I would also that refactoring Go code is pretty easy, the majority of compiled 
languages in fact. So don't be blocked by the fact that it is different. Maybe 
you've written your share of JavaScript so far, I know I did. It's worth the 
investment.