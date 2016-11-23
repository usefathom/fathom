'use strict';

import m from 'mithril';

function handleSubmit(e) {
  e.preventDefault();

  fetch('/api/session', {
    method: "DELETE",
    credentials: 'include',
  }).then((r) => {
    if( r.status == 200 ) {
      this.cb();
      console.log("No longer authenticated!");
    }
  });
}

const LogoutButton = {
  controller(args) {
    this.cb = args.cb;
    this.onSubmit = handleSubmit.bind(this);
  },

  view(c) {
    return m('a', { href: "#", onclick: c.onSubmit }, 'Sign out');
  }
}

export default LogoutButton
