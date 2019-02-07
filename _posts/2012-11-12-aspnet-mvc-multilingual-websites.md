---
redirect_from: "/2012/11/aspnet-mvc-multilingual-websites.html"
layout: post
title: ASP.NET MVC Multilingual websites
date: 2012-11-12 16:37:50 UTC
updated: 2012-11-12 16:37:50 UTC
comments: true
summary: "A complete ASP.NET MVC multilingual site covering database and querying data, controller and views, images and javascript."
categories: web
tags: asp-net mvc
---

I used to have a series of posts on how to achieve a functional multilingual
website, unfortunately it has been lost. Here is a single to the point post on
the subject, including how I did JavaScript translation when I was running
Bunker App.

### The database

It might look simple at first to just add a new column on your tables
containing the language or culture for the data. But it’s actually a little
bit bigger. It’s normally a little of the columns that needs to be translated,
so you often have to split your tables. An example:

#### Website content page

You might have the ”master” data like who created the page, the date of
creation and last modification etc. Those meta information does not need to be
translated. The page title, URL and the page body content on the other hand
needs to be translated. You might end up with the following tables layout:

|Pages|PageContent|
|----|----|
|id|PageId|
|ParentId|Culture|
|CreatedBy|URL|
|Created|Title|
|ModifiedBy|Body|
|Modified| |

That way you can have 2 or 20 translated version of a page without repeating
the non-translated information. You might also consider adding an index to the
Culture field since you would probably have lots of filtering by that field.

### MVC multilingual routes

I’m usually using the simple {lang} parameter to differentiate the route in a
website with multiple language supported.

```c#
    routes.MapRoute(
      name: "ML",
      url: "{lang}/{controller}/{action}/{id}",
      defaults: new { lang = "en", controller = "Pages", action = "Show", id = UrlParameter.Optional }
    );
```


#### You have a base controller right?



On a base controller class you might override the
[OnActionExecuting](http://msdn.microsoft.com/en-
us/library/system.web.mvc.controller.onactionexecuting(v=vs.98).aspx) and grab
the requested language, set the current thread’s
[CurrentCulture](http://msdn.microsoft.com/en-
us/library/system.threading.thread.currentculture.aspx) and
[CurrentUICulture](http://msdn.microsoft.com/en-
us/library/system.threading.thread.currentuiculture.aspx) appropriately:

```c#
    protected string Language { get; set; }

    protected override void OnActionExecuting(ActionExecutingContext filterContext)
    {
      base.OnActionExecuting(filterContext);

      if (filterContext.RouteData.Values.ContainsKey("lang"))
        Language = filterContext.RouteData.Values["lang"].ToString().ToLower();
      else
        Language = "en";

      ViewBag.language = Language;

      Thread.CurrentThread.CurrentCulture = new CultureInfo(Language);
      Thread.CurrentThread.CurrentUICulture = new CultureInfo(Language);
    }
```

Now you have the Language property pointing to the requested culture (i.e. en-
US / fr-CA, etc) and on the ViewBag.language dynamic property to use on your
views.



### Querying data



No matter how you proceed to query any data from your data store you will
simply need the Language property from your controller to get the desired data
filtered in the right culture.



```c#
     public PageController : BaseController

     {

       private IPageRepository pages = null;

     &nbsp_place_holder;

       public PageController() : this(new SqlPageRepository()) { }

     &nbsp_place_holder;

       public PageController(IPageRepository repo)

       {

         pages = repo;

       }

     &nbsp_place_holder;

       public ActionResult Show(string id)

       {

         var vm = pages.Get(id, Language);

         return View(vm);

        }

     }
```

On line 1 by inheriting from the BaseController the Language property will be
set the the requested culture, then on line 14 we can see that we can pull the
right data from the data store.



### One view to rule them all



You do not need multiple views per language supported. That option would not
scale right and would have lots of overhead. On a view you can easily have the
following elements properly translated:



#### The text



This is the easiest part. The [Resource Files](http://msdn.microsoft.com/en-
us/library/7zxb70x7(v=vs.80).aspx) are very handy for that situation. I used
two approach in the past.



**For simple website, global resource files**



If you have a website that is not too big in number of pages needed to be
translated, you might want to use the easy way. Add the App_GlobalResources to
your project and put one or more file there with translation.



Home.resx | Home.fr.resx | Home.es.resx



The fail over file “Home.resx” is the one with your default language and the
one that will be used for a specific resource value if the key does not exists
on the language file.



I find it easier to create a website with only the default file involved, and
once the site gets to the v1 release I copy the file to the other supported
languages. That way it prevent from having to manually maintain files while
you constantly add new keys to the main file.



To get a value from a view you simply use that syntax:

```html
    <h1>@Resources.Home.Title</h1>
```




Once all your website if fine in your default language you can start
translating to other languages. It’s easy enough to send the .RESX files to
professional translator.



**For larger website, use local resource files**



Sometime things get out of control when you start to have too much global
resource files. You might want to use local resource files on each of your
views subfolder.



#### Images



My approach for handling images on your views in simple enough. Ask the
designer to maximize images without text. For images that need to have text on
it, I simply ask the designer to create multiple file like image-name-en-
US.png, image-name-fr-CA.png, image-name-es-ES.png.



From your views you already have the proper culture on the ViewBag.language
property, so it’s simple enough to get the right image like this:

```html
    <img src="/images/image-name-@(ViewBag.language).png" alt="@Resources.Home.ImageNameAlt" />
```




You just have to make sure every images has the same naming convention. Also,
if you need to set different height and with for different culture, you can
still use the resource files for that, like I did for the alt attribute.



#### Multilingual JavaScript



Ho boy! First of all I’m not really what we could call a JavaScript expert,
even not a JS ultra knowledgeable person. But I come to really love the
language over time and here is what I did to have a functional multilingual
one page JavaScript app when I was the owner of Bunker App.



I tried to replicate what resource files was to the .NET world and simply
created one JS file per culture containing a dictionary of translated value.
Here is a simple example:

```javascript
    // languagepack-en.js
    var Lang = {
      globalDateFormat = 'mm-dd-yy';
      greeting: 'Hello'
    };

    // languagepack-fr.js
    var Lang = {
      globalDateFormat = 'yy-mm-dd';
      greeting: 'Bonjour'
    };
```




Just create an object containing the key and value of the text you are using
on your JavaScript scripts. You can use this like this:

```javascript
    alert(Lang.greeting + ' ' + userName);
```




You only need to include the right language pack file before your scripts,
like this:

```html
    <script language="javascript" type="text/javascript" src="/content/js/langpack-@(ViewBag.language.SubString(0, 2)).js"></script>

    <script language="javascript" type="text/javascript" src="/content/js/yourscript.js"></script>

```



We only use the left two character of the culture, because you might not have
multiple JS file for English and French for different countries. But you can
if you want.



### This is just my way of doing this



Since we have two official language here in Quebec (French and English) almost
every site I've built required to be at least support those two language. I
have been successful using those technique, but I’m not implying in any way
that this is the only or a good approach to use.



Even though having a multilingual website appear to be simple at first, it
clearly add lots of gotcha and things to think a bit. But in the end it’s not
so bad to implement.



Technorati Tags: [Multilingual
website](http://technorati.com/tags/Multilingual+website),[ASP.NET
MVC](http://technorati.com/tags/ASP.NET+MVC),[Multilingual
JavaScript](http://technorati.com/tags/Multilingual+JavaScript)




1er 2014 - 30 sept 2014

1.877.217.5118 694 4823