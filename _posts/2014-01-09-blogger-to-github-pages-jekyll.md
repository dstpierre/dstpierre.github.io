---
layout: post
title: "Moving from Blogger to GitHub Pages and Jekyll"
date: 2014-01-09 2:38:36 UTC
updated: 2014-01-09 2:38:36 UTC
comments: true
summary: "This is how I switched from Blogger to GitHub Pages and Jekyll for my blog. Was fun."
image: "gh-jekyll-dir-structure.png"
categories: tips
tags: github jekyll
--- 

I stumble upon a [tweet from Rob Conery](https://twitter.com/robconery/status/411265853955403776) recently 
mentioning that Jekyll was now free on GitHub. I was curious to 1) know what Jekyll could be (being 
developing on the .NET ecosystem for the last 13 years, I’m not super knowledgeable of the ruby world) and 2) 
would that be my exit solution of Blogger.

I started to do my due diligence and got all the information I could read on the subject. I was super excited to 
try it out myself. I started by creating a project page for my [LiveChat Starter Kit project](http://www.dominicstpierre.com/lcsk) 
just to see what it is to have a pre-built GitHub Page built.

### Some useful links to get you started

1. [GitHub Pages](http://pages.github.com/)
2. [Using Jekyll with Pages](https://help.github.com/articles/using-jekyll-with-pages)
3. [Jekyll's website](http://jekyllrb.com/)
4. [Markdown Cheatsheet](https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet)

### Friday 4am, the kids and wife are asleep, time to migrate my blog

> It does not appears to be that complex, let’s do this.

#### Step 1: Create your website’s repository.

You’ve got one site per GitHub account available at YOURUSERNAME.github.io. You can create a 
repository at GitHub with the same name and start adding your content. I created my repository as 
`dstpierre.github.io`.

#### Step 2: Hello Bootstrap

I than created a simple HTML page with Bootstrap 3. I was lurking on [Rob’s repos](https://github.com/robconery/robconery.github.io) 
to see how he did that and what was the parts that I needed.

> A config file, some specially named directory, fair enough!

I than created a file structure similar to this one:

![Jekyll Directory structure](/images/gh-jekyll-dir-structure.png "Jekyll directory structure")

_layouts: page layouts
_includes: small reusable snippets
_posts: contains all the posts
_config.yml: configuration file for Jekyll

Check [my repository](https://github.com/dstpierre/dstpierre.github.io) if you want to have an idea 
of how things could be setup.

The idea is when you push your changes the site is built and static pages are created.

I decided to have one category per post and multiple tags. You can go whichever direction you'd like. The category and 
tag pages are very similar, here's a snippet of my lcsk tag page:

```html
{% include nav.html %}

<div class="container">
    <div class="blog-header">
        <h1>LiveChat Starter Kit</h1>
        <p class="lead blog-description">
          All post related to my open source project LCSK.
        </p>
    </div>

    {% include post_list.html param = site.tags.lcsk %}
</div>
```

And here's the post_list.html:

```html
<div class='row'>
    {% for post in include.param %}
    <div class="col-sm-4">
        <h4>
            <a href="{{ post.url }}">
                {{ post.title }}
            </a>
        </h4>
        <p>
            <small>
                <span class="glyphicon glyphicon-calendar"></span> {{ post.date | date: "%A, %B %d, %y" }} &mdash;
                {% for tag in post.tags %}
                {% unless forloop.last %}
                <span class="glyphicon glyphicon-tag"></span> <a href="/tags/{{tag}}">{{tag}}</a> |
                {% else %}
                <span class="glyphicon glyphicon-tag"></span> <a href="/tags/{{tag}}">{{tag}}</a>
                {% endunless %}
                {% endfor %}
            </small>
        </p>
        <p class="post-summary">
            {%if post.image %}
            <a href='{{post.url}}'><img src="/images/{{post.image}}" style="max-height: 160px;overflow:hidden" /></a>
            {% else %}
            {{post.summary}}
            {% endif %}
        </p>
    </div>
    {% endfor %}
</div>
```

The syntax if farily simple to understand. This is just how you could create a tag page and 
reusing the same HTML for repeating post across all tag and category pages.

#### Step 3: Debugging Jekyll in Windows, ho boy.

I’ve looked at couple of resources, but nothing was really appealing. After a couple of pushes, the site 
broke completely and the page directives like '{% include nav.html %}' were not interpreted anymore.

I compared, re-check, carefully looked at every lines of HTML. I was not able to find the source of the problem.

### Every Windows developer need to have a Linux VM ready at hand

In 1998 I ran an entire year on Red Hat Linux (no graphical interface). My main activities, which consisted of 
created [eggdrop bot](http://en.wikipedia.org/wiki/Eggdrop) to protect my IRC channel in TCL and chatting using 
[BitchX](http://en.wikipedia.org/wiki/BitchX) was really comfortable using only the console.

I have a [Debian](http://www.debian.org) console only on a [VirtualBox](https://www.virtualbox.org/) VM. 
I guess I can install the GitHub Pages gem quicker than trying all sorts of [trickery on Windows](http://bradleygrainger.com/2011/09/07/how-to-use-github-pages-on-windows.html) ;).

I followed the instruction from GitHub https://help.github.com/articles/using-jekyll-with-pages. I had an error 
installing a dependency, [RedCloth](http://redcloth.org/).

I needed to install ruby 1.9.1-dev, thanks to this SO answer: http://stackoverflow.com/a/14246303/316855.

> It’s working, I’m now able to build my site locally on my Debian machine and finally I will be 
> able to see what’s wrong.

I just clone my repository, this is the command I was using to build my site:

```sh
Bundle exec jekyll build –safe
cp -r _site/* /media/sf_linuxshared/_site/
```

*remember, I'm running only the console, to see the resulting HTML pages, I need 
to copy them to my Windows shared folder.*

It’s such easier to see what’s wrong, but hey, **there were no errors**, everything run with no problem. But 
yet, the directive was *still* not interpreted.

> Ok, time to switch project, I have some client works to complete after all… And frankly I 
> do not see how I will be able to fix this.

### Saturday 5am, now I’m ready to make this work

I started by removing all kinds of things, from '{% include %}' to letting only HTML with only 
a '{{ page.title }}'. No compile error but Jekyll was not rendering properly.

I than decided to start from scratch. New layout page, new includes, new index page. At guess what, 
it **worked**.

> I hate when something like that happens, and now I wanted to understand.

I started to try the old files one by one, to finally found that it was the layout.html page that were not working. 
But as of now, I still cannot find why.

Long story short, you’ll need to have Jekyll installed locally to debug your website, and having a 
Linux VM Is the way to go for non-Mac/Linux users.

### Converting your Blogger posts to Jekyll and Markdown

Clearly I was not going to do this myself. There is probably a tool out there that do this.

https://github.com/kennym/Blogger-to-Jekyll

Again, thanks to my Debian VM. But I had an error when installing this gem. The dependency 
[feedzirra](https://github.com/pauldix/feedzirra) needed libcurl3-dev. So I install it.

```sh
sudo apt-get install libcurl3-dev
```

Than Blogger to Jekyll did its magic and converted my posts to markdown files with the proper header parsed by Jekyll. 
Sweat.

> Renaming files, keeping the same URL schema, et voilà.

The last task is to make sure the file name matches your actual Blogger URL, I fixed this easily in my _config.yml file:

```
permalink: /:year/:month/:title.html
```

### Conslusion

It took me longer that what I thoughts / would have expected. But I had a great time doing this and really enjoyed 
using Vim again. For the last 2-3 years I found myself missing more and more the Linux world.