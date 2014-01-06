---
layout: post
title: "My own RSS reader in 1 hour and 172 line of JavaScript"
date: 2013-12-04 17:42:20 UTC
updated: 2013-12-04 17:42:20 UTC
comments: true
summary: "I decided to create my own RSS reader using JavaScript as an experiment."
image: "connectica.png"
categories: open-source
tags: javascript
---

**Major update: I’ve removed jQuery and now I’m at 101 lines of vanilla Javascript.**

> I’ve posted this on Reddit Javascript and Wince made me realize that by having
> jQuery as dependency, the 172 lines of code was not a fact, since jQuery have
> 10k LoC.

> So I’ve took the challenge and rewrite the RSS reader by removing the jQuery
> reference and rewriting lots of the DOM, events code that I had for jQuery.
> And I’m now at 101 lines of JavaScript code, a save of 42 %, it’s crazy.

> I’m leaving the article as it were at first, and I’m posting the entire
> JavaScript code in the bottom. The code at GitHub is updated as well.

Like many others, I find myself “RSS reader-less” when Google decided to [shut
down Reader](http://googlereader.blogspot.ca/2013/03/powering-down-google-
reader.html). Yes, I’ve tried a couple of alternatives, but nothing was really
appealing to me. When I woke up this morning I had a thought, would that be
that long to write a one page HTML / JavaScript RSS reader?

