# Install Fathom on Render

[Render](https://render.com) is a new cloud provider that makes it incredibly easy and quick to host Docker apps, with fully-managed SSL, custom domains and continuous deploys from GitHub.

Fathom Community Edition can be installed on Render in just a few minutes using the steps below.

> The canonical version of this guide lives at https://render.com/docs/deploy-fathom-analytics.


1. Create a new [PostgreSQL database](https://render.com/docs/databases) on Render. Set the **Name**, **Database**, and **User** to `fathom`.

2. Fork [render-examples/fathom-analytics](https://github.com/render-examples/fathom-analytics) and create a new **Web Service** on Render, giving Render's GitHub app permission to access your forked repo.

3. On the service creation page, click on **Advanced** and add a new **secret file** with filename **`fathom.env`** and the following content:

   > Make sure to update the values in the highlighted lines below.

   ```shell{6-8}
   FATHOM_GZIP=true
   FATHOM_DEBUG=false
   FATHOM_DATABASE_DRIVER="postgres"
   FATHOM_DATABASE_NAME="fathom"
   FATHOM_DATABASE_USER="fathom"
   FATHOM_DATABASE_PASSWORD="db password from step 1"
   FATHOM_DATABASE_HOST="internal db hostname from step 1"
   FATHOM_SECRET="a sufficiently strong secret"
   ```

   See [Configuration](https://github.com/usefathom/fathom/blob/master/docs/Configuration.md) for more details.

   Click on **Save web service** and Fathom will be available on your `onrender.com` URL in less than a minute.

4. Once Fathom is built, go to the **Shell** tab in your Render dashboard and create a Fathom admin user by entering this command:

   > Make sure to replace the email and password.

   ```shell
   ./fathom --config /etc/secrets/fathom.env user add --email="you@your-email.com" --password="strong-password"
   ```

That's it! You can now create a new site by visiting your Fathom `onrender.com` URL. Fathom will start displaying analytics as soon as you copy the tracking snippet to your website.

You can also add a [custom domain](https://rendre.com/docs/custom-domains) to your service and Render will automatically issue and manage SSL certificates for your domain.
