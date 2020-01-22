---
permalink: "/using-staticbackend-live-feedback-tool"
layout: post
title: "Showcasing StaticBackend while building a feedback tool"
date: 2020-01-22 07:55:40 UTC
updated: 2020-01-22 07:55:40 UTC
comments: false
summary: "..."
---

Last week I 
[posted on Twitter](https://twitter.com/dominicstpierre/status/1216673403450601472) 
saying that I was going to build and give a SaaS to get my first paying 
customer at StaticBackend.

That turned out to be hard to execute as a claim. Here's 
[more explanation](https://www.youtube.com/watch?v=aMMXtzuvgew&).

I took five days to think about this and decided that the best way to showcase 
the usefulness of StaticBackend was to build open-source tools.

I'm trying to create a win-win situation where the tool would be bound by 
default to StaticBackend, hence requiring a paid account to use the tool as-is. 
But nothing stops someone from writing an open-source backend that would be a 
replacement for my backend as a service.

I've decided to build a feedback and changelog embeddable widget.

![ClearUser mock-up](/assets/img/clearuser-proto.png)

I've started coding at around 7 am. I started a new Elm frontend project.

```bash
elm init
```

I'm going to post a live update here and on Twitter during the day to showcase 
the evolution of the project.

This is the project structure so far.

```
├── clearuser.js
├── demo-host-app.html
├── elm.json
├── index.html
├── main.js
├── src
│   ├── Api
│   │   ├── DB.elm
│   │   ├── Endpoint.elm
│   │   └── Membership.elm
│   ├── HttpRequest.elm
│   ├── Main.elm
│   └── User.elm
└── start.sh

4 directories, 20 files
```

I'm using an HTML page to emulate how a real web application would embed the 
project in their application.

```html
<html>
<head>
	<title>Demo host app</title>
</head>
<body>
	<h1>Demo host app</h1>
	<p>This emulate how to embed the ClearUser app into your own application.</p>

	<button onclick="clearuser.show()">Click here to open ClearUser</button>

	<script src="/clearuser.js"></script>
	<script>
		clearuser.init("user@domain.com", "some-uniq-id", "https://yourapp.com/their/avatar.png")		;
	</script>
</body>
</html>
```

It's using an `iframe` to prevent from having CSS and layout issues. The Elm 
application communicates with the backend using HTTP requests.

I've created some helpers module here to use StaticBackend URL endpoints.

```elm
module Api.DB exposing (create, delete, fetch, list, save)

import Api.Endpoint exposing (Endpoint, url)


create col =
    url [ "add", col ] []


list col =
    url [ "list", col ] []


fetch col id =
    url [ "get", col, id ] []


save col id =
    url [ "update", col, id ] []


delete col id =
    url [ "delete", col, id ] []
```

I wanted to be in a working state, project is not compiling because I just 
introduce the `HttpRequest` module:

### 8:20: Tweet and blog post

I just posted this blog post and tweeted about my plan for the day.

But now taking a small break for breakfast and talk with my daughters and wife.

When I'm back, I'll attack the user management (register and login) and the 
initial views.


### 10:07: User login/register completed

I've used a `User` module in Elm and added two functions for `login` and 
`register`:

```elm
login : ( String, String ) -> (Result Http.Error String -> msg) -> Cmd msg
login ( email, pin ) msg =
    post
        ""
        Endpoints.login
        (Http.jsonBody
            (Encode.object
                [ ( "email", Encode.string email )
                , ( "password", Encode.string pin )
                ]
            )
        )
        (Http.expectJson msg Decode.string)


register : ( String, String ) -> (Result Http.Error String -> msg) -> Cmd msg
register ( email, pin ) msg =
    post
        ""
        Endpoints.register
        (Http.jsonBody
            (Encode.object
                [ ( "email", Encode.string email )
                , ( "password", Encode.string pin )
                ]
            )
        )
        (Http.expectJson msg Decode.string)
```

The are HTTP requests sending the following JSON:

```json
{
	"email": "current_user@email.com",
	"password": "current_user_pin_from_host"
}
```

They both receive a `string` which is the authentication token. This `token` 
will be used for all authenticated calls to the API:

```elm
credHeader : String -> List Http.Header
credHeader tok =
    [ Http.header "Authorization" ("Bearer " ++ tok)
    , Http.header "SB-PUBLIC-KEY" "5e285f6bfe98b7b19450baad"
    ]
```

The host app will pass the necessary information and Elm will initialized 
the app receiving the info as `flags`.

```elm
init : Decode.Value -> ( Model, Cmd Msg )
init flags =
    case Decode.decodeValue U.decoder flags of
        Ok usr ->
            let
                u =
                    { usr | isAdmin = False }
            in
            ( { user = u
              , token = ""
              }
            , U.login ( u.email, u.pin ) GotLogin
            )

        Err _ ->
            ( { user = U.User "" "" "" False
              , token = ""
              }
            , Cmd.none
            )
```

We're calling a `login` when the Elm application initialize.

From there we can have our starting point `update` function:

```elm
type Msg
    = GotLogin (Result Http.Error String)
    | GotRegister (Result Http.Error String)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        GotLogin res ->
            case res of
                Ok tok ->
                    ( { model | token = tok }, Cmd.none )

                Err _ ->
                    ( model
                    , U.register ( model.user.email, model.user.pin ) GotRegister
                    )

        GotRegister res ->
            case res of
                Ok tok ->
                    ( { model | token = tok }, Cmd.none )

                Err err ->
                    ( model, Cmd.none )
```


I'm now going out doing some skating with the kids, will be back after lunch.

### 13:52: After a long skating/lunch break, attacking the Feedback

Side note, this year I decided to lower my consulting days to three per week to 
have time to work on my product(s) and be with my two daughters as much as 
possible before they go to school next September. We are homeschooling them 
for the last 10 years.

![went skating to take some fresh air](/assets/img/clearuser-live-skating.jpeg)

I've started the `Feedback` layout. I'm usually writing the `model` first when 
using Elm. This is what I have so far for the `Feedback` model:

```elm
type Filter
    = Trending
    | Top
    | New

type Status
    = All
    | New
    | Pending
    | UnderReview
    | Planned
    | Completed
    | Rejected

type alias Feedback =
    { id : String
    , accountId : String
    , title : String
    , desc : String
    , likes : Int
    , owner : User
    , users : List User
    , comments : List C.Comment
    , pinnedReply : Maybe C.Comment
    , status : Status
    , postedAt : Iso8601.decoder
    }
```

Currently writing the JSON encoders/decoders for all those model, here's a 
piece of code:

```elm
decoder : Decoder Feedback
decoder =
    Decode.succeed Feedback
        |> required "id" Decode.string
        |> required "accountId" Decode.string
        |> required "title" Decode.string
        |> required "desc" Decode.string
        |> required "likes" Decode.int
        |> required "owner" U.decoder
        |> required "users" (Decode.list U.decoder)
        |> required "comments" (Decode.list C.decoder)
        |> optional "pinnedReply" C.decoder Nothin
        |> required "status" statusDecoder
        |> required "postedAt" Iso8601.decoder
```

I'm starting to think that I will change the API of StaticBackend's database 
a little bit.

That was also the goal of creating this little real-world project, to make sure 
StaticBackend can handle lots of scenarios and stay flexible yet simple.

I need to have a way to have the feedback viewable by multiple users. The way 
I designed StaticBackend was that a users could only view their own data.

To handle this scenario I'm introducing collection prefixed with 
`pub_` which will required to have a valid authentication token, but will not 
limit to only the current user's data. In the case where mutliple users need to 
view the data, this is useful.

Next up is starting to have some view rendered, for instance, my `Main` module 
have this as `model`:

```elm
type Layout
    = Loading
    | Failed String
    | Feedback
    | Changelog
    | Roadmap
    | NewPost


type alias Model =
    { user : U.User
    , token : String
    , layout : Layout
    }
```

I'll be using that `Layout` type to control what is displayed in the main 
view. More on that on my next update.

### 17:52: What a crazy day

It has been a long time that I did not had a productive day like today.

I did not reached my goal of having the first view rendered. But I'm extremely 
satisfied of where the project is at. I was kind of stock for the Feedback 
query.

To let users select if they want to sort by Trending, Top or New and filter by 
status like Under Review, Planned.

Here's the piece of code without the type definitions:

```elm
list : String -> Filter -> Status -> (Result Http.Error (List Feedback) -> msg) -> Cmd msg
list tok filter status msg =
    let
        sortBy =
            case filter of
                Trending ->
                    "trendingLikes"

                Top ->
                    "likes"

                Latest ->
                    ""

        body =
            if status /= All then
                Encode.list queryEncode
                    [ ValueString "status", ValueString "=", ValueString (statusToString status) ]

            else
                Encode.list queryEncode
                    [ ValueString "status", ValueString "!=", ValueString (statusToString New) ]
    in
    post
        tok
        (Endpoints.query "pub_feedback" sortBy)
        (Http.jsonBody body)
        (Http.expectJson msg (Decode.list decoder))
```

Project is compiling for now. I can't wait to continue on this tomorrow morning.