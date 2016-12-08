'use strict';

var Client = {};

Client.request = function(resource, args) {
  args = args || {};
  args.credentials = 'include'
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

  return fetch(`/api/${resource}`, args).then((r) => {
    if( r.ok ) {
      return r.json();
    }

    throw new Error(r);
  })
}

export default Client
