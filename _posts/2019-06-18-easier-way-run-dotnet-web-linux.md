---
permalink: "/easier-way-run-dotnet-linux-with-go"
layout: post
title: "An ~easier way to run ASP.NET MVC 5 apps on Linux with Go"
date: 2019-06-18 09:09:45 UTC
updated: 2019-06-18 09:09:45 UTC
comments: false
summary: "..."
---

First of all, ASP.NET Core runs excellent on Linux, especially for new web 
applications. This article focus on legacy applications that were written in 
the 2008-2012 timeframe against the full .NET framework. What are your options 
if you need to run those applications on a Linux server?

Had you tried porting an ASP.NET MVC 5 to ASP.NET Core? We tried that at my 
current consulting client, and after two months, we decided it was not worth it, 
and we would need a better solution.

### What happened?

Obviously what I will describe here does not apply to any legacy .NET 
applications. But if you do have those issues, you might face the same 
roadblocks that we did.

**LINQ to SQL**

If you're porting a classical enterprise application, your data layer is 
probably one of the first pieces that you will try to port. We had two projects 
that needed to run on Linux.

The first one had a dead standard data access library written with the standard 
`System.Data.SqlClient`. The best decision we made without really knowing it 
in 2007. Crazy how using bare metal is often the safer choice. This DLL could run 
almost as-is against the .NET Core.

On the other hand, the other project use LINQ to SQL and not Entity Framework. 
At the time of writing LINQ to SQL seemed the right choice. Sadly the .NET Core 
does not support this anymore, so we had two options:

* Rewrite it using the standard System.Data.SqlClient
* Rewrite it using the new Entity Framework

Those two options frankly did not sound appealing at all.

We decided to postpone that decision and continue on the project that ran 
as-is to see what kind of work we need to port the ASP.NET MVC app to .NET Core.

**Custom filters, attributes and base controllers**

We wrote custom filters that gave us a massive headache because they used 
missing pieces that are not present on the .NET Core version.

We had lots of `OnActionExecuting` and `OnActionExecuted` that were totally 
broken. This could have been doable to fix the issue, but was it worth the time 
investment?

By a happy chance, we did not use the ASP.NET Membership thing so I cannot 
comment on that, but from the LINQ to SQL experience it might give some trouble 
if you did.

For us though the significant pain came when porting the controllers and all the 
helpers around them was related to the changes that were made to Razor. We 
were generating lots of content via helper functions like this one, for example:

```csharp
protected string GetViewHtml(string viewName, object model)
{
	var content = string.Empty;
	var view = ViewEngines.Engines.FindView(ControllerContext, viewName, null);
	using (var writer = new StringWriter())
	{
		ViewData.Model = model;
		var context = new ViewContext(ControllerContext, view.View, ViewData, TempData, writer);
		view.View.Render(context, writer);
		writer.Flush();
		content = writer.ToString();
	}
	return content;
}
```

Our controllers were mostly OK, then the rest of our derived classes and 
objects relating to ASP.NET MVC were not portable directly. Again a couple of 
choices:

* We need to find how to glue everything and make .NET Core happy.
* Not even sure we had other option at this point.

For completeness, we used an empty ASP.NET MVC Core template, and we added 
pieces by pieces to reform the web application.

At some point, it just felt wrong, and we decided to investigate other options, 
maybe a straight port to .NET Core was not the best way to have that legacy 
ASP.NET MVC apps running on Linux. Since the company wants to move to the cloud 
and because they decided to transition to Linux for easier management and 
cost-effective reasons.

### Can Go be an option here?

I already sold Go's benefits to the company, and there were already a couple of 
processes that had been rewritten to Go. So we decided to explore how we could 
transition an old ASP.NET MVC 5 web application to Go. Can it be more comfortable 
and faster than trying to have it run on .NET Core?

Of course not, claiming this would be silly. We picked the project that had its 
data access layer using System.Data.SqlClient which were already working on 
.NET Core.

We came with the idea of wrapping this around a .NET Core ASP.NET API dummy 
project that could be auto-generated.

It's certainly not a good idea to rewrite legacy code. This DLL is battle 
tested even though there's no unit test. The fact that it's in production since 
2008 and that looking at the Git repo, there's not much bug fixes that were 
applied in the last two years. We can safely say that this piece is solid.

By having an auto-generated wrapper, we could have this DLL, which consisted of 
all data access and business rules available to any HTTP client we wanted. Not 
only on the .NET world anymore.

We used Python for the generator, here's a quick example that you could adapt 
to your use case and have a working ASP.NET MVC Core wrapper around a full 
.NET Framework 4.5 DLL.

The template:

