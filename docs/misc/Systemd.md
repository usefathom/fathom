# Managing the Fathom process with Systemd

To run Fathom as a service (so it keeps on running in the background and is automatically restarted in case of a server reboot) on Ubuntu 16.04 or later, first ensure you have the `fathom` binary installed and in your `$PATH` so that the command exists.

Then, create a new service config file in the `/etc/systemd/system/` directory.

Example file: `/etc/systemd/system/fathom.service`

The file should have the following contents, with `$USER` substituted with your actual username.

```
[Unit]
Description=Starts the fathom server
Requires=network.target
After=network.target

[Service]
Type=simple
User=$USER
Restart=always
RestartSec=6
WorkingDirectory=/etc/fathom # (or where fathom should store its files)
ExecStart=fathom server

[Install]
WantedBy=multi-user.target
```

Save the file and run `sudo systemctl daemon-reload` to load the changes from disk. 

Then, run `sudo systemctl enable fathom` to start the service whenever the system boots.

### Starting or stopping the Fathom service manually
```
sudo systemctl start fathom
sudo systemctl stop fathom
```

### Using a custom configuration file

If you want to [modify the configuration values for your Fathom service](https://github.com/usefathom/fathom/blob/master/docs/Configuration.md), then change the line starting with `ExecStart=...` to include the path to your configuration file.

For example, if you have a configuration file `/home/john/fathom.env` then the line should look like this:

```
ExecStart=fathom --config=/home/john/fathom.env server --addr=:9000
```

#### Start Fathom automatically at boot
```
sudo systemctl enable fathom
```

#### Stop Fathom from starting at boot

```
sudo systemctl disable fathom
```
