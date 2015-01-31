---
layout: post
title: "LiveChat Starter Kit v2 is in development"
date: 2012-11-12 14:25:09 UTC
updated: 2012-11-12 14:25:09 UTC
comments: true
summary: "The rewrites of LCSK using SignalR as the communication mechasim."
categories: open-source
tags: lcsk
---

I happened to have try two of the currently popular live chat app out there
recently for one of my product [Bunker App](http://www.bunkerapp.com). The
first one was [Olark](http://www.olark.com/) and the second was
[SnapEngage](http://www.snapengage.com/). The experience I went thru while
testing those two apps bring me to the conclusion that I’m not done yet with
my open source [live chat project LiveChat Starter
Kit](http://livechatstarterkit.codeplex.com) (LCSK). Here is what have been my
experience.

### Olark was having a major problem for me

Bunker App is a single page JavaScript app. Put simple, there is no full page
refresh made when requesting a view or posting data. Those kind of app are
more and more popular with HTML5 being adopted.

I’ve contacted multiple times the Olark’s support to try to explain that their
chat widget was not refreshing the correct status of the online of offline
operators since they were requiring a full page request, hence a new HTTP
request to their JavaScript file. After some discussions, a support person
came to my site and try the widget and report back that he was not having that
situation. I than explain that it’s completely fine on page where visitors
browse from page to page where the full page is refreshed, but not on the
application where the main HTML page is loaded only once.

They ended up telling me that their engineer told them that it is working /
not possible (despite me having sent proofing screenshots). I ended up wasting
time with them, and returned back with my original solution (my own project).

Note that Olark for standard website is working perfectly fine, and they have
a great product and user interface, you just can’t told them that their widget
is not refreshing correctly without a full refresh, because they will not
accept that truth ![Winking smile](http://dominicstpierre.net/content/blogimgs
/Windows-Live-Writer/LCSK-v2-is-in-development_12F80/wlEmoticon-
winkingsmile_2.png).

### SnapEngage has the same problem but…

After a month or two, I decided to try SnapEngage, same issue with the non
refreshing status if there is not a full page refresh. I’ve contacted their
support (via chat, compared to email for the other situation). They
immediately understood the situation and told me right away that their widget
was indeed not going to refresh properly within a single page JavaScript app.
They have a workaround for that situation though that involve calling some
JavaScript to force an update etc.

I was really impressed by their app and UI, I was sad that it did not work out
straight out of the box.

### LCSK is not that far compared to those established products

I then return to my good old project, the one with no designer, no sexy UI,
but the one that is refreshing correctly and display the current state of the
operators ;). I decided to improve the product and start the development of
the version 2. I realized that LCSK could be a viable commercial product
almost as is, with the help of a good web designer and some hard hours, I
could easily turn this into a product and try to market this. But, I will not.
I’m a huge fan of open source, this project exists since 2007, I’m just more
motivated of getting a solid result.

The goal of this project was to get you started with a basic live chat app
that give you real-time visitor details, chat from visitors and engaging
visitors to chat. I have interesting plan for the next version.

Technorati Tags: [Live Chat](http://technorati.com/tags/Live+Chat),[Open Sourc
e](http://technorati.com/tags/Open+Source),[LCSK](http://technorati.com/tags/L
CSK)

