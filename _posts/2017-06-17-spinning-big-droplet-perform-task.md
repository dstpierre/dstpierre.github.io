---
permalink: "/spinning-a-big-digitalocean-droplet-to-perform-a-task-than-destroy-it-from-go-5652c0b64686"
layout: post
title: "Spinning a big DigitalOcean droplet to perform a task than destroy it from Go"
date: 2017-06-17 11:26:37 UTC
updated: 2017-06-17 11:26:37 UTC
comments: false
summary: "..."
---

My development machine has more RAM than our servers at
[Roadmap](https://roadmap.space/). I typically prefer to scale horizontally for
a web app with a higher number of smaller servers compare to having less servers
with higher resources.

This led to a scenario where a specific process we have was eating up all memory
once deploy to the servers. That one task needed vertical scaling.

For background, this is the task in question:

*We’re calculating the engagement of our customers’ users regarding their
product roadmap for a handful of metrics. It’s done once per day and it needs to
traverse the entire users we’re tracking and do some interesting calculation and
month over month comparison. As new customers signs up as well as when existing
customers have more and more users, our total of tracked users grows.*

My first thought as a programmer was:

> I can optimize this, how hard can it be?

It’s in Go (of course), and frankly it is as simple as it can be. Using channels
it is doing the calculation concurrently for all our customers and to be
perfectly honest I was not certain how I could optimize simple calculations like
percentage of usage and simple math like month over month diffs.

I knew that the memory usage was coming from the fact that I was loading the
data in memory, otherwise I was brutalizing the database. But I don’t have time
to find another way for now.

My second thought as a programmer was:

> It worked on my machine right? I just need a similar server to run this process.

How much would it cost to spin a 16GB DigitalOcean droplet and run the task.
It’s actually $0.23 cents per day, I think we can survive that. Our task
currently runs in ~1 minute so we’re paying for the minimum of 1 hour.

This is the code that creates a new DigitalOcean droplet, run a bash script
passing the private network IP of the new droplet and destroy it once everything
is done.

You’ll need two environment variables to execute this Go code, `DO_KEY` which is
a DigitalOcean private access token and `DO_FP` which is a fingerprint of an SSH
key you have there.

#### How we’re executing this

We have a server that’s responsible of running some cron jobs. One of them is a
daily call to a bash script that starts the deployment:

{% include push-content.html %}

```sh
#!/bin/bash

export DO_KEY=pat-key-here
export DO_FP=fingerprint-here
./engagement -init
```

The `vm.go` above is part of our binary that we want to deploy to the new
server. This Go command takes 2 possible arguments:

```sh
-init That basically run the start function in vm.go
-destroy That run the end function in vm.go
```

In `vm.go` we’re creating the droplet and waiting a hard-coded time of 1 minute
and 30 seconds for the droplet to get created so we can grab its private IP
address. We are than passing this IP to another bash script:

```sh
#!/bin/bash
# process.sh

ssh -q -o "StrictHostKeyChecking no" -o UserKnownHostsFile=/dev/null "root@$1" mkdir /root/roadmap
scp -q -o "StrictHostKeyChecking no" -o UserKnownHostsFile=/dev/null ../backup/roadmap/* "root@$1:/root/roadmap"
scp -q -o "StrictHostKeyChecking no" -o UserKnownHostsFile=/dev/null ./engagement* "root@$1:/root/"
ssh -q -o "StrictHostKeyChecking no" -o UserKnownHostsFile=/dev/null "root@$1" "cd /root && ./engagement.sh"
```

This is just an example of what we’re doing. We basically create a new directory
called `roadmap` we than copy a database backup and the binary we want to
execute on the new droplet and than we execute it via another bash script that
simply set some other environment variables and execute the `engagement` Go
command on the new droplet.

The `-q` and the two `-o` arguments ensures that the new host key checking will
not prompt and allows you to copy files via `scp`from the orchestrating server
to the new droplet and run some commands via SSH.

In `vm.go` we’re creating a file named `id.todelete` where `id` is the new
droplet id. We use it to cleanup created droplet(s). In case of any errors, our
`-init` option ensure we’re calling the `end` function before starting the
process. In the event that one droplet would not be destroyed, it would be the
day after.

I would strongly recommend handling the `error` that is returned by the `start`
and `end` functions and act on it to prevent you paying for hours of usage
because the droplet did not destroy properly.

#### I bought time to optimize the process

At time of writing the biggest RAM droplet at DigitalOcean is priced at $0.95
per hour for 64GB. At some point I will need to be way more creative than that
for this specific process. But for now, I have one thing to change to sustain
our growth and it’s the `size: "16gb"` , bumping that number and I’m good and
can continue focusing on other things.

It kind of reminds me how in 1999 we were treating Microsoft SQL Server 7
performance issue by just adding more and more RAM to the server :).

* [Golang](https://dominicstpierre.com/tagged/golang?source=post)
* [Digitalocean](https://dominicstpierre.com/tagged/digitalocean?source=post)
* [Deployment](https://dominicstpierre.com/tagged/deployment?source=post)
* [API](https://dominicstpierre.com/tagged/api?source=post)

