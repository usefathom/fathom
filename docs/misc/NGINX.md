# Using NGINX with Fathom

Let's say you have the Fathom server listening on port 9000 and want to serve it on your domain, `yourfathom.com`.

We can use NGINX to redirect all traffic for a certain domain to our Fathom application by using the `proxy_pass` directive combined with the port Fathom is listening on. 

Create the following file in `/etc/nginx/sites-enabled/yourfathom.com`

```
server {
	server_name yourfathom.com;

	location / {
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header X-Forwarded-For $remote_addr;
		proxy_set_header Host $host;
		proxy_pass http://127.0.0.1:9000; 
	}
}
```

### Test NGINX configuration
```
sudo nginx -t
```

### Reload NGINX configuration

```
sudo service nginx reload
```
