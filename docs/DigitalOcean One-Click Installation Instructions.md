# Fathom Analytics One-Click Application

> [Deploy your own Digital Ocean Droplet with our One-Click installer](https://marketplace.digitalocean.com/apps/fathom-analytics?action=deploy&refcode=a36c0e8abd93). New Digital Ocean customers will receive $100 for 60 days.

[Fathom Analytics](https://usefathom.com) is a powerful but simple way to track visitors and referrals to your website without compromising their privacy. No personal data is stored, stats are collected in aggregate, making Fathom GDPR and ePrivacy compliant.

The Fathom One-Click application provides a simple way to get started with the software by automatically installing it plus the required software on a single Ubuntu 18.04 Droplet.

> NOTE: A fully-managed, automatically updated and hosted service is also available via [Fathom PRO](https://usefathom.com), which lets you configure and manage everything from your browser.

### Quick Start

> [Deploy your own Digital Ocean Droplet with our One-Click installer](https://marketplace.digitalocean.com/apps/fathom-analytics?action=deploy&refcode=a36c0e8abd93). New Digital Ocean customers will receive $100 for 60 days.

We recommend you have a registered domain name already, and that you use a sub-domain for Fathom Analytics. You’ll need to setup an A record from the domain to point to the IP address of your Fathom Analytics One-Click Droplet, (e.g. `stats.example.com`).

After the Fathom Analytics One-Click Droplet has been created and the DNS records are setup, you’ll need to log into your new Fathom Droplet to finish the setup.

From a terminal on your local computer, connect to the Droplet as the `root` user. Make sure to substitute the IP address of the Droplet:

```ssh root@use_your_droplet_ip```

As soon as you log in, the Droplet will automatically prompt you through the setup of Fathom Analytics.

```Welcome to the Fathom Analytics setup process!```

Follow the prompts to configure Fathom Analytics.

1. You’ll type `yes` to `Will you be pointing a domain at this Fathom instance?`
2. You’ll enter your domain name, e.g. `stats.example.com`
3. You’ll type `yes` to `Do you want to password protect your Fathom stats?`
4. You’ll type an email address (this is your username) for your Fathom account
5. You’ll type and then confirm a password for your Fathom account

Once that’s finished, you’ll see a message along the lines of:

```Once your point your domain to this server (XX.XX.XX.XX), you can access Fathom at: https://stats.example.com```

Where `XX.XX.XX.XX` is the IP address to your Droplet and `https://stats.example.com` is the that you pointed to it.

### Secure your Fathom Dashboard with SSL

Now that Fathom Analytics is installed on your server, you can access it with your email address and password. But, it won’t be fully secure as it uses `http://` not `https://`

To secure your dashboard, you can run: 

```certbot --nginx -d stats.example.com```

And follow the prompts to install a Let’s Encrypt free SSL certificate on your server. Make sure you select that all requests are loaded via `https://`.

### Setup your Fathom Dashboard

Now that your dashboard is secure, you can log into it at:

```https://stats.example.com```

Where the URL is the one you pointed to your Droplet via an A record.

You’ll log in with your email address and password, and give your dashboard a `Site name`. The Site name is just for your reference (in case you add multiple dashboards to the same Fathom instance).

Next you’ll be given your code snippet, with your unique `site id` to add to the footer of any/all pages on the website you’d like Fathom to track.

### Next steps

Log into your Fathom dashboard every few days to see the most popular content and biggest referrers to your website. You can also add infinite sites to track from the same dashboard.

> [Deploy your own Digital Ocean Droplet with our One-Click installer](https://marketplace.digitalocean.com/apps/fathom-analytics?action=deploy&refcode=a36c0e8abd93). New Digital Ocean customers will receive $100 for 60 days.