```html
{%raw%}
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using System.Globalization;
using Microsoft.AspNetCore.Mvc;

namespace Your.Namespace.Controllers
{
    [ApiController]
    public class {{ cls['name'] }}Controller : ControllerBase
    {
{% for method in cls['methods'] %}
        // POST api/{{ cls['name'] | lower }}/{{ method['name'] | lower }}
        [Route("api/{{ cls['name'] | lower }}/{{ method['name'] | lower }}")]
        [HttpPost]
        public IActionResult {{ method['name'] }}{% if method['singleton'] %}([FromBody] Your.Namespace.{{ cls['name'] }} {{ cls['name'] | lower }}){% else %}(){% endif %}
        {
    {% for param in method['params'] -%}
        {% if param['type'] == 'string' %}
            var {{ param['name'] }} = Request.Query["{{ param['name'] | lower }}"].ToString();
        {% elif param['type'] == 'char' %}
            var {{ param['name'] }} = Request.Query["{{ param['name'] | lower }}"].ToString()[0];
        {% elif param['type'] == 'bool' %}
            var {{ param['name'] }} = Boolean.Parse(Request.Query["{{ param['name'] | lower }}"].ToString());
        {% elif param['type'] == 'decimal' %}
            var {{ param['name'] }} = Decimal.Parse(Request.Query["{{ param['name'] | lower }}"].ToString());
        {% elif param['type'] == 'int' %}
            var {{ param['name'] }} = Int32.Parse(Request.Query["{{ param['name'] | lower }}"].ToString());
        {% elif param['type'] == 'DateTime' %}
            var {{ param['name'] }} = DateTime.ParseExact(Request.Query["{{ param['name'] | lower }}"].ToString(), "yyyy'-'MM'-'dd'T'HH':'mm':'ss", CultureInfo.InvariantCulture);
        {% elif param['type'] == 'Guid' %}
            var {{ param['name'] }} = Guid.Parse(Request.Query["{{ param['name'] | lower }}"].ToString());
        {% elif param['type'] == 'List<string>' %}
            var {{ param['name'] }} = Request.Query["{{ param['name'] | lower }}"].ToList();
        {% endif -%}
    {% endfor %}
        {% if method['hasOut'] %}
            {{ method['outType'] }} outParam;
        {% endif %}
            {% if method['type'] != 'void' -%}var result = {% endif -%}{% if method['singleton'] -%}{{ cls['name'] | lower }}.{{ method['name'] }}{% else -%}Your.Namespace.{{ cls['name'] }}.Instance.{{ method['name'] }}{% endif -%}({{ method['params_vals'] }}{% if method['hasOut'] %}, out outParam{% endif %});

        {% if method['type'] == 'void' -%}
        return Ok();
        {% else -%}
        {% if method['hasOut'] %}
            return Ok(new {
                results = result,
                outParam = outParam
            });
        {% else %}
            return Ok(result);
        {% endif -%}
        {% endif -%}
        }
{% endfor %}
    }
}
{%endraw%}
```

You cannot use this tempalte as-is. You'll need to handle lots of exceptions found 
in your code base. For example we had multiple methods that were having an `out` 
parameter. We handled that in a JSON like way.

The goal to show this code is to demonstrate how little code we needed to have 
the full .NET Core wrapper and the Go client and structs matching that API.

Here's the code we had to extract the `public` method from that data access 
layer DLL.

{% include push-content.html %}

