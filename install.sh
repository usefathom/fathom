#!/bin/bash
#
# Fathom Installer script
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
# 
# Defaults
SERVER_PORT="9000"   # Port for Fathom server

# Check if running as root
if [ "$EUID" -ne 0 ]; then
  echo "Please run the Fathom installer as root (to install the necessary NGINX & systemd files)"
  exit
fi

echo "Welcome to the Fathom quick installer. Press CTRL-C at any time to abort."

function download_fathom() {
   # Download latest version of the Fathom application
   echo "Downloading Fathom"
   wget -O fathom https://usesfathom.com/downloads/fathom-latest

   # Move Fathom to $PATH so we can run the command from anywhere
   chmod +x fathom
   mv fathom /usr/local/bin/fathom
   FATHOM_PATH=$(command -v fathom)
   echo "Fathom installed to $FATHOM_PATH" 
   echo ""
}

function new_site_dir() {
   read -p "Where would you like to store your new Fathom instance? (default: new-fathom): " SITE_DIR
   if [[ "$SITE_DIR" == "" ]]; then
      SITE_DIR="new-fathom"
   fi;

   SITE_DIR_ABS="$PWD/$SITE_DIR"
   if [ -d "$SITE_DIR" ]; then 
     read -p "Warning: $SITE_DIR_ABS already exists. Are you sure? (y/N): " CONTINUE
     if [ "$CONTINUE" != "y" ]; then exit 0; fi
   fi;
   
   if [ ! -e "$SITE_DIR" ]; then
      mkdir -p "$SITE_DIR"
      chmod 755 "$SITE_DIR"
   fi;
   
   cd "$SITE_DIR"
   echo ""
}

function setup_config() {
   # Ask for configuration values
   echo "Choose database engine:"
   echo "  1) SQLite (default)"

   if [ "$(command -v mysql)" ]; then
      echo "  2) MySQL"
   fi;

   read DATABASE_CHOICE

   if [ "$DATABASE_CHOICE" == "1" ] || [ "$DATABASE_CHOICE" == ""  ]; then
      DATABASE="sqlite3"
   fi;

   if [ "$DATABASE_CHOICE" == "2" ]; then
      DATABASE="mysql"
   fi;

   # set filename if using sqlite3
   if [ "$DATABASE" == "sqlite3" ]; then
      DATABASE_NAME="$SITE_DIR_ABS/fathom.db"
   fi;

   # or ask for credentials if using postgres or mysql
   if [ "$DATABASE" == "mysql" ] || [ "$DATABASE" == "postgres" ]; then
      echo "Enter your $DATABASE credentials: "
      read -p "  Database user: " DATABASE_USER
      read -p "  Database password: " DATABASE_PASSWORD
      read -p "  Database name (default: fathom): " DATABASE_NAME

      if [ "$DATABASE_NAME" == "" ]; then
         DATABASE_NAME="fathom"
      fi;

      echo "Creating database $DATABASE_NAME"
      mysql --user="$DATABASE_USER" --password="$DATABASE_PASSWORD" --execute="CREATE DATABASE $DATABASE_NAME;"
      # TODO: Add Postgres support 
   fi;

   # Create configuration file
   TEMPLATE=$(cat <<-END
   FATHOM_DEBUG=false
   FATHOM_SERVER_ADDR=":$SERVER_PORT"
   FATHOM_DATABASE_DRIVER="$DATABASE"
   FATHOM_DATABASE_NAME="$DATABASE_NAME"
   FATHOM_DATABASE_USER="$DATABASE_USER"
   FATHOM_DATABASE_PASSWORD="$DATABASE_PASSWORD"
   FATHOM_DATABASE_HOST=""
   FATHOM_SECRET="abcdefghijklmnopqrstuvwxyz1234567890"
END
)

   echo "$TEMPLATE" > ".env"
   echo "Created configuration file: $SITE_DIR_ABS/.env"
   echo "Success! You can now run Fathom using \`fathom --config=$SITE_DIR_ABS/.env server\`"
   echo ""
}

function new_fathom_user() {
   echo "Create user account: "
   read -p "  Email address: " USER_EMAIL
   read -p "  Password: " USER_PASSWORD
   fathom --config="$SITE_DIR_ABS/.env" register --email="$USER_EMAIL" --password="$USER_PASSWORD"
   echo ""
}

