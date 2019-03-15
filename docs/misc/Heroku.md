# Running Fathom on Heroku

### Requirements

* heroku cli (logged in)
* git
* curl
* wget
* tar are required
* ~ openssl is required to generate the secret_key, but you're free to use what you want

### Create the app

First you need to choose a unique app name, as Heroku generates a subdomain for your app.

* create the app via the buildpack

```bash
heroku create UNIQUE_APP_NAME --buildpack https://github.com/ph3nx/heroku-binary-buildpack.git
```

* locally clone the newly created app

```bash
heroku git:clone -a UNIQUE_APP_NAME
cd UNIQUE_APP_NAME
```

* create the folder that will contain fathom

```bash
mkdir -p bin
```

* download latest version of fathom for linux 64bit

```bash
curl -s https://api.github.com/repos/usefathom/fathom/releases/latest \
  | grep browser_download_url \
  | grep linux_amd64.tar.gz \
  | cut -d '"' -f 4 \
  | wget -qi - -O- \
  | tar --directory bin -xz - fathom
```

* create the Procfile for Heroku

```bash
echo "web: bin/fathom server" > Procfile
```

* create a Postgres database (you can change the type of plan if you want - https://elements.heroku.com/addons/heroku-postgresql#pricing)

```bash
heroku addons:create heroku-postgresql:hobby-dev
```

* update the environment variables, generate a secret_key

here you can change the way you generate your secret_key.

```bash
heroku config:set PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/app/bin \
  FATHOM_DATABASE_DRIVER=postgres \
  FATHOM_DATABASE_URL=$(heroku config:get DATABASE_URL) \
  FATHOM_DEBUG=true \
  FATHOM_SECRET= $(openssl rand -base64 32) \
  FATHOM_GZIP=true
```

* add, commit and push all our files

```bash
git add --all
git commit -m "First Commit"
git push heroku master
```

* the created app runs as a free-tier. A free-tier dyno uses the account-based pool
of free dyno hours. If you have other free dynos running, you will need to upgrade your app to a 'hobby' one. - https://www.heroku.com/pricing

```bash
heroku dyno:resize hobby
```

* check that everything is working

```bash
heroku run fathom --version
```

* add the first user

```bash
heroku run fathom user add --email="test@test.com" --password="test_password"
```

* open the browser to login and add your first website

```bash
heroku open
```

* ENJOY :)
