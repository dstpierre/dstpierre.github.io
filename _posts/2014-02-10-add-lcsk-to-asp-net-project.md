---
layout: post
title: "Add LiveChat Starter Kit to an ASP.NET project"
date: 2014-02-10 10:54:25 UTC
updated: 2014-02-10 10:54:25 UTC
comments: true
summary: "Small video showing how easy it is to add LCSK to an ASP.NET project."
image: "lcsk-addtoproject.jpg"
categories: open-source
tags: lcsk
---

LCSK is a small and lightweight **free and open source** live chat / live support 
application. It uses [Microsoft SignalR](http://www.asp.net/signalr) for communication mechanism. 
No database involved, only a single folder and the SignalR dependencies are added to your project.

This is an updated video on how to add my open source project 
[LiveChat Starter Kit](https://github.com/dstpierre/lcsk) to your ASP.NET project.

<iframe width="960" height="540" src="//www.youtube.com/embed/DjaO4R1knJE?rel=0" frameborder="0"></iframe>

## Steps for adding LCSK to your ASP.NET project

1. From a Package Management Console type: `Install-Package LCSK`.
2. Copy the 3 JavaScript lines wherever you want the chat box to appear.
3. Start your app and navigate to: **/lcsk/install.html**.
4. Set the password for the admin and agent role.
5. Go to: **/lcsk/agent.html** and sign-in with any name and the agent password.

From there you're up to go, you can simulate a visitor with another browser and request a chat 
by typing into the chat box.

Get more documentation and the complete source from [LCSK's github page](https://github.com/dstpierre/lcsk).