---
permalink: "/using-typescript-to-publish-a-testable-react-npm-package-f3811b3c64e3"
layout: post
title: "Using TypeScript to publish a testable React npm package"
date: 2017-04-01 11:26:37 UTC
updated: 2017-04-01 11:26:37 UTC
comments: false
summary: "..."
---

## I’ve recently migrated a [React](https://facebook.github.io/react/) component
package named [react-trix](https://github.com/dstpierre/react-trix) to
[TypeScript](https://www.typescriptlang.org/) and wanted to talk about the
experience of creating anNPM package in TypeScript.

I’m not here to convinced you to start using TypeScript,
[Flow](https://flow.org/) or [insert name of latest JavaScript type checker]
compiler. [Jason Dreyzehner](https://medium.com/u/8e9fa353daf3) wrote an
excellent article if you want to know more about [giving TypeScript a chance
here](https://medium.freecodecamp.com/its-time-to-give-typescript-another-chance-2caaf7fabe61).

I’ve been using TypeScript for ~3 years now but I never published a public NPM
package with it and I could not find much in terms of recent article other than
[this
one](https://medium.com/@mweststrate/how-to-create-strongly-typed-npm-modules-1e1bda23a7f4)
from [Michel Weststrate](https://medium.com/@mweststrate) for TypeScript < 2.0.
So here we are.

You can clone this repository if you want to follow from your code editor. This
is a quick starting point project where you can simply edit the `package.json`
file, run `npm install`and create your TypeScript React component(s).

### The project structure

I normally create a project directorywith the following directories/files
structure for NPM packages.


```shell    .
├── lib
├── LICENSE
├── main.js
├── package.json
├── README.md
├── src
├── tests
├── test.sh
├── tsconfig.json
```

Let’s install TypeScript:

`npm install typescript -g`

You might need to use `sudo`if you’re not using
[NVM](https://github.com/creationix/nvm), but I would suggest you start doing it
if you don’t ;), this is out of scope of this article though, but NVM make it
easier to have side-by-side Node versions installed in your user space.

At minimum you’ll need React and ReactDOM packages with their appropriate
typings packages:

`npm install react react-dom @types/react @types/react-dom --save-dev`

The `@types/react` and `@types/react-dom` are useful to help the TypeScript
compiler understand external libraries and also let your code editor offer you a
better experience. I’m using [vscode](https://code.visualstudio.com/) and the
TypeScript and Go experiences for writing code is what I enjoy at the moment.

### The TypeScript configuration file

It’s easier to have a TypeScript configuration file when calling the TypeScript
compiler `tsc`rather than having to supply the arguments via the CLI. The file
is named `tsconfig.json` and will be loaded automatically.

```json
{
	"compilerOptions": {
		"outDir": "./lib",
		"target": "es5",
		"module": "commonjs",
		"noImplicitAny": true,
		"removeComments": true,
		"declaration": true,
		"sourceMap": true,
		"jsx": "react"
	},
	"include": [
		"./src/**/*.tsx",
		"./src/**/*.ts"
	],
	"exclude": [
		"node_modules"
	]
}
```

To stay compatible with the “world” our TypeScript source files will need to be
compiled to JavaScript so when other developers install and use our package they
will in fact import a standard [EcmaScript
5](https://en.wikipedia.org/wiki/ECMAScript) source code.

The `target` is set to `es5` meaning that we want the compiler to produce code
compatible with this standard and we’re also saying that our `module` will be
the `commonjs` same as Node.

There’s two important properties in here, the `declaration` and the `jsx`.

The `jsx` option tells the compiler we’re going to have React’s JSX in our
TypeScript source.

The `declaration` indicates that we want to generate declaration files so other
TypeScript developers will have a nicer experience when they install our
package.

Remember those `@types/react` packages we installed previously for React? Our
package will include its own TypeScript definition files `d.ts`. They accomplish
the same goals as the `@types/react` for example, but are part of our NPM
package.

### A first TypeScript React component

It’s time to write a first React component in TypeScript, we’ll write a simple
component for now.

```js
import * as React from "react";

export class Hello extends React.Component<{}, {}> {
	constructor() {
		super();
	}
	render() {
		return <h1>Hello world</h1>;
	}
}
```

Alright, time to test our component. And no we will not use another project, at
least not for now. It’s not time to test our package, it’s time to test our
React component and we will also use **TypeScript** to do that.

We’ll start by installing the needed packages:

```
react-addons-test-utils mocha --save-dev
```

Testing React component is out of scope of this article. But this is part of
publishing an NPM package I think. So I’m not going to talk much about why I’m
using [Mocha](https://mochajs.org/) instead of
[Jest](https://facebook.github.io/jest/docs/tutorial-react.html) etc. I’m a
strong believer in:

> use whatever languages, frameworks and tools you like and #JFDI.

Once you have those packages installed we will first create a TypeScript
configuration file inside our `tests` directory:

```json
{
	"compilerOptions": {
		"outDir": "./",
		"target": "es5",
		"module": "commonjs",
		"noImplicitAny": false,
		"removeComments": true,
		"sourceMap": false,
		"jsx": "react"
	},
	"include": [
		"./*.tsx",
		"./*.ts"
	],
		"exclude": [
		"node_modules"
	]
}
```

I’m often using `npm test` with a bash script to run my tests, so here’s my
typical script:

```sh
#!/bin/bash
cd tests
tsc
cd ..
./node_modules/.bin/mocha tests/tests/*_test.js
exit 0
```

Basically we are running the TypeScript compiler against our `_test` test files
and we run Mocha on the generated JavaScript files.

Let’s create a first test for our `Hello` component.

```js
import * as React from "react";
import { expect } from "chai";
import { shallow, mount, render } from "enzyme";
import { spy } from "sinon";

import { Hello } from "../src/hello";

describe("<Hello />", () => {
	it("renders the the h1", () => {
		const wrapper = shallow(<Hello />);
		expect(wrapper.find("h1")).to.have.length(1);
	});
});
```

Thanks to [Enzyme](https://github.com/airbnb/enzyme) it’s easy to test our React
component without having to use a full/head-less browser. Depending on the
complexity of your component you might need to use the `mount` function instead
of `shallow` and use something like `jsdom` to emulate a full DOM for React.

Our `Hello` component is simple enough that we’re able to test it via `shallow`
and running `npm test` on a terminal shows us our passing test.

```shell
[dstpierre@roadmap]: ~/projects/tmp>$ npm test

> tmp@1.0.0 test /home/dstpierre/projects/tmp
> ./test.sh

<Hello />
		✓ renders the the h1

1 passing (11ms)
```

So Enzyme and Mocha allow us to develop our React component(s) and not have to
worry about packaging and testing on another project or a debug HTML page.

This pretty much resume the flow of how you can develop/test your React
component:

1.  You create `.tsx` and `.ts` files in the `src` directory
1.  You create tests in the `tests` directory for your TypeScript code
1.  You run your tests and start at #1 until your package is at a place where you
want to see it in an external app.

### Testing the installation of our package

Running `tsc` on the root of our project creates the standard EcmaScript 5
JavaScript files into the `lib` directory.

But before someone can installs our package and imports our `Hello` component we
need to make sure it’s exported and that our package.json file has an entry for
the `main` option.

I usually simply create a `main.js` at the root of my project with the following
code:

```js
exports.Hello = require("./lib/hello").Hello;
```

Pointing to this file in `package.json`:

```json
{
	...
	"main": "main.js"
	...
}
```

We also need to specify which files we want included in the package. We set this
via the `files` entry of `package.json`:

```json
{
	...
	
	...
}
```

And lastly our TtypeScript definition file via the `types` entry of
`package.json` :

```json
{
	...
	"types": "./lib/hello.d.ts"
	...
}
```

*This is just how I like to work, there’s tons of way to do this, so please if
you don’t find this to your taste, adjust to whatever makes you happy.*

In my `package.json` file I usually add a `prepublish` entry on the `scripts`
section like this:

```json
{
	...
	"scripts": [
		"build": "rm ./lib/* && tsc",
		"prepublish": "npm run build",
		"test": "./test.sh"
	]
}
```

The `prepublish` entry will be running every time the npm package is publishing
via `npm publish` or creates a local package via `npm pack`.

When I want to test my package I just create it with `npm pack` this creates a
`.tgz` file that I can install in an external app:

```shell
$> npm install ../my-component/my-component-1.0.0.tgz --save
```

And now on this testing external app we’re able to use the component as normal
via TypeScript or standard es2015:

```js
import { Hello } from "my-component";
...
render() {
	return (
		<div>
			<p>external app</p>
			<Hello />
		</div>
	)
}
```

### That’s it!

We have a fully testable TypeScript npm package with one or multiple React
components ready to be installed locally or published via the public NPM
repository.

You might want to check out `npm link` if you would prefer not to manually pack
the component and install it each time. Depending on what you are building this
might be a smoother flow. I tend to rely on Enzyme for testing and occasionally
test the package in an external app, so the `npm pack` and `npm install` is not
a big thing for me.

> Make sure to click the like button if you have appreciated this tutorial.

* [JavaScript](https://dominicstpierre.com/tagged/javascript?source=post)
* [Typescript](https://dominicstpierre.com/tagged/typescript?source=post)
* [NPM](https://dominicstpierre.com/tagged/npm?source=post)
* [React](https://dominicstpierre.com/tagged/react?source=post)

