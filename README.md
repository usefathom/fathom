Fathom Lite - simple website analytics
==============================

[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=LJ5WZVA9ER9GJ)
[![Go Report Card](https://goreportcard.com/badge/github.com/usefathom/fathom)](https://goreportcard.com/report/github.com/usefathom/fathom)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/usefathom/fathom/master/LICENSE)

[Fathom Analytics](https://usefathom.com/) is a simpler and more privacy-focused alternative to Google Analytics.

Collecting information on the internet is important, but it’s broken. We’ve become complacent in trading information for free access to web services, and then complaining when those web services do crappy things with that data.

The problem is this: _if we aren’t paying for the product, we are the product_.

Google Analytics may give you free access to their services but in turn, they’re assembling data profiles on your website visitors, which they can then use for better targeting of advertisements across their network.

We need to stop giving away our data and our users' privacy for free access to a tool.

Fathom respects the privacy of your users and does not collect any personally identifiable information. All while giving you the information you need about your site, so you can make smarter decisions about your design and content.

At present, Fathom Analytics Lite is not PECR compliant due to the fact that it uses an anonymous cookie. Our [PRO version](https://usefathom.com) is PECR compliant, and we'll be making changes to this codebase some time in the future to make it compliant.

![Screenshot of the Fathom dashboard](https://github.com/usefathom/fathom/raw/master/assets/src/img/fathom.jpg?v=7)

## Lite vs PRO

We offer a [PRO version of Fathom Analytics](https://usefathom.com/#pricing) that starts at $14 / month. If you’d like to become a customer, we’d love to have you on board. You can signup for a [free 7-day trial of Fathom Analytics here](https://app.usefathom.com/register).

![Screenshot of the PRO Dashboard](https://usefathom.com/assets/fathom-analytics.png)

| Lite | Pro |
|-----------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------|
| Uses Cookies| Cookie-free|
|-|Track goal completions|
|-|Share your site publicly or privately|
|-|Device, Browser & Country Data|
| No support (you can post an issue to our repo) | Fast and responsive support from the Fathom founders|
| No guaranteed uptime | Fully redundant cloud-based architecture that keeps your stats online                                                                               |
| Scaling requires you to power down and then upgrade your server | On-demand, automatic scaling, so even if your site goes viral, your stats won’t stop or slow down - we can handle billions of page views each month |
| Manual backups, if you know how to set them up on your server | Real time backups included in the cost|
| No data protection| Continuous data protection (in the event of a database hardware failure, we have a live standby database ready to go in another availability zone)|
| Manual updates, via the command line| Automatic updates, patches and new versions at no extra cost, no coding required| 
| Manual server hardening for security | Totally secure server, monitored and maintained by us, included in the price |
| You pay for hosting, you have to do all the work to maintain the server, the code and backups | You pay us, we take care of everything for you|
| Tracker file served via single server, from a single location | Tracker file served via our super-fast CDN, with endpoints located around the world to ensure fast page loads |
| Data aggregation performed on a single server | Super fast data-aggregation spread across our cloud architecture |
| Contribute to our repo| Supporting a privacy-focused, indie software company |    
|| [Get started for free](https://app.usefathom.com/register) |


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

## Roadmap

Find [our public roadmap here](https://trello.com/b/x2aBwH2J/fathom-roadmap). 

If you have an idea or suggestion for Fathom, [submit it as an issue here on GitHub](https://github.com/usefathom/fathom/issues).

## Donation
If you're self-hosting Fathom Analytics Lite and want to support it's development, you can:

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=LJ5WZVA9ER9GJ)

## Copyright and license

MIT licensed. Fathom and Fathom logo are trademarks of Fathom Analytics.