```python
#!/usr/bin/python3

from jinja2 import FileSystemLoader, Template, Environment
import glob
from os import path
import re

loader = FileSystemLoader('.')
env = Environment(loader=loader)

def readPublicMethods(filepath):
    cls = dict(name="", methods=[], fields=[])
    with open(filepath, 'r') as f:
        prev_line = ''
        for line in f.readlines():
            # Class
            m = re.search(r"public class (\w+)Controller", line)
            if m:
                cls['name'] = m.group(1)

            # Method
            m = re.search(r"^\W+public (?:ActionResult|JsonResult|IHttpActionResult|FileResult) (\w+)", line)
            if m:
                method = {'name': m.group(1), 'params': [], 'verb': 'Get'}

                if prev_line.strip() == '[HttpPost]':
                    method['verb'] = 'Post'

                # Params
                for (_, t, n) in re.findall(re.compile(r"^\s+public \w+ \w+\W?[\(,]\W?(([A-Za-z0-9\?]+(?:<[\w\s,]+>)?))\W(\w+)"), line):
                    method['params'].append({ 'name': n, 'type': t })

                cls['methods'].append(method)
            prev_line = line

    return cls

def cs_to_go_type(cs_type, add_package_prefix=True):
    conv = {
        'void': '',
        'string': 'string',
        'int': 'int',
        'int?': 'int',
        'char': 'rune',
        'bool': 'bool',
        'bool?': 'bool',
        'decimal': 'float64',
        'List<int>': '[]int',
        'DateTime': 'time.Time',
        'Nullable<DateTime>': 'time.Time',
        'Nullable<int>': 'int',
        'Guid': 'string',
        'sbyte': 'int8',
        'double': 'float64',
        'List<string>': '[]string',
        'Dictionary<int, int>': 'map[int]int'
    }


    if cs_type.strip() in conv.keys():
        return conv.get(cs_type)
    m = re.search(r"List<(\w+)>", cs_type)
    if m:
        if add_package_prefix:
            return '[]pkgname.' + m.group(1)
        return '[]' + m.group(1)
    if add_package_prefix:
        return "pkgname." + cs_type
    return cs_type

def generateGoControllers(template, cls):
    if cls['name'] == '':
        return
    cls['name'] =  cls['name'].lower()
    for method in cls['methods']:
        print(cls['name'], method)
        method['name'] =  method['name'][0].lower() + method['name'][1:]
        for param in method['params']:
            param['type'] = cs_to_go_type(param['type'], False)

    if cls['name'] != '':
        with open(path.join('..', 'server', cls['name'].lower() + '.go'), 'w') as f:
            f.write(template.render(cls=cls))

if __name__ == '__main__':
    go_ctrl_template = env.get_template('templates/ctrl.go')
    ctrl_files = glob.glob("../Controllers/*.cs")
    for filepath in ctrl_files:
        cls = readPublicMethods(filepath)
        generateGoControllers(go_ctrl_template, cls)
```

The general idea is to expose all public methods from the DLL as web API 
endpoint from an ASP.NET Core application that will be run on a Linux server. 
We dockerized this app, but it's not mandatory. The idea is also to generate 
the Go code that can call those API endpoints by generating the structures and 
the HTTP client code.

This is the Go code template:

```html
{%raw%}
package legacydal

import (
	"net/url"
	"strconv"
	dal "github.com/org/pkgname/dal"
)

// {{ cls['name'] }}Service ...
type {{ cls['name'] }}Service struct {}

{% for method in cls['methods'] -%}
// {{ method['name'] }} ...
func (svc *{{ cls['name'] }}Service) {{ method['name'] }}({{ method['params_vals'] }}) ({% if method['type'] != '' -%}{{ method['type'] }}, {% endif -%}{% if method['hasOut'] -%}{{ method['outType'] }},{% endif -%} error) {

	q := url.Values{}

	{% for param in method['qparams'] -%}
        {% if param['type'] == 'string' %}
	q.Add("{{ param['name']|lower }}", {{ param['name'] }}) 
        {% elif param['type'] == 'rune' %}
	q.Add("{{ param['name']|lower }}", string({{ param['name'] }})) 
        {% elif param['type'] == 'bool' %}
	if {{ param['name'] }} {
		q.Add("{{ param['name']|lower }}", "true") 
	} else {
		q.Add("{{ param['name']|lower }}", "false") 
	}
        {% elif param['type'] == 'float64' %}
	q.Add("{{ param['name']|lower }}", strconv.FormatFloat({{ param['name'] }}, 'f', -1, 64) );
        {% elif param['type'] == 'int' %}
	q.Add("{{ param['name']|lower }}", strconv.Itoa({{ param['name'] }})) 
		{% elif param['type'] == 'time.Time' %}
	q.Add("{{ param['name']|lower }}", {{ param['name'] }}.Format("2006-01-02T15:04:05") )
		{% elif param['type'] == '[]string' %}
	for _, s := range {{ param['name'] }} {
		q.Add("{{ param['name']|lower }}", s) 
	}
  {% endfor %}

	{% if method['hasOut'] -%}
	var result struct{ 
		OutParam {{ method['outType'] }} `json:"outParam"`
		Results {{ method['type'] }} `json:"results"`
	}
	{% elif method['type'] == '' -%}
	{% else -%}
	var result {{ method['type'] }}
	{% endif -%}

	{% if method['singleton'] -%}
		{% if method['type'] == '' -%}
	err := post("/api/{{ cls['name'] | lower }}/{{ method['name'] | lower }}", q, &entity, nil)
		{% else -%}
	err := post("/api/{{ cls['name'] | lower }}/{{ method['name'] | lower }}", q, &entity, &result)
		{% endif %}
	{% elif method['type'] == '' -%}
	err := post("/api/{{ cls['name'] | lower }}/{{ method['name'] | lower }}", q, nil, nil)
	{% else -%}
	err := post("/api/{{ cls['name'] | lower }}/{{ method['name'] | lower }}", q, nil, &result)
	{% endif %}

	{% if method['hasOut'] -%}
	return result.Results, result.OutParam, err
	{% elif method['type'] == '' -%}
	return err
	{% else -%}
	return result, err
	{% endif %}
}

{% endfor %}
{%endraw%}
```

