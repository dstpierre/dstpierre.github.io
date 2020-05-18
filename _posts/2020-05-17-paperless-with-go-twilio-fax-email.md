---
permalink: "/paperless-go-twilio-fax-email"
layout: post
title: "Build paperless remote friendly process with Go, Twilio and Google Cloud Run"
date: 2020-05-17 07:55:40 UTC
updated: 2020-05-17 07:55:40 UTC
comments: false
summary: "..."
---

Today I'll share how I helped an old 25 years old credit bureau relying on 
faxes and printed papers to go fully remote with all their employees.

Somewhere in March, the province of Québec in Canada announced that all 
non-essential companies must stop doing business on their location. If they 
can offer remote work to their employees, that would be the recommended way to 
work for a while.

This tutorial might help if you need to turn faxes and printed paper processes 
into digital remote-friendly processes. And maybe gain some productivity with 
automation.

Here's what I'll cover in this tutorial:

1. [Ditch those old fax machines](#ditch-those-old-fax-machines)
2. [Build a fax machine with Twilio and Google Cloud](#build-a-fax-machine-with-twilio-and-google-cloud)
3. [Add QR-Code to your document](#add-qr-code-to-your-document)
4. [Dispatching the faxes](#dispatching-the-faxes)
5. [Cost comparison of on-premise fax machines and Twilio](#cost-comparison-of-on-premise-fax-machines-and-twilio)


## Ditch those old fax machines

Banks are still working with faxes. I needed a quick way to replicate the fax 
sending and receiving operations. There's mainly one tool that comes up when 
asking for telephony; it's Twilio.

I wanted something that can scale up and down. Something that I could plug and 
forget, just like a fax machine. I did not want to increase infrastructure 
complexity. I settle for Twilio with Google Cloud.


### Build a fax machine with Twilio and Google Cloud

Here's what you'll need for this part:

* A [Twilio](https://www.twilio.com/) account
* A [Google Cloud Platform](https://console.cloud.google.com/getting-started) account and project ready
* Docker
* Go

Twilio is fantastic, and it's a joy to use their product. Create an account, 
pick a phone number that can send/receive faxes.

Here is the Google Cloud Function that gets called when a new call is made. We 
return the URL of a Google Cloud Run application we'll build next.

```go
package faxreceiver

import (
        "fmt"
        "net/http"
)

// Receiver handle Twilio for new fax and set handler for received faxes
func Receiver(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/xml")
  fmt.Fprint(w, `
    <Response>
      <Receive action="https://your-cloud-run-url-ue.a.run.app/"/>
    </Response>
  `)
}
```

On your Twilio configuration for your fax number, you set the Webhook URL for  
new call received to this function. This function indicates to Twilio what to 
do with the fax. We want to call a Google Cloud Run application.

The reason we want a Google Cloud Run application to run is that we're going to 
detect a QR-Code presence on the fax and add to a queue so it can be dispatch.

Physical fax machines accumulate paper, and someone needs to dispatch those 
papers for further processing. We're going to automate this by placing a 
QR-Code on each document, thus saving that human step in the process.

Here's the Cloud Run components:

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
)

// Fax represent a received fax
type Fax struct {
	From            string `json:"from"`
	To              string `json:"to"`
	RemoteStationID string `json:"remoteStationId"`
	FaxStatus       string `json:"faxStatus"`
	ErrorCode       string `json:"errorCode"`
	NumberOfPages   string `json:"numberOfPages"`
	MediaURL        string `json:"mediaUrl"`
	Error           string `json:"error"`
	QRCode          string `json:"qrcode"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("parsing form")

	if err := r.ParseForm(); err != nil {
		fmt.Println("error parsing form: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("creating fax instance")
	fax := Fax{
		From:            r.Form.Get("From"),
		To:              r.Form.Get("To"),
		RemoteStationID: r.Form.Get("RemoteStationId"),
		FaxStatus:       r.Form.Get("FaxStatus"),
		ErrorCode:       r.Form.Get("ErrorCode"),
		NumberOfPages:   r.Form.Get("NumPages"),
		MediaURL:        r.Form.Get("MediaUrl"),
		Error:           r.Form.Get("ErrorMessage"),
	}

	// if there's any issue with the fax we just return without any
	// errors.
	if fax.FaxStatus == "failed" || len(fax.MediaURL) == 0 || len(fax.ErrorCode) > 2 {
		fmt.Println("received error", fax.FaxStatus, fax.ErrorCode, fax.Error)
		w.Write([]byte("ok"))
		return
	}

	var b []byte
	retry := 3

	for {
		if retry == 0 {
			break
		}

		fmt.Println("downloading fax PDF file: ", retry, fax.MediaURL)
		content, err := download(fax.MediaURL)
		if err != nil {
			fmt.Println("error downloading media file: ", err)
			retry--

			time.Sleep(250 * time.Millisecond)
			continue
		}

		b = content
		break
	}

	if len(b) == 0 {
		fmt.Println("unable to download the media url")
		http.Error(w, "unable to download the media file", http.StatusInternalServerError)
		return
	}

	fmt.Println("detecting qr code")
	code, err := detect(b)
	if err != nil {
		fmt.Println("error while detecting qrcode: ", err)
	}

	fax.QRCode = code

	fmt.Println("adding fax to queue topic")
	if err := queue(fax); err != nil {
		fmt.Println("error while adding to queue: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("returning OK")
	w.Write([]byte("ok"))
}

func main() {
	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func download(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func detect(b []byte) (string, error) {
	filename := fmt.Sprintf("/tmp//%d.pdf", time.Now().UnixNano())

	if err := ioutil.WriteFile(filename, b, 0666); err != nil {
		fmt.Println("error while writing PDF file: ", err)
		return "", err
	}

	out, err := exec.Command("zbarimg", "-q", filename).Output()
	if err != nil {
		fmt.Println("error while executing zbar: ", err)
		return "", err
	} else if len(out) == 0 {
		return "", nil
	}

	output := string(out)
	fmt.Println("zbarimg output: ", output)

	for _, buf := range strings.Split(output, "\n") {
		if strings.HasPrefix(buf, "QR-Code") {
			return strings.Replace(buf, "QR-Code:", "", -1), nil
		}
	}

	return "", nil
}

func queue(fax Fax) error {
	bgCtx := context.Background()
	psc, err := pubsub.NewClient(bgCtx, os.Getenv("GCP_PROJECTID"))
	if err != nil {
		log.Println("unable to create the pubsub client: ", err)
		return err
	}

	topic := psc.Topic("fax_received")
	if err != nil {
		log.Println("error returned by CreateTopic: ", err)
		return err
	}

	b, err := json.Marshal(fax)
	if err != nil {
		log.Println("error while encoding fax to json: ", err)
		return err
	}

	topic.Publish(bgCtx, &pubsub.Message{Data: b})
	return nil
}
```

This is what's happening:

1. When a new fax arrives we get the values from the form post.
2. We try to download the PDF file.
3. We try to detect the QR-Code using zbarimg.

This is the Dockerfile for our Google Cloud Run container:

```dockerfile
# Use the official Go image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.13 as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o server

# Use the official Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3
RUN apk add --no-cache ca-certificates imagemagick
RUN sed -i -e 's/v[[:digit:]]\..*\//edge\//g' /etc/apk/repositories
RUN apk update
RUN apk add --no-cache zbar

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /server

# Run the web service on container startup.
CMD ["/server"]
```

This builds our Go server and uses Alpine to install zbar for the runtime 
container.

To deploy the function and the cloud run container I use Makefile:

```make
build:
	gcloud builds submit --tag gcr.io/your-project-id/name-of-your-cloud-run

deploy: build
  gcloud run deploy --image gcr.io/your-project-id/name-of-your-cloud-run --platform managed --allow-unauthenticated --region=us-east1

function:
	gcloud functions deploy NameYourFunction --runtime go111 --trigger-http
```

I deploy my Google Cloud Function with:

```bash
$> make function
```

And I deploy the Cloud Run application with:

```bash
$> make deploy
```

## Add QR-Code to your document

{% include push-content.html %}

Now it's time to add those QR-Code to the document you send, so when they get 
back, they can be automatically dispatched to the right employee or the right 
place in your system.

The majority of the document templates where I needed these are HTML files. I 
created this simple QR-Code render to an image that can be added to the 
template.

```go
...
import (
	"image/png"
	"net/http"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)
...
func renderQRCode(w http.ResponseWriter, r *http.Request) {
	width := r.URL.Query().Get("w")
	height := r.URL.Query().Get("h")
	code := r.URL.Query().Get("code")

	iw, err := strconv.Atoi(width)
	if err != nil {
		log.Println("error converting with to int: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ih, err := strconv.Atoi(height)
	if err != nil {
		log.Println("error converting height to int: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	qrcode, err := qr.Encode(code, qr.L, qr.Auto)
	if err != nil {
		log.Println("error generating qrcode: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	qrcode, err = barcode.Scale(qrcode, iw, ih)

	png.Encode(w, qrcode)
}
```

This is a simple HTTP handler that returns PNG with the QR-Code. We can add a 
QR-Code in our HTML templates like this:

```html
<img src="https://your-url/qr?w=150&h=150&code=your-code-here" />
```

I'm using QR-Code as I found that they are better at keeping their quality when 
the faxes are printed than re-faxed compare to the traditional barcode. As 
always, YMMV.

## Dispatching the faxes

You'll need to create a Google Cloud Queue. We've seen that our Google Cloud 
Run app is adding items to our queue when a fax arrives.

I'll leave the implementation to you to decide what you do with faxes that have 
and haven't a QR-Code. But here's the code to get you started dequeuing the 
messages.

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
)

// InboundFax represent a received fax
type InboundFax struct {
	From            string `json:"from"`
	To              string `json:"to"`
	RemoteStationID string `json:"remoteStationId"`
	FaxStatus       string `json:"faxStatus"`
	NumberOfPages   string `json:"numberOfPages"`
	MediaURL        string `json:"mediaUrl"`
	Error           string `json:"error"`
	QRCode          string `json:"qrcode"`
}

func receiverSubscription() {
	bgCtx := context.Background()
	psc, err := pubsub.NewClient(bgCtx, os.Getenv("GCP_PROJECTID"))
	if err != nil {
		log.Fatal("unable to create the pubsub client: ", err)
	}

	topic := psc.Topic("fax_received")
	if err != nil {
		log.Fatal("error returned by CreateTopic: ", err)
	}

	sub := psc.Subscription("bg-faxreceiver")
	ok, err := sub.Exists(bgCtx)
	if err != nil {
		log.Fatal("unable to create the pubsub subscription: ", err)
	} else if !ok {
		config := pubsub.SubscriptionConfig{
			Topic: topic,
		}
		sub, err = psc.CreateSubscription(bgCtx, "bg-faxreceiver", config)
		if err != nil {
			log.Fatal("error getting pubsub subscription: ", err)
		}
	}

	fmt.Println("message queue receiver subscription established.")
	err = sub.Receive(bgCtx, func(ctx context.Context, m *pubsub.Message) {
		fmt.Println("received a new fax")
		var fax InboundFax
		if err := json.Unmarshal(m.Data, &fax); err != nil {
			log.Println("error while decoding the fax msg: ", err)
			m.Nack()
			return
		}

		go receiveFax(fax)
		m.Ack()
	})

	if err != nil {
		log.Println("error establishing the Receive handler: ", err)
	}
}

func receiveFax(fax InboundFax) {
	ib, err := json.Marshal(fax)
	if err != nil {
		log.Println("unable to marshal received fax: ", err)
	}

	// ... you have the ib.QRCode to determine what to do next
}
```

In my case I start this in my `main` function as a goroutine:

```go
func main() {
	//...
	go receiverSubscription()
	// ...
}
```

Here you have it. A fully working fax machine that can scale with an automatic 
dispatch of document based on a QR-Code value.

## Cost comparison of on-premise fax machines and Twilio

It's been almost two months that this solution is in production, so I don't 
have a huge dataset to base some conclusions.

Twilio's cost is higher since it's based on volume. Here are some aspects that 
differ from paperless compared to the printed faxes:

* Locally there were eight fax lines. There's one Twilio number.
* Locally they were printing all the documents. Now, nothing is printed. This 
is a significant environmental gain I'd say and a cost-saving for printers and 
paper.
* Locally they were manually handling the dispatch of the document. With the 
QR-Code, it's automated.
* Locally there was still some small maintenance, although fax machine tends to 
be plug and forget kind of devices. But the printers. Now there's nothing to 
maintain in terms of physical devices. On the other hand, there's more code 
to maintain.

Overall it will take one or two years before a real cost comparison can be 
made, for now, it feels higher due to the monthly cost.

The significant benefit here was that it allows all employee to continue their 
work from home, and the company can operate almost normally.