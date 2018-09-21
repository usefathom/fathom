'use strict';

var Client = {};
Client.request = function(url, args) {
  args = args || {};
  args.credentials = 'same-origin'
  args.headers = args.headers || {};
  args.headers['Accept'] = 'application/json';

  if( args.method && args.method === 'POST') {
    args.headers['Content-Type'] = 'application/json';

    if(args.data) {
      if( typeof(args.data) !== "string") {
        args.data = JSON.stringify(args.data)
      }
      args.body = args.data
      delete args.data
    }
  }

  // trim leading slash from URL
  url = (url[0] === '/') ? url.substring(1) : url;

  return window.fetch(`/api/${url}`, args)
    .then(handleRequestErrors)
    .then(parseJSON)
    .then(parseData)
}

function handleRequestErrors(r) {
  // if response is not JSON (eg timeout), throw a generic error
  if (! r.ok && r.headers.get("Content-Type") !== "application/json") {
    throw { code: "request_error", message: "An error occurred" }
  }

  return r
}

function parseJSON(r) {
  return r.json()
}

function parseData(d) {

  // if JSON response contains an Error property, use that as error code
  // Message is generic here, so that individual components can set their own specific messages based on the error code
  if(d.Error) {
    throw { code: d.Error, message: "An error occurred" }
  }

  return d.Data
}

export default Client
