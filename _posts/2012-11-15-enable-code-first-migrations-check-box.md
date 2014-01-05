---
layout: post
title: Enable Code First Migrations check box is hidden from Publishing Wizard
date: 2012-11-15 10:25:08 UTC
updated: 2012-11-15 10:25:08 UTC
comments: false
tags: azure
--- 

I’ve been really excited with the new Windows Azure Websites recently. As a
long time [AppHarbor](http://www.appharbor.com) user, deploying ASP.NET MVC
app to the cloud via a git push is extremely nice.

Something strange happened to me on my last project though. After creating a
basic working prototype of the app I was ready to deploy the application.
Since it’s a small prototype for now I was using the Visual Studio Publishing
Wizard.

I had enable migrations on the project, created some models and my
[DbContext](http://msdn.microsoft.com/en-
us/library/system.data.entity.dbcontext(v=vs.103).aspx) class, everything was
working nicely. On the Publishing Wizard I was seeing this:

![without-code-first](/images/without-code-first.png "Publishing Wizzard without code first")

The **Enable Code First Migrations** check box was not present, and my
DbContext class was not displayed. The wizard should look like this:

![with-code-first](/images/with-code-first.png "Publishing Wizzard with code first")

After a couple of time (hours?) of digging, I finally find where the thing is
configured.

![publish-settings](/images/publish-settings.png "Publishing settings")


### Manually editing the pubxml file to make the check box appear

You can enable that check box by editing the Properties\YourProject -
WebDeploy.pubxml look for the PublishDatabaseSettings

```xml    
    <PublishDatabaseSettings>  
      <Objects xmlns="">  
        <ObjectGroup Name="Namespace.Models.YourDBClass" Order="1" Enabled="True">  
          <Destination Path="your-connection-string-goes-here" />  
          <Object Type="DbCodeFirst">  
            <Source Path="DBMigration" DbContext="Namespace.Models.YourDBClass, AssamblyName" MigrationConfiguration="Namespace.Migrations.Configuration, Assambly" Origin="Convention" />  
          </Object>  
        </ObjectGroup>  
      </Objects>  
    </PublishDatabaseSettings>
```
  

Change Namespace.Models.YourDBClass by your class that inherits DbContext,
change Namespace.Migratins.Configuration to fit your migration configuration
namespace and Assambly with your Assambly name.

Save and open the publish wizard you will have that check box.

I’ve took the VS Publishing Wizard from [that stackoverflow
question](http://stackoverflow.com/questions/13230902/publishing-my-database-
package-deploy-with-code-first-migrations/13388315#13388315) I’ve answered.

  

Technorati Tags: [Windows
Azure](http://technorati.com/tags/Windows+Azure),[Publishing
Wizard](http://technorati.com/tags/Publishing+Wizard),[EF Code First
Migrations](http://technorati.com/tags/EF+Code+First+Migrations)

