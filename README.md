Fathom - simple website analytics
==============================

[![Go Report Card](https://goreportcard.com/badge/github.com/usefathom/fathom)](https://goreportcard.com/report/github.com/usefathom/fathom)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/usefathom/fathom/master/LICENSE)

[Fathom Analytics](https://usefathom.com/) is a simpler and more privacy-focused alternative to Google Analytics.

Collecting information on the internet is important, but it’s broken. We’ve become complacent in trading information for free access to web services, and then complaining when those web services do crappy things with that data.

The problem is this: _if we aren’t paying for the product, we are the product_.

Google Analytics may give you free access to their services but in turn, they’re assembling data profiles on your website visitors, which they can then use for better targeting of advertisements across their network.

We need to stop giving away our data and our users' privacy for free access to a tool.

Fathom [respects the privacy of your users and does not collect any personally identifiable information](https://usefathom.com/data/). All while giving you the information you need about your site, so you can make smarter decisions about your design and content.

![Screenshot of the Fathom dashboard](https://github.com/usefathom/fathom/raw/master/assets/src/img/fathom.jpg?v=7)

## Installation

For getting a development version of Fathom up & running, go through the following steps.

1. Ensure you have [Go](https://golang.org/doc/install#install) and [NPM](https://www.npmjs.com) installed
1. Download the code: `git clone https://github.com/usefathom/fathom.git $GOPATH/src/github.com/usefathom/fathom` 
1. Compile the project: `make build` 
1. (Optional) Set [custom configuration values](docs/Configuration.md)
1. (Optional) Register a user account: `fathom user add --email=<email> --password=<password>`
1. Start the webserver: `fathom server` and then visit **http://localhost:8080** to access your analytics dashboard

To install and run Fathom in production, [see the installation instructions](docs/Installation%20instructions.md).

## Docker

### Building

Ensure you have Docker installed and run `docker build -t fathom .`.
Run the container with `docker run -d -p 8080:8080 fathom`.

### Running

To run [our pre-built Docker image](https://hub.docker.com/r/usefathom/fathom/), run `docker run -d -p 8080:8080 usefathom/fathom:latest`

## Tracking snippet

To start tracking, create a site in your Fathom dashboard and copy the tracking snippet to the website(s) you want to track.

### Content Security Policy

If you use a [Content Security Policy (CSP)](https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP) to specify security policies for your website, Fathom requires the following CSP directives (replace `yourfathom.com` with the URL to your Fathom instance):

```
script-src: yourfathom.com;
img-src: yourfathom.com;
```

## Roadmap

Find [our public roadmap here](https://trello.com/b/x2aBwH2J/fathom-roadmap). 

If you have an idea or suggestion for Fathom, [submit it as an issue here on GitHub](https://github.com/usefathom/fathom/issues).

## Copyright and license

MIT licensed. Fathom and Fathom logo are trademarks of Fathom Analytics.
