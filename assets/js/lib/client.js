'use strict';

var Client = {};

Client.request = function(resource, args) {
  args = args || {};
  args.credentials = 'include'

  return fetch(`/api/${resource}`, args).then((r) => {
    if( r.ok ) {
      return r.json();
    }

    throw new Error(r);
  })
}

export default Client
