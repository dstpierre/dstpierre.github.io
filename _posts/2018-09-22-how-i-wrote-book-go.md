---
permalink: "/how-i-wrote-a-book-on-go-your-turn-now-103a7b06873"
layout: post
title: "How I wrote a book on Go, your turn now"
date: 2018-09-22 11:26:37 UTC
updated: 2018-09-22 11:26:37 UTC
comments: false
summary: "..."
---

Pure vanity was the reason I decided back in October 2017 to write a book.

Now let me elaborate a little bit. I wanted to have a hard-copy printed book
with my face and name on it standing on my desk. I felt I had built enough SaaS
to have sufficient domain knowledge to write a book. Since 2014 I’m enjoying Go
like no other language I’ve worked with so far, why not write a book? How *hard*
can it be after all?

It turns out it’s pretty hard, **way more** than what I could have imagined. But
if I did it, **you can too** and to be frank, it feels very good. I sometimes
call myself an author.

### Using markdown to write the book’s content

I wanted something comfortable, fluid and familiar. I’m running Arch Linux, and
I don’t have any word processor installed. Using Google Document felt wrong.

I started looking around to see if using markdown was an option. Then I stumbled
on [Pandoc](https://pandoc.org/). A superb package that is capable of turning
markdown files into eBook generating PDF and ePub version. No time to test if it
works or not, time to write the book.

I trusted the tool and started writing the book in markdown. I’ve read the
minimum documentation to know that using the # Header 1 and ## Sub-Header 2 was
enough for Pandoc to generate a beautiful table of contents using the H1 as
chapters and H2 as sub-chapter.

I created one chapter per `.md` file. This is an example of chapter one and how
the table of contents is rendered on the PDF.

![Showing a sample markdown file and the table of contents from the PDF.](https://cdn-images-1.medium.com/max/2000/1*AORMYXQ9WsH84fXGdLeomA.png)


I used a GitHub private repository because I knew I wanted to invite some people
that could help from the beginning and give valuable feedback while I was
writing the book.

Collaborators can create a pull request and help you get your idea clear while
the book gets written. I was writing a chapter and could get some feedback right
away before starting the next one. Markdown is very easy to merge pull requests
into the master branch.

### The imposter syndrome

In December 2017 I started to second guess myself. Am I writing a book? I have
no audience, almost no one know me in the Go community. I live 2 hours from
Montreal (in the middle of the wood), so I’m not doing any meet-up talks, etc.

And then you stop. Fear, uncertainty, and *fear*. Some people have those
feelings when they are about to launch their product; I never felt that.
Excitement, happiness, and relieves are my feelings when launching a product.

To me it was kind of over, I wrote the sample chapter and created a simple
landing page. One chapter took me about two months to write, and I had 11
chapters left to write 11 x 2 = 22 months if my math is correct. Discouragement
added up to the other feelings.

If you reach a similar position, you probably made the same mistake I did.
You’re not talking enough about you writing this book. Tweet about it, blog
about it, talk to friends, families, colleagues everywhere you have a chance say
you’re writing this book.

One day [Dean Layton-James](https://medium.com/@DeanLJ) asked me about the book.
He was interested in pre-ordering the book because he was starting a project in
Go and wanted to double check the progress of my book. He pointed me to some
authors that enabled pre-order while they were writing their book. Sometimes it
could take up to one year before the book would be ready. Could that be a good
path for me?

I was intrigued and decided to try it out. I created the quickest Stripe payment
page you can think of, emailed my list, tweet about this announcing that the
book was in pre-order starting at $15 and going up at each chapter release.

To my big surprise, this worked. 💥 Ten minutes after emailing my list I already
had two sales. The entire day I received new orders. I was excited, happy and
energized to continue with the book. From there I wrote about 1 chapter per
month with the accompanying source code to the release on September 2018. In
continually received new order each time I was releasing a new chapter, it felt
really great and kind of validated that there’s a demand for that book.

### What I used as the landing page and email capture

For the landing page, I used [Carrd](https://carrd.co/) and Mailchimp to capture
email. I was sending the sample chapter via MailChimp’s confirmation email. This
work great. People would need to double opt-in to receive the sample chapter.

To capture payment, I created a page with Stripe’s checkout and a quick Go
handler to process the payment, send the email with a link to the book zip file.
No more than 15 lines of Go code and I was capturing credit card, sending the
book and building a csv file of customers that I could use to send downloadable
file when new chapter were added.

It was just enough to start. I’ve looked at services like Gumroad and similar,
but it did not worked for me, YMMV.

I was sending an email to my list each time a new chapter was released and that
the price went up + $2. If I ever write another book, I’ll most certainly re-use
this approach. I was able to start taking order for the book in February 2018, 8
months before the book was released.

### How I generate the book files

At first, I used a virtual machine with Ubuntu and Pandoc installed, and then I
discovered this helpful docker image that can generate the PDF straight from my
Arch Linux machine. Note that I did not wanted to install Pandoc from the Arch
repo, the package is huge and install lots of Haskell dependencies and what not.
The docker image is perfect. This is how I generate my book.

```sh
#!/bin/bash
echo "Generating pdf"
docker run -v `pwd`:/source jagregory/pandoc --listings -H listings-setup.tex -N -S -s -o basaig.pdf title.txt 01.md 02.md 0x.md --epub-cover-image=cover.jpg --toc --latex-engine=xelatex --top-level-division=chapter --base-header-level=1 --highlight-style espresso

echo "Generating epub"
docker run -v `pwd`:/source jagregory/pandoc --listings -H listings-setup.tex -N -S -s -o basaig.epub title.txt 01.md 02.md 0x.md --epub-cover-image=cover.jpg --toc --toc-depth=2 --latex-engine=xelatex
```

You simply list your markdown files in the order you want your chapters to be
generated.

This is the template I use to render the code blocks with the light gray
background and the line numbers.

```latex
% Contents of listings-setup.tex
\usepackage{xcolor}

\lstset{
		basicstyle=\ttfamily,
		numbers=left,
		numberstyle=\footnotesize,
		stepnumber=1,
		numbersep=5pt,
		backgroundcolor=\color{black!10},
		showspaces=false,
		showstringspaces=false,
		showtabs=false,
		tabsize=2,
		captionpos=b,
		breaklines=true,
		breakatwhitespace=true,
		breakautoindent=true,
		linewidth=\textwidth
```

This render code blocks like this one:

![Code block sample](https://cdn-images-1.medium.com/max/1000/1*Huw88HQjVfkvCABMVXerSg.png)


It cannot be easier. Everything is inside clear markdown files. One per chapter.
And generating the PDF and ePUB takes less than a couple of seconds.

### How to market your book

I’m still trying to figure this one out sorry :). This is the actions I did so
far:

1.  Tweet as much as possible about the book.
1.  Posted on Reddit `/r/golang`in my case to announce the book.
1.  Tried to link from some blog posts (I should be blogging more often)
1.  Announced it on some forums I’m a member of.
1.  Emailing my subscribers at least once a month to talk about the book’s progress.
1.  Contacted one influencer [Rob Conery](https://medium.com/@robconery) in the hope
that he could talk about the book, but still no news from him so far (I’ll not
stop trying Rob sorry haha).
1.  I’ve started a Discourse community to stay available for questions and other
related to the book from readers.
1.  I’m going to start live streaming building a tiny real-world example using the
content from the book

In any case, if you’re interested in writing a book or creating an online course
but are not sure if you should or think you cannot do it, I’d say **just try**.
I’m a non-native English speaker without an audience, and I did it. It’s hard
but doable.

* [Writing](https://dominicstpierre.com/tagged/writing?source=post)
* [Golang](https://dominicstpierre.com/tagged/golang?source=post)
* [Books](https://dominicstpierre.com/tagged/books?source=post)
* [Pandoc](https://dominicstpierre.com/tagged/pandoc?source=post)

