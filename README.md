Fathom Lite - simple website analytics
==============================
[![Go Report Card](https://goreportcard.com/badge/github.com/usefathom/fathom)](https://goreportcard.com/report/github.com/usefathom/fathom)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/usefathom/fathom/master/LICENSE)

Fathom Lite is a previous and open-source version of [Fathom Analytics](https://usefathom.com) (a paid, hosted [Google Analytics alternative](https://usefathom.com/google-analytics-alternative)). It was the very first version of our software, and has been downloaded millions of times!

While we are no longer adding features to this Lite version, we will be continuing to maintain it long-term and fix any bugs that come up.

![Screenshot of the Fathom dashboard](https://github.com/usefathom/fathom/raw/master/assets/src/img/fathom.jpg?v=7)

## Fathom Lite vs Fathom Analytics
Fathom Analytics is much more detailed, feature-rich, and even more focused on [privacy-law compliance](https://usefathom.com/compliance), than Fathom Lite. 

If you’d like to become a customer of Fathom Analytics, and not have to worry about servers, maintenance, security, you can give our software a try with a [30-day free trial](https://usefathom.com/ref/GITHUB) **(this link will give you $10 credit).**

![Screenshot of the Fathom Analytics Dashboard](https://usefathom.com/assets/images/fathom-screenshot.png)

| Lite | Pro |
|-----------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------|
| Uses Cookies| Cookie-free|
|-|[EU Isolation](https://usefathom.com/features/eu-isolation)|
|-|[Bypass ad-blockers](https://usefathom.com/features/custom-domains)|
|-|[Email reports](https://usefathom.com/docs/features/email-reports)|
|-|[Track event completions](https://usefathom.com/docs/features/events)|
|-|[Share your dashboard publicly or privately](https://usefathom.com/docs/features/shared-dashboards)|
|-|[Track UTM campaigns](https://usefathom.com/docs/features/campaigns)|
|-|Device, Browser & Country Data|
|No support|Fast and responsive support from the Fathom founders|
|No guaranteed uptime|Fully redundant cloud-based architecture that keeps your analytics online|
|Scaling requires you to power down and then upgrade your server|On-demand, automatic scaling, so even if your site goes viral, your stats won’t stop or slow down - we handle billions of page views each month|
|Manual backups by you|Real time backups included in the cost|
|No data protection|Continuous data protection (in the event of a database hardware failure, we have a live standby database ready to failover to)|
|Manual updates, via the command line|Automatic updates, patches and new versions at no extra cost, no coding required|
|Manual server hardening for security|Totally secure servers, monitored and maintained by us, included in the price|
|You pay for hosting, you do all the work to maintain the server, the code and backups|You pay us, we take care of everything for you|
|Embed script served via single server, from a single geographical location|Embed script served via our super-fast CDN, with endpoints located around the world to ensure fast page loads|
|Data aggregation performed on a single server|Super fast data-aggregation spread across our cloud architecture|
|Fork this repo|Supporting a privacy-focused, indie software company|
|Offered as-is|New features added all the time|
|-|[Get started for free](https://app.usefathom.com/register)|


## Installation


### Production

You can install Fathom on your server by following [our simple instructions](docs/Installation%20instructions.md).

### Development

For getting a development version of Fathom up & running, go through the following steps.

1. Ensure you have [Go](https://golang.org/doc/install#install) and [NPM](https://www.npmjs.com) installed
1. Download the code: `git clone https://github.com/usefathom/fathom.git $GOPATH/src/github.com/usefathom/fathom`
1. Compile the project into an executable: `make build`
1. (Optional) Set [custom configuration values](docs/Configuration.md)
1. (Required) Register a user account: `./fathom user add --email=<email> --password=<password>`
1. Start the webserver: `./fathom server` and then visit **http://localhost:8080** to access your analytics dashboard

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

## Copyright and license

MIT licensed. Fathom and Fathom logo are trademarks of Fathom Analytics.
