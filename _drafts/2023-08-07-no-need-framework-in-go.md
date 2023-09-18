---
permalink: "/no-need-web-framework-in-go"
layout: post
title: "No need for web framework in Go"
date: 2023-08-07 07:55:40 UTC
updated: 2023-08-07 07:55:40 UTC
comments: false
summary: "..."
---

I used frameworks and libraries when I started building web applications in Go. Coming from other stacks where there's not such a fully-features HTTP server package as in Go in the core library, it's normal to want to use 3rd party packages.

After over seven years, I appreciated writing my things and removing the unnecessary dependencies on 3rd party packages.

Before we begin, note that it's okay if you prefer to use frameworks or libraries, I want to present another option that's not using external dependencies, which for long-living projects can be a lifesaver in ten years.

### Routing

Let's start with the elephant in the room, the routing aspect of web API and applications.

In all discussions I've been part of, it's usually one of the top reasons one uses a library like [chi](https://github.com/go-chi/chi) or [httprouter](https://github.com/julienschmidt/httprouter).

Of course, the [net/http](https://pkg.go.dev/net/http) package isn't as good  regarding URL segmentation of parameters. Here's what I've been using for years, and I don't see why I'd need more.

```go
func GetURLParam(r *http.Request, at int) string {
  arts := strings.Split(r.URL.Path, "/")
  if len(parts) <= at {
    return ""
  }
  return parts[at]
}
```

Take the following URL: `/messages/1234/reply/4`.

If we would want to grab the presumably message id and reply id we'd use this in our handler:

```go
func handler(w http.ResponseWriter, r *http.Request) {
  msgID := GetURLParam(r, 2)
  userID := GetURLParam(r, 4)
}
```

That's really all there is to get those. Now granted, they're `string` and maybe you'd want them to be `int`.

```go
func GetURLParam(r *http.Request, at int, v any) {
  parts := strings.Split(r.URL.Path, "/")
  if len(parts) <= at {
    return
  }

  switch v.(type) {
  case *int:
    i, err := strconv.Atoi(parts[at])
    if err != nil {
      return
    }

    if val, ok := v.(*int); ok {
      *val = i
    }
  case *string:
    if val, ok := v.(*string); ok {
      *val = parts[at]
    }
  }
}
```

And you may use this as:

```go
var msgID int

GetURLParam(r, 2, &msgID)
```

By splitting the URL path via the slash `/` character, we can grab the desired segment and convert it to an integer or leave it as a string. From the caller's point of view, it's still just a one-liner.

The next issue with the standard library regarding routing is the way URLs might look limited. You may fix this one by changing how you structure your URLs. Let's see.

Instead of having this route:

```
/messages/1234/reply/4
```

You could have:

```
/messages/reply/1234/4
```

It's not RESTful, or it's not standard. Arguably, URLs with parameters are not required to be human-friendly. Having all your parameters at the end of the URLs is easier to use the stdlib package.

```go
http.HandleFunc("/messages/reply/", msg.reply)
```

The last slash `/` at the end of this URL enables you to have other URL segments retrieved via the GetURLParam above. What about URLs with a lot of parameters? Well, there's such a thing as too many parameters as well.

### Verb based handler

Let's take an example with the `chi` router:

```go
r := chi.NewRouter()
r.Get("/", func(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("GET request"))
})
r.Post("/", func(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("POST request"))
})
```

Here's the equivalent using the standard library:

```go
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodGet {
    w.Write([]byte("GET request"))
  } else if r.Method == http.MethodPost {
    w.Write([]byte("POST request"))
  }
})
```

Again, less comfortable, granted. But is it really that less? There will often be code repeated for different methods like POST, PUT, and PATCH that might share the same parsing logic for the request body.

Suppose you're keeping your handler light on your domain logic and keep them mostly doing the parsing of parameters, calling another package that handles those parameters and returns the response. Having one handler and a conditional block to check what's the HTTP method is pretty convenient.

Here's an example of a potential handler for a `Message` type:

```go
type Message struct {
  ID string `json:"id"`
  Body string `json:"body"`
}
func message(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodGet {
    list, err := message.List()
    respond(w, http.StatusOK, list)
  } else if r.Method == http.MethodDelete {
    http.Error(w, "not supported", http.StatusMethodNotAllowed)
    return
	} else {
    var msg Message
    if err := parseBody(w.Body, &message); err != nil {
      //...
    }
    msg.Save(msg)		
    respond(w, http.StatusOK, msg)
  }
}
```

