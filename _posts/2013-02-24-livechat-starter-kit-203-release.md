---
layout: post
title: LiveChat Starter Kit 2.0.3 release
date: 2013-02-24 14:17:08 UTC
updated: 2013-02-24 14:17:08 UTC
comments: true
summary: "First release of LCSK using SignalR, no database and single folder deployment."
categories: open-source
tags: lcsk
---

This release is a small one. Since the Microsoft.AspNet.SignalR package is now
a production ready one, I’ve updated the
[LCSK](http://livechatstarterkit.codeplex.com) NuGet package to use that
dependency.

### Some Changes

The 1.0.0 production release introduce a couple of changes that broke the
latest code from the visitor experience’s video.

The main changes was in the way javascript client code call server methods and
the way the server hub call client side function.

**Before the changes**
    
    myhub.DoSomething();  
    myhub.clientDoOtherThing = function() { }

  

**Now with the 1.0.0**
    
    myhub.server.DoSomething();  
    myhub.client.clientDoOtherThing = function() { }

  

Not a huge changes, but still, enough to break the javascript code files.

  

**On the server, you need to call the client methods like this:**
    
    Clients.Client(connectionId).clientDoOtherThing();

  

The basics functionalities are currently working nicely on LCSK. The following
months will be some small improvement and addition of features.

  

### Available on GitHub

  

I’m currently using both CodePlex and GitHub to host the source code of LCSK.
Since the project is hosted at CodePlex since 2007, I don’t want to loose any
SEO / traffic juice I already have. But on the other hand, I’m more attracted
by GitHub lately so I’m using both for now. Here are the links:

  

GitHub: [https://github.com/dstpierre/lcsk](https://github.com/dstpierre/lcsk)

  

CodePlex: [http://livechatstarterkit.codeplex.com](http://livechatstarterkit.c
odeplex.com)

  

&nbsp_place_holder;

  

As always, feedback are appreciated. If you want to contribute please make
sure you use the GitHub reposotory.

  

Technorati Tags: [LCSK](http://technorati.com/tags/LCSK),[Live
Chat](http://technorati.com/tags/Live+Chat),[Live
Support](http://technorati.com/tags/Live+Support),[Open Source](http://technor
ati.com/tags/Open+Source),[C#](http://technorati.com/tags/C%23),[SignalR](http
://technorati.com/tags/SignalR)