function new_nginx_server() {
   read -p "Server name (full domain): " SERVER_NAME
   if [ "$SERVER_NAME" == "" ]; then exit 0; fi

   # Make sure we're not overwriting existing server blocks
   if [ -e "/etc/nginx/sites-available/$SERVER_NAME" ]; then 
     read -p "Warning: /etc/nginx/sites-available/$SERVER_NAME already exists. Are you sure? (y/N): " CONTINUE
     if [ "$CONTINUE" != "y" ]; then exit 0; fi
   fi;

   TEMPLATE=$(cat <<-END
server {
   server_name $SERVER_NAME;

   location / {
      proxy_set_header X-Real-IP \$remote_addr;
      proxy_set_header X-Forwarded-For \$remote_addr;
      proxy_set_header Host \$host;
      proxy_pass http://127.0.0.1:$SERVER_PORT;
   }
}
END
)
   echo "$TEMPLATE" > "/etc/nginx/sites-available/$SERVER_NAME"
   ln -s "/etc/nginx/sites-available/$SERVER_NAME" "/etc/nginx/sites-enabled/$SERVER_NAME" || true
   nginx -t
   service nginx reload
   echo ""
}

function new_systemctl_service() {
   echo ""
   # Use server name or ask for new service name   
   SERVICE_NAME="$SERVER_NAME"   
   if [ "$SERVICE_NAME" == "" ]; then
      read -p "Service name: " SERVICE_NAME      
   fi;

   # Make sure we're not overwriting existing service files
   if [ -e "/etc/systemd/system/$SERVICE_NAME.service" ]; then 
     read -p "Warning: /etc/systemd/system/$SERVICE_NAME.service already exists. Are you sure? (y/N): " CONTINUE
     if [ "$CONTINUE" != "y" ]; then exit 0; fi
   fi;

    TEMPLATE=$(cat <<-END
[Unit]
Description=Fathom service for $SERVICE_NAME
Requires=network.target
After=network.target

[Service]
Type=simple
User=$USER
Restart=always
WorkingDirectory=$SITE_DIR_ABS
ExecStart=$FATHOM_PATH --config=$SITE_DIR_ABS/.env server
SyslogIdentifier=$SERVICE_NAME

[Install]
WantedBy=multi-user.target
END
)     
   
   echo "$TEMPLATE" > "/etc/systemd/system/$SERVICE_NAME.service"
   systemctl daemon-reload
   systemctl enable "$SERVICE_NAME"
   systemctl start "$SERVICE_NAME"

   echo "Success! Service $SERVICE_NAME created & started."
   echo "You can manually start the service using systemctl start $SERVICE_NAME"
   echo ""
}

function new_cerbot_cert() {
   certbot --nginx -d "$SERVER_NAME"
   echo ""
}

# Download Fathom if command does not exist
FATHOM_PATH="$(command -v fathom)"
if [ "$FATHOM_PATH" != "" ]; then
   echo "Fathom detected in $FATHOM_PATH"
   echo "Installed: $(fathom --version)"
   read -p "Download latest Fathom anyway? (Y/n): " CONTINUE
fi

if [ "$FATHOM_PATH" == "" ] || [ "$CONTINUE" != "n" ]; then
   download_fathom
fi;

new_site_dir
setup_config
new_fathom_user

# Ask to setup new NGINX server block
if [ "$(command -v nginx)" ]; then
   read -p "NGINX detected. Create a new server block? (Y/n): " CONTINUE
   if [ "$CONTINUE" != "n" ]; then
      new_nginx_server
   fi;
fi;

# Ask to setup new Systemctl service file
if [ "$(command -v systemctl)" ]; then
   read -p "Systemctl detected. Create a new service? (Y/n): " CONTINUE
   if [ "$CONTINUE" != "n" ]; then 
      new_systemctl_service
   fi
fi;

# Ask to request new LetsEncrypt certificate
if [ "$SERVER_NAME" != "" ] && [ "$(command -v certbot)" ]; then
   read -p "LetsEncrypt detected. Request SSL certificate for $SERVER_NAME? (Y/n): " CONTINUE
   if [ "$CONTINUE" != "n" ]; then 
      new_cerbot_cert
   fi
fi;

# Try to be helpful
if [ "$SERVER_NAME" != "" ] && [ "$SERVICE_NAME" != "" ]; then
   echo "Success! You should now see Fathom running on $SERVER_NAME"
   exit 0
fi;
