---
layout: post
title: "Visual Studio 2013 & Windows 8.1 day 1"
date: 2013-10-25 00:31:21 UTC
updated: 2013-10-25 00:31:21 UTC
comments: true
summary: "My first day using Windows 8.1 and Visual Studio 2013. I rarely rant, but this post is one."
image: "mswmp-link-error.png"
categories: tips
tags: visual-studio windows
---

Visual Studio 2013 released on October 17th 2013, I decided to give it a try
and replace my current main development environment from Windows 8 and Visual
Studio 2012. Here’s a series of blog post on what I found worth talking and
some problems I’m having along the way.

## Why Windows 8.1?

On the install screen there was a warning message saying that I would need
Windows 8.1 in order to develop Windows 8 app. Alright, let’s do this, why
not. Windows 8.1 could not be that destructive, it’s just .1 more than Windows
8.

If you have an MSDN subscription like I have (BizSpark), you’ll have to pick
the Windows 8.1 N version, which from what I understand is the bulk version.
No way for me to update Windows via the store. Let’s burn that iso and install
Windows 8.1.

On the installation I had two choices:

  1. Keep your files.
  2. Keep nothing.

Well, “keep my files” would probably mean / imply keeping my installed apps as
well. I knew that the Windows store apps were going to be gone, but I was
convinced that the apps I installed, the desktop one, were going to stay.
Programs like Visual Studio 2012.

At my horrify surprises, after the installation of Windows 8.1, BAM, no more
store apps but no more desktop apps either.

Wait a minute, are all my files still there?

Yep, but no software remains, no more Visual Studio 2012. And lots of other
things… What a time waster, but hey I should have paid more attention to that
details I guess.

Ok, it’s not too bad, let just install VS2013 only and hope for the best. Am I
that naïve? I think so.

I installed Visual Studio 2013, Git, SQL Server Express, you know all those
things that I needed to continue my work. All those un-paid hours not fun at
all. I was not enjoying my Windows 8.1 upgrade so far.

Ok, the machine is starting to look decent, time to start my music and start
working with our new friend, VS2013.

Wait a minute, there’s no music player on the base Windows 8.1 install.

No problem, let’s have a look at their Music app. I download it and try to
open it.

Impressive experience. The app load with the loading screen (logo) and closes
immediately.

Arf. I don’t have time for this, let’s try to download Windows Media Player
like I did for Windows 8. It appears like the Windows 8/8.1 N version does not
include it by default. You need the Media Feature Pack.

This is the page in question:

![mswmp-link-error](/images/mswmp-link-error.png "Windows Media Player link error")

I click on the link that says Media Feature Pack, but I arrive at that page
that mention Windows 8 and not Windows 8.1.

![mswmp-wrong-version](/images/mswmp-wrong-version.png "Windows Media Player wrong version")

I try the download, even if I have doubt that it will work. And it did not. I
got the error message saying “The update is not application for your
computer.”

![mswmp-error](/images/mswmp-error.png "Windows Media error on Windows 8.1")

At that time I’m starting to get really discourage. I need music to work, I
cannot complete anything without my music. Working from home, there’s too much
noise to cover.

Alright Microsoft, you are against me today, no problem I will craft my own
player just to get this day and maybe tomorrow going. Here the simple WPF app
I wrote to play wma and mp3 file.

But the day could not have gone worst. It appears that I need at least Windows
Media Player 10 for the MediaElement to play wam and mp3.

This is starting to get ridiculous.

Let’s try our oldest and best friend Google to fix the day, again. I decided
to try and search for the Windows 8.1 N version of the Media Feature Pack. The
link that should have been there all the time in the page above.

FINALLY. I came to the Windows 8.1 page to download Windows Media Player.
Here’s the correct link.

[http://blogs.msdn.com/b/robmar/archive/2013/10/21/download-media-feature-
pack-for-n-and-kn-versions-of-
windows-8-1.aspx](http://blogs.msdn.com/b/robmar/archive/2013/10/21/download-
media-feature-pack-for-n-and-kn-versions-of-windows-8-1.aspx)

## What it has to do with Visual Studio 2013 you say?

Well, I admit it’s not directly his fault, but still I have lost considerable
amount of time for extremely lazy errors.

I would have expect a smoother upgrade, one that would have kept my desktops
software from Windows 8 to 8.1

I would have expect a better QA for links on the Microsoft site. They clearly
did a copy / paste and did not changed the link.

What the heck is wrong with the music in Windows 8. How come there’s no way to
play music as is and easily. What’s wrong with their Music app, which open,
load first init screen and simply die, no error message nothing.

Are the quality are that low these days. Apple is not better, my wife and kids
both have iMac and I must admit that the quality of the OSes (Windows & Mac)
are clearly declining.

True more and more users do have tablet and small devices. But there still
business people working with desktops and we (developer) who make that echo
system around the OS cannot lost that amount of time because of lowering
quality in the main piece of a desktop computer.

Let’s hope the following days with Visual Studio 2013 will be happier.

/end of rant…

