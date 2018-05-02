'use strict';

import { h, render, Component } from 'preact';
import LoginForm from '../components/LoginForm.js';

class Login extends Component {
  render() {
    return (
      <div class="wrapper">
         <header class="section">
           <nav class="main-nav animated fadeInDown">
               <ul>
                 <li class="logo"><a href="/">Fathom</a></li>
             </ul>
           </nav>
         </header>
         <section>
            <LoginForm onSuccess={this.props.onLogin}/>
         </section>
         <footer class="section"></footer>
       </div>
    )
  }
}

export default Login
