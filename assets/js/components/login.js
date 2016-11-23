import m from 'mithril';

function handleSubmit(e) {
  e.preventDefault();

  fetch('/api/session', {
    method: "POST",
    data: {
      email: this.data.email(),
      password: this.data.password()
    },
    credentials: 'include'
  }).then((r) => {
    if( r.status == 200 ) {
      this.onAuth();
      console.log("Authenticated!");
    }

    // TODO: Handle errors
  });
}

const Login = {
  controller(args) {
    this.onAuth = args.onAuth;
    this.data = {
      email: m.prop(''),
      password: m.prop(''),
    }
    this.onSubmit = handleSubmit.bind(this);
  },

  view(c) {
    return m('div.block', [
      m('h2', 'Login'),
      m('p', 'Please enter your credentials to login to your Ana dashboard.'),
      m('form', {
        method: "POST",
        onsubmit: c.onSubmit
      }, [
        m('div.form-group', [
          m('label', 'Email address'),
          m('input', {
            type: "email",
            name: "email",
            required: true,
            onchange: m.withAttr("value", c.data.email )
          }),
        ]),
        m('div.form-group', [
          m('label', 'Password'),
          m('input', {
            type: "password",
            name: "password",
            required: true,
            onchange: m.withAttr("value", c.data.password )
          }),
        ]),
        m('div.form-group', [
          m('input', {
            type: "submit",
            value: "Sign in"
          }),
        ]),
      ])
    ])
  }
}

export default Login
