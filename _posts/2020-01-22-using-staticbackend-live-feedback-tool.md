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
