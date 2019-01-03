---
layout: post
title: "How to start with TypeScript and Preact"
date: 2017-08-24 11:26:37 UTC
updated: 2017-08-24 11:26:37 UTC
comments: false
summary: "..."
--- 

![](https://cdn-images-1.medium.com/max/600/1*IXH_QkfEOBbzQE9R6OnU3g.png)

The last [two](https://www.leadfuze.com/) [startups](https://roadmap.space/)
I’ve built used [TypeScript](https://www.typescriptlang.org/) and
[React](https://facebook.github.io/react/) for the front-end.

As we can see in this screenshot, Roadmap’s bundle is 1,327,161 bytes, which is
1.327161 megabytes if I remember my elementary math.

I know, I know, this is before gzip compression. Still a heck of a first payload
before it’s in the browser’s cache.

And this is only the app’s own bundle, there’s lots of other resources to load
on top of that like [Intercom](https://intercom.io/), Google Analytics, the CSS,
all the images, etc. I’d like my SPA to be way smaller than that, since React
account for at least 90% of the size of the bundle.js, there’s not much I could
do.

A couple of months ago I’ve started to hear good things about
[Preact](https://preactjs.com/) and their tag line (or their home page HTML
title) at the time of writing this is:

> Preact: Fast 3kb React alternative with the same ES6 API. Components & Virtual
> DOM.

I’m usually not a fan of taking others’ products and building a case based on
one aspect the other thing do better, but I admit that React is HUGE and 3kb for
advertised ~same functionalities was intriguing.

It’s pretty quick to test Preact with TypeScript, as in the 5 minutes quick,
we’ll do just that.

### Creating a TypeScript project that uses Preact

First step is to create a new project with `npm init` and install the required
packages:

```shell
$> mkdir testing-preact && cd testing-preact
$> npm init
$> npm install typescript preact webpack ts-loader --save
```
We will now create our `tsconfig.json` configuration file to have the TypeScript
compiler do the transpilation:

```json
{
  "compilerOptions": {
    "outDir": "./public",
    "target": "es5",
    "module": "commonjs",
    "noImplicitAny": false,
    "removeComments": true,
    "sourceMap": true,
    "jsx": "react",
    "jsxFactory": "h"
  },
  "include": [
    "./src/**/*.tsx",
    "./src/**/*.ts"
  ]
}
```

One important option here is the **“jsxFactory”: “h”**. Preact use its own class
to render the JSX code to JavaScript, and it’s called “h”. We will see later how
it’s done.

Now lets create our webpack configuration file `webpack.config.js`so we can
build our app bundle:

```js
module.exports = {
  entry: "./src/bootstrap.tsx",
  output: {
    filename: "./public/bundle.js",
  },
  // Enable sourcemaps for debugging webpack's output.
  devtool: "eval",
  resolve: {
    // Add '.ts' and '.tsx' as resolvable extensions.
    extensions: [".webpack.js", ".web.js", ".ts", ".tsx", ".js"]
  },
  module: {
    loaders: [
      // Handle .ts and .tsx file via ts-loader.
      { test: /\.tsx?$/, loader: "ts-loader" }
    ],
  },
};
```

This will load our bootstrap file called `bootstrap.tsx` inside our `src`
director. Now lets create the structure for our test project:

```shell
$> mkdir src
$> mkdir public
```

The project structureshould look like this:

```shell
- public/
- src/
- index.html
- package.json
- tsconfig.json
- webpack.config.js
```

We will now create our app entry point `bootstrap.tsx` inside our `src`
directory:

```js
import { h, render } from "preact";
import { App } from "./app";

render(<App name="cool working" />, document.getElementById("app"));
```
We’re basically simply loading our main app (which does not exists yet) and we
tell Preact to start loading it into our app container with the id of app.

We’re using `render` instead of ReactDOM to accomplish this.

In **all** your Preact components you need to import `h` from `preact` so it can
render JSX to JavaScript via the jsxFactory discussed above inside the
`tsconfig.json` file.

Lets create a quick app to test that what we have so far even work. This is
`app.tsx` inside the `src` directory:

```js
import { h, Component } from "preact";

export interface AppProps {
  name: string;
}

interface AppState {
  name: string;
}

export class App extends Component<
, AppState> {
  constructor(props: AppProps) {
    super(props);

    this.state = { name: props.name };
  }
  componentDidMount() {
    setTimeout(() => { 
      var state = this.state;
      state.name = "Preact's componentDidMount worked as expected";
      this.setState(state);
    }, 2000);
  }
  render(props: AppProps, state: AppState) {
    return <h1>props: {props.name} state: {state.name}</h1>;
  }
}
```

Alright, pretty similar to a React component right? We’re basically simply
declaring the `props` and `state` interfaces for the component and the component
implementation itself.

To test a basic example of the life cycle of component, I’m simply having a
timeout that will set the state of the component after 2 seconds.

One interesting aspect to notice is that Preact does pass the `props` and
`state` as parameter to the `render` function which I find particularly pleasant
to have.

The only missing pieces to test our simple app is our `index.html` in the root
directory:

```html
<html>
<head>
  <title>testing Preact</title>
</head>
<body>
  <div id="app"></div>
  <script>
  // loading the bundle and making sure we're not using cache
  var rnd = (new Date()).getTime();
  var b = document.createElement("script");
  b.src = "/public/bundle.js?v=" + rnd;
  b.type = "text/javascript";
  document.body.appendChild(b);
  </script>
</body>
</html>
```

Now we can build our simple app with:

```shell
$> webpack
```

This should build your app along with Preact lightweight library and the end
result would be in the ~45 000 bytes. A **major** boost I must admit.

### How can they achieve this size you say?

Of course Preact is not supporting 100% of what React have implemented. Here is
what’s not included and my experiences in not missing them and why.

**PropType validation**: With TypeScript we have interfaces to give us
validation and similar effect.

**React.Children**: I typically never use this on my two last production React
apps.

**Synthetic events**: Old Internet Explorer should be updated anyway, sorry :).

I suggest you have a look [here to see the difference from Preact and
React](https://github.com/developit/preact/wiki/Differences-to-React).

### Now as I’m sure everyone is waiting for my opinion ;)

From a back-end Go developer that needed to built 2 SaaS using React, here’s my
thoughts on Preact.

**Pros**

* Lightweight, this means faster download, but faster build when developing as
well which is always nice.
* All accumulated knowledge having worked with React for the last 3 years can be
applied. This is a **major benefits** with the state of JavaScript UI libraries
in the last decade.
* Do have TypeScript type definition baked in, no need to install a separate type
def.
* Work great with TypeScript, and TypeScript do bring some order inside larger
JavaScript project.
* We can use `class` instead of `className` I cannot say how many times I still
continue to type class in JSX, and worst I’ve started to type className when
writing normal HTML which is not good muscle memory to have.

**Cons**

* Not much coverage for CSS UI library like Sementic-UI, Bootstrap, etc. Those
tends to have pretty good React coverage, but Preact is still fairly new.
* Pretty hard to migrate an existing TypeScript React code base if you typed the
synthetic event arguments, for example of a input changed handler
`handleChange(e: React.FormEvent<HtmlInputElement>)` this make the migration
tedious because those types do not exists in Preact.

I’m currently testing [Bulma](http://bulma.io/) as a CSS UI framework to pair
with Preact. The fact that Bulma does not include any JavaScript make it a good
candidate.

For migration, Preact have this [switching to Preact
guide](https://preactjs.com/guide/switching-to-preact) that could help depending
on your current code base.

Even though I’m feeling good when writing Go code (back end), I’m feeling good
when I know that the front-end is as lightweight and as easy as possible for new
devs to get started. Preact is a marvelous option, check it out.

* [Typescript](https://dominicstpierre.com/tagged/typescript?source=post)
* [Preact](https://dominicstpierre.com/tagged/preact?source=post)
* [JavaScript](https://dominicstpierre.com/tagged/javascript?source=post)