It turns out to be pretty quick and easy, thanks to the [Google Feed
API](https://developers.google.com/feed/) which allow you to get a nice JSON
object from an RSS feed.

![connectica](/images/connectica.png "Connectica screenshot")

## No non-sense amount of files

This is the first project I’m starting since the sale of my last SaaS. I
wanted something quick. I was using [goread.io](https://www.goread.io/)
recently as my Google Reader replacement and I thought, what do I really need
when reading my news?

  1. I need to see the last post for blog I follow
  2. I need to be able to read the post and keep track of which I’ve read
  3. I need to be able to add and remove blog feed

I felt like I was missing posts from my favorite “crew” recently and since I
like to be on the edge as much as possible I had to do something to make sure
I got those articles in front of me.

Let’s open Visual Studio and start a new HTML page where everything will be
there to consume RSS feeds.

## The magic piece of code

_this is the original code using jQuery_

Here’s the most important JavaScript function:

```javascript    
    function parseRSS(url, callback) {  
      $.ajax({  
          url: 'https://ajax.googleapis.com/ajax/services/feed/load?v=1.0&num=10&callback=?&q=' + encodeURIComponent(url),  
          dataType: 'json',  
        success: function (data) {  
          callback(data.responseData.feed);  
        }  
      });  
    }
```
Easy enough. This API convert an RSS feed to JSON.

### Using localStorage and sessionStorage

I’m using a simple JSON object { name, url } to hold all feeds that have been
added in the localStorage.

```javascript    
    function loadFeeds() {  
      $('#feed-list').text('');  
      var feeds = JSON.parse(localStorage.getItem('feeds'));  
      if (feeds == undefined || feeds == null || feeds.length == 0) {  
        addFeed();  
      } else {  
        feeds.forEach(function(f) {  
          $('#feed-list').append(  
            '<span class="feed" data-url="' + f.url + '">' +  
            '  <a href="#" class="switch-feed">' + f.name + '</a>' +  
            '  <a href="#" class="delete-feed">X</a>' +  
            '</span>');  
        });  
      }  
    }
```

> I know, why are you not using any template mechanism. This thing took me in
> all 1.5 hours to write, add to GitHub and deploy to Azure.

I’m simply iterating on all saved feed and displaying a link to load the
latest posts. When we click on a feed the idea is to check if the posts were
already loaded into the sessionStorage and display it.

```javascript    
    $('body').on({  
      click: function (e) {  
        e.preventDefault();  
        var url = $(this).parent().data('url');  
        var feed = sessionStorage.getItem(url);  
        if(feed == undefined || feed == null) {  
          parseRSS(url, function (data) {  
          sessionStorage.setItem(url, JSON.stringify(data));  
          showFeed(data);  
        });  
        } else {  
          console.log('cached');  
          var data = JSON.parse(feed);  
          showFeed(data);  
        }  
      }  
    }, '.switch-feed');
```
This will keep all feed data during all the time of the session and will be
deleted once the browser is closed.

I’m also using localStorage to track which post were opened by using an Array
of URL. Here the two functions that take care of tagging and telling which
post has been read.

```javascript    
    function hasRead(url) {  
      var read = JSON.parse(localStorage.getItem('read'));  
      if (read == undefined || read == null) {  
        return false;  
      }  
      var found = false;  
      for (var i = 0; i < read.length; i++) {  
        if (read[i] === url) {  
          found = true;  
          break;  
        }  
      }  
      return found;  
    }  
    function flagAsRead(url) {  
      var read = JSON.parse(localStorage.getItem('read'));  
      if (read == undefined || read == null) {  
        read = [];  
      }  
      read.push(url);  
      localStorage.setItem('read', JSON.stringify(read));  
    }
```
The only thing missing for my needs was to display the blog post in question
when we click on the title. This also flag the post as read and get the post
from the cached version in sessionStorage.

```javascript    
    $('body').on({  
      click: function(e) {  
        e.preventDefault();  
        $('.content').hide();  
        var parent = $(this).parent();  
        if(parent.hasClass('active')) {  
          parent.find('.snippet').show();  
          parent.find('.content').hide();  
          parent.removeClass('active');  
        } else {  
          parent.find('.snippet').hide();  
          parent.find('.content').show();  
          parent.addClass('active');  
          var url = $(this).data('url');  
          if ($(this).hasClass('new-post')) {  
            flagAsRead(url);  
            $(this).removeClass('new-post').addClass('old-post');  
          }  
        }  
      }  
    }, '.show-post');  
    
	function showFeed(feed) {  
      $('#blogs').html('');  
      feed.entries.forEach(function (f) {  
      var d = new Date(f.publishedDate);  
      $('#blogs').append(  
        '<div class="blog">' +  
        '  <div class="show-post ' + (hasRead(f.link) ? 'old-post' : 'new-post') + '" data-url="' + f.link + '">' +  
          d.toDateString() + ' - ' + f.title +  
        '</div>' +  
        '<p class="snippet">' + f.contentSnippet + '</p>' +  
        '<div class="content">' + f.content + '</div>' +  
        '</div>');  
      });  
      
	  $('.content').hide();  
    }
```
The showPost function is called when we click on a blog name link. This load
all the available posts and hide the content and only display a snippet of the
post. When we click on the post title the snippet is hidden and the content is
shown. It’s not pretty how the HTML is render and most of the code is probably
not optimal. The points was not to create a nice code but to get a working RSS
reader in one file with only HTML, CSS and JavaScript.

## Connectica is born

I named this project connectica and it’s available on
[GitHub](https://github.com/dstpierre/connectica) and it’s also live at
[http://connecti.ca](http://connecti.ca).

**This is the code without jQuery.**

```javascript    
    google.load("feeds", "1");  
    var addNewFeed = null;  
    
	// helpers from Todd Motto: http://toddmotto.com/creating-jquery-style-functions-in-javascript-hasclass-addclass-removeclass-toggleclass/  
    var hasClass = function (elem, className) {  
        return new RegExp(' ' + className + ' ').test(' ' + elem.className + ' ');  
    }  
    var addClass = function (elem, className) {  
        if (!hasClass(elem, className)) {  
            elem.className += ' ' + className;  
        }  
    }  
    var removeClass = function (elem, className) {  
        var newClass = ' ' + elem.className.replace(/[\t\r\n]/g, ' ') + ' ';  
        if (hasClass(elem, className)) {  
            while (newClass.indexOf(' ' + className + ' ') >= 0) {  
                newClass = newClass.replace(' ' + className + ' ', ' ');  
            }  
            elem.className = newClass.replace(/^\s+|\s+$/g, '');  
        }  
    }  
    // my own helper function  
    function isModern() {  
        if ('querySelector' in document && 'addEventListener' in window && Array.prototype.forEach)  
            return true;  
        else  
            return false;  
    }  
    var setVisibility = function(elm, visible)  {  
        elm.style.display = visible ? 'block' : 'none';  
    }  
    var changeInputColor = function (color) {  
        var allInputs = document.querySelectorAll('input[type="text"]');  
        [].forEach.call(allInputs, function (i) {  
            i.style.backgroundcolor = color;  
        });  
    }  
    var resetInputText = function () {  
        var allInputs = document.querySelectorAll('input[type="text"]');  
        [].forEach.call(allInputs, function (i) {  
            i.value = '';  
        });  
    }  
    window.onload = function () {  
        if (!isModern()) {  
            return alert('Your browser cannot run this app.');  
        }  
        addNewFeed = document.querySelector('#add-new-feed');  
                  
        setVisibility(addNewFeed, false);      
        var newFeed = document.querySelector('#new-feed');  
        newFeed.addEventListener('click', function(e) {  
            e.preventDefault();  
            addFeed();  
        });  
        var addFeed = document.querySelector('#add-feed');  
        addFeed.addEventListener('click', function(e) {  
            e.preventDefault();  
            changeInputColor('white');  
            var name = document.querySelector('input[name="feed-name"]').value;  
            var url = document.querySelector('input[name=feed-url').value;  
            if((name == undefined || name == '') ||  
                (url == undefined || url == '')) {  
                changeInputColor('yellow');  
                return;  
            }  
            var feeds = JSON.parse(localStorage.getItem('feeds'));  
            if (feeds == undefined || feeds == null) {  
                feeds = [];  
            }  
            feeds.push({ name: name, url: url });  
            localStorage.setItem('feeds', JSON.stringify(feeds));  
            setVisibility(addNewFeed, false);  
            resetInputText();  
            loadFeeds();  
        });  
        loadFeeds();  
    };  
    function addFeed() {  
        setVisibility(addNewFeed, true);  
        document.querySelector('input[name=feed-name').focus();  
    }  
    function loadFeeds() {  
        var list = document.querySelector('#feed-list');  
        list.innerText = '';  
        var feeds = JSON.parse(localStorage.getItem('feeds'));  
        if (feeds == undefined || feeds == null || feeds.length == 0) {  
            addFeed();  
        } else {  
            var toRemove = document.querySelectorAll('.switch-feed');  
            if (toRemove != undefined && toRemove != null && toRemove.length > 0) {  
                [].forEach.call(toRemove, function (l) {  
                    l.removeEventListener('click');  
                });  
            }  
            toRemove = document.querySelectorAll('.delete-feed');  
            if (toRemove != undefined && toRemove != null && toRemove.length > 0) {  
                [].forEach.call(toRemove, function (l) {  
                    l.removeEventListener('click');  
                });  
            }  
            var buffer = '';  
            feeds.forEach(function (f) {  
                buffer +=  
                    '<span class="feed" data-url="' + f.url + '">' +  
                    '  <a href="#" class="switch-feed">' + f.name + '</a>' +  
                    '  <a href="#" class="delete-feed">X</a>' +  
                    '</span>';  
            });  
            list.innerHTML = buffer;  
            var items = document.querySelectorAll('.switch-feed');  
            if (items != undefined && items != null && items.length > 0) {  
                [].forEach.call(items, function (l) {  
                    l.addEventListener('click', function (e) {  
                        e.preventDefault();  
                        var url = this.parentNode.getAttribute('data-url');  
                        var feed = sessionStorage.getItem(url);  
                        if (feed == undefined || feed == null) {  
                            parseRSS(url, function (data) {  
                                sessionStorage.setItem(url, JSON.stringify(data.feed));  
                                showFeed(data.feed);  
                            });  
                        } else {  
                            console.log('cached');  
                            var data = JSON.parse(feed);  
                            showFeed(data);  
                        }  
                });  
            }  
            items = document.querySelectorAll('.delete-feed');  
            if (items != undefined && items != null && items.length > 0) {  
                [].forEach.call(items, function (l) {  
                    l.addEventListener('click', function (e) {  
                        e.preventDefault();  
                        var parent = this.parentNode;  
                        var url = parent.getAttribute('data-url');  
                        parent.parentNode.removeChild(parent);  
                        var feeds = JSON.parse(localStorage.getItem('feeds'));  
                        if (feeds != undefined && feeds != null) {  
                            var newFeeds = [];  
                            feeds.forEach(function (f) {      
                                if (f.url != url)  
                                    newFeeds.push(f);  
                            });  
                            localStorage.setItem('feeds', JSON.stringify(newFeeds));  
                        }  
                                  
                    });      
                });  
            }  
        }  
    }  
    function showFeed(feed) {  
        var blogs = document.querySelector('#blogs');  
        var posts = document.querySelectorAll('.show-post');  
        if (posts != undefined && posts != null && posts.length > 0) {  
            [].forEach.call(posts, function (p) {  
                p.removeEventListener('click');  
            });  
        }  
        blogs.innerText = '';  
        var buffer = '';  
        [].forEach.call(feed.entries, function (f) {  
            var d = new Date(f.publishedDate);  
            buffer +=   
                '<div class="blog">' +  
                '  <div class="show-post ' + (hasRead(f.link) ? 'old-post' : 'new-post') + '" data-url="' + f.link + '">' +  
                    d.toDateString() + ' - ' + f.title +  
                '</div>' +  
                '<p class="snippet">' + f.contentSnippet + '</p>' +  
                '<div class="content" style="display: none;">' + f.content + '</div>' +  
                '</div>';  
        });  
        blogs.innerHTML = buffer;  
        posts = document.querySelectorAll('.show-post');  
        if (posts != undefined && posts != null && posts.length > 0) {  
            [].forEach.call(posts, function (p) {  
                p.addEventListener('click', function (e) {  
                    e.preventDefault();  
                    var contents = document.querySelectorAll('.content');  
                    [].forEach.call(contents, function (c) {  
                        setVisibility(c, false);  
                    });  
                    var parent = this.parentNode;  
                    var snippet = parent.querySelector('.snippet');  
                    var content = parent.querySelector('.content');  
                    if (hasClass(parent, 'active')) {  
                        setVisibility(snippet, true);  
                        setVisibility(content, false);  
                        removeClass(parent, 'active');  
                    } else {  
                        setVisibility(snippet, false);  
                        setVisibility(content, true);  
                        addClass(parent, 'active');  
                        var url = this.getAttribute('data-url');  
                        if (hasClass(this, 'new-post')) {  
                            flagAsRead(url);  
                            removeClass(this, 'new-post');  
                            addClass(this, 'old-post');  
                        }  
                    }  
                });  
            });  
        }  
    }  
    function flagAsRead(url) {  
        var read = JSON.parse(localStorage.getItem('read'));  
        if (read == undefined || read == null) {  
            read = [];  
        }  
        read.push(url);  
        localStorage.setItem('read', JSON.stringify(read));  
    }  
    function hasRead(url) {  
        var read = JSON.parse(localStorage.getItem('read'));  
        if (read == undefined || read == null) {  
            return false;  
        }  
        var found = false;  
        for (var i = 0; i < read.length; i++) {  
            if (read[i] === url) {  
                found = true;  
                break;  
            }  
        }  
        return found;  
    }  
    function parseRSS(url, callback) {  
        var feed = new google.feeds.Feed(url);  
        feed.load(callback);  
    }  
    function debugObject(obj) {  
        var keys = Object.keys(obj);  
        for (var i = 0; i < keys.length; i++) {  
            console.log(keys[i]);  
        }  
    }
```


## The querySelector and querySelectorAll

I’ve almost switch 1-for-1 the jQuery selector with the document.querySelector
and document.querySelectorAll. I find this extremely confortable after a
couple of minutes. I might not use this the proper way, but at first I tried
to to forEach directly on what the function querySelectorAll return. But I was
getting “Object does not support forEach”. So I ended up doing that trick, but
I would prefer another way, clearly I’m doing something wrong here.

```javascript    
    var list = document.querySelectorAll('.class-name');      
    // this work  
    [].forEach.call(list, function(item) {  
      item.innerHTML = 'this work';  
    });  
    //this does not work  
    list.forEach(function(item) {  
      item.innerHTML = 'this work';  
    });
```

_If someone can explain why the second block is not working, that would be
appreciated._

## Event listeners

Adding and removing event on element is as simple as that:

```javascript    
    var button = document.querySelector('#add-feed');      
    button.addEventListener('click', function() {  
    });  
    // removing the event  
    button.removeEventListener('click'); 
```

The rest of the adventure removing jQuery was just small thing (for this
simple app that is). Like getting the parent of an element with .parentNode.
Checking, adding and removing CSS class, I took helper function from [Todd
Motto](http://toddmotto.com/creating-jquery-style-functions-in-javascript-
hasclass-addclass-removeclass-toggleclass/).

At the end the first iteration was fun, but removing jQuery as a challenge was
more challenging and entertaining. And this is exactly what I’m after now.
