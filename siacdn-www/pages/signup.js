import React from 'react';
import Head from 'next/head';
import Link from 'next/link';
import Router from 'next/router';
import cookies from 'next-cookies';
import {
  Segment,
  Header,
  Button,
  Form,
} from 'semantic-ui-react';
import Client from '../lib/client';
import redirect from '../lib/redirect';
import {StripeProvider, Elements, CardElement, injectStripe} from 'react-stripe-elements';
import Nav from '../components/nav';

const IS_SERVER = typeof window === 'undefined';

class SignupForm extends React.Component {
  static async getInitialProps(ctx) {
    const { authTokenID } = cookies(ctx);
    return { authTokenID };
  }

  handleSubmit = async (ev, err) => {
    ev.preventDefault();
    if (this.password1.value != this.password2.value) {
      alert('Passwords must match')
      return;
    }
    const { token } = await this.props.stripe.createToken({type: 'card', name: 'Jenny Rosen'});
    const { authTokenID } = this.props;
    const client = new Client(authTokenID);
    const account = await client.createAccount(this.username.value, this.password1.value, token)
    Router.push('/dashboard');
  }

  render() {
    return (
      <Form onSubmit={this.handleSubmit}>
        <Form.Field>
          <label>Username</label>
          <input placeholder='Username' ref={e => (this.username = e)} />
        </Form.Field>
        <Form.Field>
          <label>Password</label>
          <input type="password" placeholder='Password' ref={e => (this.password1 = e)} />
        </Form.Field>
        <Form.Field>
          <label>Password (Repeat)</label>
          <input type="password" placeholder='Password (Repeat)' ref={e => (this.password2 = e)} />
        </Form.Field>
        <Form.Field>
          <label>Card details <span className="hint">(this will not initiate a charge)</span></label>
          <div className="fieldWrapper">
            {IS_SERVER ? null : <CardElement style={{base: {fontSize: '16px', fontFamily: "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif", lineHeight: '24px'}}} />}
          </div>
        </Form.Field>
        <Button type='submit'>Sign up</Button>
        <style jsx>{`
          .fieldWrapper {
            padding: 6px;
            border: 1px solid rgba(34,36,38,.15);
            border-radius: .28571429rem;
          }
          .hint {
            color: rgba(100,100,100,.4) !important;
            font-weight: normal;
          }
        `}</style>
      </Form>
    )
  }
}

if (!IS_SERVER) {
  const TempSignupForm = injectStripe(SignupForm);
  SignupForm = () => (
    <Elements><TempSignupForm /></Elements>
  )
}

const render = () => (
  <div>
    <Head>
      <link
        rel="stylesheet"
        href="//cdnjs.cloudflare.com/ajax/libs/semantic-ui/2.2.2/semantic.min.css"
      />
      <link rel="stylesheet" href="/static/css/global.css" />
      <script src="https://js.stripe.com/v3/"></script>
    </Head>
    <div className="holder">
      <Nav activeItem="signup" authAccount={null} />
      <Segment padded>
        <Header as="h1">Sign up</Header>
        <SignupForm />
      </Segment>
    </div>
  </div>
);

export default () => (
  IS_SERVER ? render() : <StripeProvider apiKey="pk_test_zMFraFeAYdlJGMqNzSq1Bw5o">
    {render()}
  </StripeProvider>
);
