---
layout: post
title: "Deploying Go apps to Azure Web Apps is now easy"
date: 2015-11-13 11:44:00 UTC
updated: 2015-11-13 11:44:00 UTC
comments: true
summary: "I just discovered that Azure's Kudu is handling Go (golang) apps now when pushing a repo to your Azure web apps."
image: "go-azure-deployment.png"
categories: tips
tags: go azure
---

![Go editing with atom](/images/go-azure-deployment.png "Go editing with atom")

You know when you do something by habits and don't really have time to see if there's a better way to do it?

That's exactly what I was doing regarding my Go deployment to my Azure Web Apps.

### Manual deployment ###

Here were my flow for deploying a Go program to Azure Web Apps.

1. Development
2. Git push to Github
3. `go test` and `go build`
4. I had a web.config that specified that I was deploying an httpPlateformHandler app.
5. Was incrementing the exe's filename => v-0.0.1, v-0.0.2
6. Was FTPing to Azure Web Apps and sending the new exe with web.config updated app

The last steps was automatically picking the new exe as the new app to handle requests.

Here's the web.config that I had.

```xml
<?xml version="1.0" encoding="UTF-8"?>
<configuration>
  <system.web>
        <customErrors mode="Off"/>
    </system.web>
    <system.webServer>
        <handlers>
            <add name="httpplatformhandler" path="*" verb="*" modules="httpPlatformHandler" resourceType="Unspecified" />
        </handlers>
        <httpPlatform stdoutLogEnabled="true" processPath="d:\home\site\wwwroot\v-0.0.10.exe"
                      arguments=""
                      startupTimeLimit="60">
            <environmentVariables>
              <environmentVariable name="GOROOT" value="d:\home\site\wwwroot\go" />
            </environmentVariables>
        </httpPlatform>
    </system.webServer>
</configuration>
```

Although this was taking me like 45-60 seconds of manual commands, it was fairly quick to have a new version in production.

### Then I needed to implement websocket ###

I'm currently porting some Node apps and I needed to have websocket handler. On Node when you want to enable websocket you need to eneable it on the Azure portal, and you need to add this to your `web.config`:

```xml
<webSocket enabled="false" />
```

So I try that and deployed my Go app, was not working, I tried to removed it and it was not working as well.

Long story short, I exhausted all my idea and start fearing that Azure websocket and Go app were not going to play nicely. Than I posted a [StackOverflow question](http://stackoverflow.com/questions/33675348/azure-websocket-not-working-with-go-httpplatformhandler/33684891#33684891).

**Xiaomin Wu** replied and I was shocked to see on his Github repo the presence of the Azure button "Deploy to Azure" and he also mentioned that he did not created any `web.config`.

> What? Can we just `git push` to Azure for Go app now?

### Deploying Go app on Azure is now as easy as a C# or Node app ###

I than plugged my Azure Web App to my Github repository and tried to deploy. **It worked**.

This is a major improvement in the flow of deploying changes to production.

I compared the `web.config` that is automatically created by Kudu with mine and there's not much difference. Still don't really understand why my manual deployment did not work with websocket, but I'm just glad we can now deploy Go apps to Azure from a `git push`.

**Flow is now**:

1. Developement
2. `go test`
3. `git push`

### Been loving Go so far ###

This past year was super exciting for me in terms of technologies that I had the chance to discover and work  with. Go is the most recent for me in production and so far I like it very much.

The learning curve is steep, especially from someone like me that's coming from a ~15 year C# background.

But I think it is making me become a better programmer, and it has been a long time since I felt this from a language.
