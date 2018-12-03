# Updating Fathom to the latest version

To update your existing Fathom installation to the latest version, first rename your existing Fathom installation so that we can move the new version in its place.

```
mv /usr/local/bin/fathom /usr/local/bin/fathom-old
```

Then, [download the latest release archive suitable for your system architecture from the releases page](https://github.com/usefathom/fathom/releases/latest) and place it in `/usr/local/bin`.

```
tar -C /usr/local/bin -xzf fathom_$VERSION_$OS_$ARCH.tar.gz
chmod +x /usr/local/bin/fathom
``` 

If you now run `fathom --version`, you should see that your system is running the latest version. 

```
$ fathom --version
Fathom version 1.0.0
```


### Restarting your Fathom web server

To start serving up the updated Fathom web application, you will have to restart the Fathom process that is running the web server.

If you've followed the [installation instructions](Installation%20instructions.md) then you are using Systemd to manage the Fathom process. Run `systemctl restart <your-fathom-service>` to restart it.

```
systemctl restart my-fathom-site
```

Alternatively, kill all running Fathom process by issuing the following command.

```
pkill fathom
```
