There was an interesting (and popular) [post](https://www.reddit.com/r/golang/comments/18sncxt/go_nil_panic_and_the_billion_dollar_mistake/) on the Go subreddit last week. It told a story of a company losing money due to a null pointer error in production caused by introducing a NULL value into a database table in production.

I've published a [YouTube video](https://www.youtube.com/watch?v=GYaCfpwgbzE) to speculate on what could have been the error, as the post was not very clear on that. This article will propose actionable guide rails you can use to prevent this from happening.

The reason it touched me like that was because the OP mentioned something at the end of the article:

> To me this was a big eye opener. I'm pretty experienced with Go and was previously recommending it to everyone. Now I am not so sure anymore.

Null pointer errors have been plaguing our industry for a long time, and they've nothing to do with Go Perse. I also published a [podcast episode](https://share.transistor.fm/s/2ad776aa) if you want more of my opinion on that subject.

## Handle database NULL values

The `sql` package already solves that exact specific error. Multiple types can be used instead of primitive type pointers for instance:

```go
type YourTable struct {
  Col sql.NullString
}

```

The package offers multiple replacement types to gracefully handle the possibility of having a NULL value in the database.

```go
if !mytable.Col.Valid {
  return ""
}
return mytable.Col.String
```

You use the `Valid` field to determine if the database value is valid. If it is, you can use the `String` field. The `String` field would be an empty string even if the value in DB were NULL. There's no risk of having any panic on your Go program for dereferencing a pointer that's nil.

That's the issue in the article: if you're going to use pointers, you'll need protection to prevent dereferencing null pointers.

## Guard rails when using pointers

Go makes it very easy for programmers to use pointers. Does that mean you need to have everything a pointer, not at all?

On the contrary, you should think twice before using a pointer, and if you do, here are some guard rails you'll need to have to ensure your program runs smoothly.

### 1. Check for nil

Dah, WTH, of course!

I know this sounds basic, but since there are still so many errors in backend programs caused by null pointers, one has to point out the obvious.

For programmers not liking the Go's `if err != nil`, well, using a pointer that you know might potentially be nil will require you to write a lot of code blocks like this:

```go
type Thing struct {
  Value *string
}
var thing *Thing

func (t *Thing) DoSomething() string {
  if t == nil || t.Value == nil {
    return "<nil>"
	}
  return "OK: " + *t.Value
}

thing.DoSomething() // prints: <nil>
thing = &Thing{Value: "no panic"}
thing.DoSomething() // prints: OK: no panic
```

There are a couple of important things to unpack here.

There are two possibilities for a nil pointer reference error, so we have to check both the structure and the field for a nil pointer before using them.

Secondly, do you notice how we must dereference the pointer for the concatenation `"OK: " + *thing.Col``? The LS (language server) often does this automatically, so you need to be careful and check for the possibility of nil before accessing the value pointing by the pointer.

Have you written this everywhere you're using pointers? I certainly have not, if I'm 100% honest. Using pointers can bite you in production multiple months, even years after deployment.

### 2. Use error to communicate no data

Often, we create a function that returns a structure and an error. Generally, the structure is a pointer indicating there's no data.

```go
func GetSomething() (*Thing, error) {
  // some code path that give
  return nil, nil
}
```

The issue happens when there's no error returned, and the structure pointer isn't initialized, so it's nil. In that case, the caller of the function must be disciplined and check for the null pointer; otherwise, well, panic and crash.

```go
import somepkg

func main() {
  thing, err := somepkg.GetSomething()
  if err != nil {
    log.Fatal(err) // I know it's not imported
	}
  print(thing.Col) // BOOM, this panic and crash
}
```

What if we used our good old friend, the error value, to communicate empty data and not use a pointer for the structure?

```go
var ErrNoData = errors.New("no data available")
func DoSomething() (Thing, error) {
  thing := Thing{} // initialize
  // actually do something
  // if a code path find no data
  return thing, ErrNoData
}
```

Now, the caller will receive an error, and they can check if that error is the `ErrNoData` and have a better way to handle that scenario without the risk of crashing the program.

```go
import somepkg
func main() {
  thing, err := somepkg.DoSomething()
  if err != nil {
    if errors.Is(err, somepkg.ErrNoData) {
      // handle this
		}
    // handle all other errors
	}
  // thing is a non-pointer struct ready to be used
}
```

### 3. Code reviews

As stated in the Reddit post, the programmers who wrote the program were new to Go. However, these errors should have been caught in code reviews, even tests.

More experienced programmers are used to dealing with null pointers and should already have a defense mechanism to protect themselves.

I recommend introducing discipline as the #1 core value you want your team to develop. It often takes a lot of work to balance fast shipping and correctness in the code. Discipline will help prevent some of those production issues.