---
layout: post
title: "The day I ditched Windows 8 for Linux"
date: 2014-06-11 11:25:45 UTC
updated: 2014-06-11 11:25:45 UTC
comments: true
summary: "I've been developing in C# on Windows for the last 14 years and today is the day I decided enough is enough."
image: "linux-desktop.png"
categories: tips
tags: linux, windows
---

![simple C color inverter][/images/linux-desktop.png]

My main source of revenue is by getting paid to architect and develop C# software 
for companies running in Windows infrastructure. Still, since I installed Windows 8 
when it launched, I've been slowed down by my OS. I'm not writing this to start any 
kind of OS war whatsoever. Its just that I decided to try something else.

### I tried to switch before

Those thoughts started in 2008 when I decided to try Mac. I bought an iMac and ran 
on Mac OS X for about 1.5 year. I guess Vista annoyed me at that time. But Mac 
was not for me, too visual.

In 1998 I ran only on a console based Debian Linux, at the time I was almost just 
using BitchX (irc client) and text-based web browsers. I tried many times 
to return to Linux since then. But two days ago, June 9th 2014, I was 
ready to do whatever I had to do to have a Linux OS that should allow me 
to continue supporting my clients and get my work done. But more importantly 
to get an OS that shall let me be in control again, productive and stable.

### What went wrong with Windows 8 for me

1. **Stability**: I never had that many hard drive issues until Windows 8. The store
apps was all corrupted two times. Countless times the OS scanned and tried to repair 
disk errors before opening Windows.
2. **Responsiveness**: Too often I opened Visual Studio, Chrome, Firefox and all of a sudden 
clicking on an icon to make the app appearing just did nothing. All the desktop was 
frozen, ALT+TAB, ALT+F4 nothing was responding, but the current app continued to work.

### What was keeping me from ditching Windows

1. **Camtasia Studio**: I bought Camtasia Studio a couple of months ago. I was sad to have paid $300 and used 
it for about 15-20 screen casts. But still, the issue with Windows 8 greatly out-weight the benefits 
of a single software.
2. **My Steam Games**: I have lots of game that I bought on Steam when Steam did their weekend deals. One 
of my favourite game is Civilization, I have a lots of turn-based games. I'll have to keep them on an older 
PC, that's not terrible.
3. **Visual Studio**: As I've said, I'm primarily a C# developer, Visual Studio is a big part of my revenue stream. 
I don't really have alternative on Linux like Camtasia has, but I may manage to use it via a virtual machine.

### How to continue using Visual Studio in Linux

I have no problem running Visual Studio in a virtual machine using VirtualBox. I picked a Windows 7 
OS to run Visual Studio 2013. I'm planning on installing nothing more in that virtual machine. I will 
try to keep everything else on Linux. I'm willing to take my time and find alternative ways to do 
things that I'm used to do in Windows.

This is the hard part. When you do something that you know it should take you X minutes in Windows. Now you 
have to dig a bit and maybe learn new ways of accomplishing the task. But I enjoy learning new things, I 
don't think that would be much of a pain for me.


### What distribution I choose

Here's how it went:

#### June 9th 2014

**8:00 pm**: The family routine is over, kids are OK, time to attack the installation of ArchLinux.
I was running Arch in a VM on Windows and I liked the lightweight aspect and being in control.

**10:XX pm**: My 1-month-old PC gives me hard time installing Arch, I downloaded Manjaro 
and tried to see if their installer was doing a better job than me.

**10:50 pm**: Looks like I will not be able to work tomorrow. Linux is still not installed. Manjaro also 
failed installing. WTF is that new PC.

**11:XX pm**: Let's try my good old friend Debian.

**11:3X pm**: My network card is not supported. Haaaaa, that UEFI bios, WTH is that. Alright, just 
for curiosity, would Ubuntu work.

**11:5X pm**: Yeah :(. Well... It was certainly not my first choice. But Ubuntu installed without 
any glitch. I need to install VirtualBox, Skype, Windows 7, Visual Studio transfer my backup source 
code to be ready for work in about 7 hours.

**12:XX am**: Alright, I'm too old for this, I will take the A.M. tomorrow and finish this up.

#### June 10th 2014

**5:5X am**: Fresh and clean, I'm now top shape, let's do this... Installing VS 2013 and get my source code in place.

**6:3X am**: Alright, everything seems to be in place, I don't really like how VirtualBox display my 
Windows 7 VM, but for now it works.

**7:XX am**: Then, I decided to try if I could installed Xfce instead of Unity. You know 
that kind of thoughts couple of minutes before starting to work.

What a crap. I don't really remember what I did, but I added couple of packages, did not liked the results 
than I removed couple of packages, and boom! The thing was all crappy suddenly. Holly sh*t. Why did 
I try something like this, I was in a working state, but no...

**9:XX am**: Fine, I will quickly install Windows 7, that should do the trick, I will not 
use Windows 8 and I will be up and ready to work without loosing too much of billable hours.

**9:XX am**: What, my network card is not detected by Windows 7. I don't have time to search for all those 
drivers...

OK, I took a long deep breath and I said:

>> You wanted to go with Arch, you should try it until it work.

**10:XX am**: Let's give Arch another chance.

**11:XX am**: Nope, I failed. Well, OK let's return to Ubuntu and try not to do stupid thing this time.

Yep, that's me. Went to Linux, than Windows 7 than back to Linux in 12 hours. But I did not wanted 
to dual-boot. If you dual-boot it's too easy to be tempted to always return to the easy path and 
it take a longer time to learn.

The moral of the story is, if you are thinking about switching to Linux, make sure you have the following.

1. **Patience**: You'll need a lot a patience.
2. **Read a lot**: Don't expect everything to just work as is immediately. You'll need to dig.
3. **Accept that it's not perfect**: This is a major change, and as with anything we are afraid of changing
something. Accept that it's going to be bumpy and not really the way you've planned or imagined.
4. **Research before acting**: Make sure you search before doing anything. Check which Linux packages you'll need 
to replace software you use in Windows.


When you are in an hurry, it's often not the best time to install a Linux distribution that does not 
work out of the box with your current hardware.


### I'm now in control

I've already my feet wet. I was unable to get a tool _xcalib_ to invert colour of both of my screen. 
Only one of them was inverted. I took some code here and there and compiled a small C utility that 
invert my screens. For a programmer Linux feels right at home. As anything it take times to really 
be comfortable with a new environment. My first task on my list is to learn Kdenlive and continue 
my screen casts the way I was using Camtasia Studio.

If you are a .NET dev, you should at least consider creating a Linux virtual machine on your Windows box. 
Developing Python, NodeJS, Go is much appealing on that OS to say the least.
