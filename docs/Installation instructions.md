# Installation instructions for Fathom

To install Fathom on your server: 

1. [Download the latest Fathom release](https://github.com/usefathom/fathom/releases) suitable for your platform.
2. Extract the archive to `/usr/local/bin`

```sh
tar -C /usr/local/bin -xzf fathom_$VERSION_$OS_$ARCH.tar.gz
chmod +x /usr/local/bin/fathom
```

Confirm that Fathom is installed properly by running `fathom --version`

```sh
$ fathom --version
Fathom version 1.0.0
```

## Configuring Fathom

> This step is optional. By default, Fathom will use a SQLite database file in the current working directory.

To run the Fathom web server we will need to [configure Fathom](Configuration.md) so that it can connect with your database of choice. 

Let's create a new directory where we can store our configuration file & SQLite database.

```
mkdir ~/my-fathom-site
cd ~/my-fathom-site
```

Then, create a file named `.env` with the following contents.

```
FATHOM_SERVER_ADDR=9000
FATHOM_GZIP=true
FATHOM_DEBUG=true
FATHOM_DATABASE_DRIVER="sqlite3"
FATHOM_DATABASE_NAME="fathom.db"
FATHOM_SECRET="random-secret-string"
```

If you now run `fathom server` then Fathom will start serving up a website on port 9000 using a SQLite database file named `fathom.db`. If that port is exposed then you should now see your Fathom instance running by browsing to `http://server-ip-address-here:9000`.

Check out the [configuration file documentation](Configuration.md) for all possible configuration values, eg if you want to use MySQL or Postgres instead.

## Register your admin user

> This step is required.

To register a user in the Fathom instance we just created, run the following command from the directory where your `.env` file is. 

```
fathom user add --email="john@email.com" --password="strong-password"
```

**Note:** if you're running Fathom v1.0.1 or older, the command is `fathom register --email="john@email.com" --password="strong-password"`

## Using NGINX with Fathom

We recommend using NGINX with Fathom, as it simplifies running multiple sites from the same server and handling SSL certificates with LetsEncrypt.

Create a new file in `/etc/nginx/sites-enabled/my-fathom-site` with the following contents. Replace `my-fathom-site.com` with the domain you would like to use for accessing your Fathom installation.

```sh
server {
	server_name my-fathom-site.com;

	location / {
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header X-Forwarded-For $remote_addr;
		proxy_set_header Host $host;
		proxy_pass http://127.0.0.1:9000; 
	}
}
```

Test your NGINX configuration and reload NGINX.

```
nginx -t
service nginx reload
```

If you now run `fathom server` again, you should be able to access your Fathom installation by browsing to `http://my-fathom-site.com`.

## Automatically starting Fathom on boot

To ensure the Fathom web server keeps running whenever the system reboots, we should use a process manager. Ubuntu 16.04 and later ship with Systemd.

Create a new file called `/etc/systemd/system/my-fathom-site.service` with the following contents. Replace `$USER` with your actual username.

```
[Unit]
Description=Starts the fathom server
Requires=network.target
After=network.target

[Service]
Type=simple
User=$USER
Restart=always
RestartSec=3
WorkingDirectory=/home/$USER/my-fathom-site
ExecStart=/usr/local/bin/fathom server

[Install]
WantedBy=multi-user.target
```

Reload the Systemd configuration & enable our service so that Fathom is automatically started whenever the system boots.

```
systemctl daemon-reload
systemctl enable my-fathom-site
```

You should now be able to manually start your Fathom web server by issuing the following command.

```
systemctl start my-fathom-site
```

## Tracking snippet

To start tracking pageviews, copy the tracking snippet shown in your Fathom dashboard to all pages of the website you want to track.


### SSL certificate

With [Certbot](https://certbot.eff.org/docs/) for LetsEncrypt installed, adding an SSL certificate to your Fathom installation is as easy as running the following command.

```
certbot --nginx -d my-fathom-site.com
```


