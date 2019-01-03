---
permalink: "/see-your-stripe-mrr-and-monthly-stats-in-your-terminal-c3110c0891b9"
layout: post
title: "See your Stripe MRR and monthly stats in your terminal"
date: 2017-06-17 11:26:37 UTC
updated: 2017-06-17 11:26:37 UTC
comments: false
summary: "..."
---

![Code editor showing the main calculation logic and the terminal output for
termrr](https://cdn-images-1.medium.com/max/2000/1*iTIzOn0vchY7fkXh4q-uBA.png)


Terminal for me is where I’m feeling at home, certainly because I’m [visually
impaired](https://makermatters.space/going-blind-inspired-me-to-quit-my-own-startup-to-build-something-new-df3911be27b),
nonetheless, I’m always either reviewing/writing code or doing sysops
operations.

At [Roadmap](https://roadmap.space/) we do have Stripe notifications in Slack.
We’re not yet using a Stripe analytics tool like
[Baremetrics](https://baremetrics.com/) for example and I wanted a simple and
quick way to answer the question:

> What’s our current MRR?

Having a manually calculated spreadsheet MRR is all fine, but let’s be honest it
will not get updated frequently and it’s hard to calculate the right amount.

I searched for a command-line tool that would do something like that. I’ve found
some nice [dashboards by
Segment](https://github.com/segmentio/metrics-stripe-subscriptions) and some
Python and [Ruby](https://gist.github.com/siong1987/97b5d8f083675f5641de) GitHub
gist. Nothing either in Arch repo.

I opened my editor and wrote a Go package and I named it
[termrr](https://github.com/dstpierre/termrr) for Terminal MRR. This is a sample
of the output it generates:

```shell
$> ./termrr -key sk_test_32ajah7QA-not-a-real-key
MRR is  425.00
Month over month stats
=====================================
2017-05 New customers: 1 MRR: 25.00
2017-04 New customers: 1 MRR: 50.00
2017-03 New customers: 1 MRR: 55.00
2017-02 New customers: 2 MRR: 80.00
2017-01 New customers: 1 MRR: 95.00
2016-12 New customers: 2 MRR: 120.00
```

I wanted to have a properly calculated MRR. Handling customer and subscription
discounts as well as yearly subscriptions.

The month over month stats show how many new customers were added and how many
new MRR it brought to your total.

#### How to install

You may download a pre-built binary for Linux, Mac and Windows from the
[Releases](https://github.com/dstpierre/termrr/releases) page of the GitHub
project.

If you already have Go installed, just get the package and run `termrr` from
your terminal:

```shell
go get github.com/dstpierre/termrr
```

#### How to use it

```shell
./termrr -key your_test_or_live_stripe_key
```

You use the `-key` argument and pass a Stripe live or test key. It is basically
only iterating over all your customers and get their subscriptions to calculate
the MRR and monthly stats. No write is sent, only read-only of course.

#### What’s next?

I’m not sure, for now I needed that information and that’s what I built. As the
requirements change I might adjust. I’m switching from UI to terminal for some
of the tools I’m using to gain productivity. Email client is my next target.

I must admit that a real-time terminal dashboard could be interesting ;).
Thoughts?

* [Startup](https://dominicstpierre.com/tagged/startup?source=post)
* [SaaS](https://dominicstpierre.com/tagged/saas?source=post)
* [Stripe](https://dominicstpierre.com/tagged/stripe?source=post)
* [Terminal](https://dominicstpierre.com/tagged/terminal?source=post)

