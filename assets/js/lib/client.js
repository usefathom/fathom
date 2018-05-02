'use strict';

var Client = {};
Client.request = function(resource, args) {
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

  return fetch(`/api/${resource}`, args)
    .then(handleRequestErrors)
    .then(parseJSON)
    .then(parseData)
}

function parseJSON(r) {
  return r.json()
}

function handleRequestErrors(r) {
    if (!r.ok) {
        throw new Error(r.statusText);
    }
    return r;
}

function parseData(d) {
  if(d.Error) {
    throw new Error(d.Error)
  }

  return d.Data
}

export default Client