We wrapped that implementation inside interfaces to allow us to port the code 
to native Go at some point slowly. If we ever need to do that, we could do this 
gradually service by service.

This is the template that render the `struct` and the `interface`:

```html
{%raw%}
package dal


type {{ cls['name'] }} struct {
	{% if cls['parent'] %}
	{{ cls['parent'] }}
	{% endif %}
	{% for field in cls['fields'] -%}
	{{ field['name'] }} {{ field['type'] }} `json:"{{ field['json_name'] }}"`
	{% endfor -%}
}

{% if cls['methods']|length > 0 %}
type {{ cls['name'] }}Service interface { {% for method in cls['methods'] %}
	{{ method['name'] }}({{ method['params_vals'] }}) ({% if method['type'] != '' -%}{{ method['type'] }}, {% endif -%} error){% endfor %}
}
{% endif %}
{%endraw%}
```

Just like the C# generator, we read all .cs files and look for public method 
and the public property of all the classes to generate the Go code.

We were fortunate that all our classes were formatted the same. Code discipline 
is very important, and we could generate those projects in about three days.

For the more problematic project that uses LINQ to SQL, we are thinking of 
using the same technique, but we're not there yet. This will probably be a 
follow-up article. But we are hesitating between those options:

* Create a parser for the LINQ to SQL queries and have them turn to SQL. 
* The LINQ to SQL Context has a way to output the generated SQL. Maybe there's 
something to dig there.

### What about the views?

We did the same for the view, a simple Python script that transforms the Razor 
views into Go HTML templates. This process produces 80% working template, and 
it was quick to have them in a working state.

This is an example if you'd like to craft one that works for your Razor views:

```python
{%raw%}
#!/usr/bin/python3

import glob
import os
import re


def convert_line(line):
    line = re.sub(r'@?Resources\.(\w*.\w*)', r'{{ _t "\1"}}', line)
    line = re.sub(r'@?else if\s?\((.*)\)', r'{{ else if \1 }}', line)
    line = re.sub(r'@?if\W?\((.*)\)', r'{{ if \1}}', line)
    line = re.sub(r'@\{Html\.RenderPartial\((.*)\);\}', r'{{ template \1 . }}', line)
    line = re.sub(r'@?else(?:\s+)?$', r'{{ else }}', line)
    line = re.sub(r'@?foreach\s*\((.*)\)', r'{{ range \1 }}', line)
    line = re.sub(r'@?this.LanguageString\((".*"),\s?(".*")\)', r'{{ onFr \1 \2 }}', line)
    line = re.sub(r'@?for\((.*)\)', r'{{ range \1 }}', line)
    line = re.sub(r'@(Model[\w\.\(\)]*)', r'{{ \1 }}', line)
    line = re.sub(r'([^@\.])(Model[\w\.\(\)]*)', r'\1.\2', line)
    line = re.sub(r'@?Html.Raw\((.*)\)', r'{{ rawhtml \1 }}', line)
    line = re.sub(r'this\.Preferences\(\)', r'.Preferences', line)
    
    return line
    

def convert_view(src, dst):
    print(src, "->", dst)
    lines = []
    with open(src) as f:
        name = os.path.basename(os.path.dirname(src)) + "_" + os.path.basename(src)[:-7]
        for line in f.readlines():
            lines.append(convert_line(line))
        lines.append("\n{{ end }}")
        lines = [('{{ define "%s" }}\n' % name)] + lines[1:]
    if not os.path.exists(os.path.dirname(dst)):
        os.mkdir(os.path.dirname(dst))
    with open(dst, mode="w") as f:
        f.write("".join(lines))


if __name__ == '__main__':
    view_src = glob.glob("../prj/Views/**/*.cshtml") + glob.glob("../prj/Views/*.cshtml") + glob.glob("../prj/Views/Shared/DisplayTemplates/*.cshtml")
    for src in view_src:
        dst = "../views/" + src[19:]
        convert_view(src, dst)
{%endraw%}
```

### In conclusion

Overall it was quicker and easier to have the V and C part of the MVC 
application handled by Go than fighting with .NET Core. Plus we now have a 
fully functional interface package for all the data access. Ready to create 
unit tests and maybe re-implement them in Native Go or not.

Again YMMV and depending on the legacy ASP.NET MVC project, it might be the 
opposite for some. We were happy to upgrade to .NET Core. But it was just way 
too painful for those two use-cases.

I guess we can resume this article in one sentence. If you're planning an 
ASP.NET MVC 5 migration to .NET Core and are having a tough time, it might be 
worth checking at the auto-generating code and turn your legacy application 
into a more maintainable application.