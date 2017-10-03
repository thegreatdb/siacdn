import React from 'react';
import Head from 'next/head';
import Link from 'next/link';
import Router from 'next/router';
import cookies from 'next-cookies';
import HttpsRedirect from 'react-https-redirect';
import { Segment, Header, Button, Form, Message } from 'semantic-ui-react';
import Client from '../lib/client';
import redirect from '../lib/redirect';
import {
  StripeProvider,
  Elements,
  CardElement,
  injectStripe,
} from 'react-stripe-elements';
import Nav from '../components/nav';

const IS_SERVER = typeof window === 'undefined';

class SignupForm extends React.Component {
  state = { error: null, submitting: false };

  static async getInitialProps(ctx) {
    const { authTokenID } = cookies(ctx);
    return { authTokenID };
  }

  handleSubmit = async (ev, err) => {
    ev.preventDefault();
    if (this.password1.value != this.password2.value) {
      await this.setState({ error: { message: 'Passwords must match' } });
      return;
    }
    await this.setState({ error: null, submitting: true });
    try {
      const { token } = await this.props.stripe.createToken({
        type: 'card',
        name: this.name.value,
      });
      const { authTokenID } = this.props;
      const client = new Client(authTokenID);
      const account = await client.createAccount(
        this.email.value,
        this.password1.value,
        this.name.value,
        token
      );
      // TODO: Use the account for something here?
      await Router.push('/dashboard');
      // Updating state directly because component is unmounted
      this.state['authAccount'] = account;
      this.state['submitting'] = false;
    } catch (error) {
      await this.setState({ error, submitting: false });
    }
  };

  render() {
    const { submitting, error } = this.state;
    const hasError = Boolean(error);
    return (
      <Form error={hasError} loading={submitting} onSubmit={this.handleSubmit}>
        {hasError ? (
          <Message header="Whoops!" content={'' + error.message} error />
        ) : null}
        <Form.Field>
          <label>E-Mail</label>
          <input
            placeholder="E-Mail"
            type="email"
            ref={e => (this.email = e)}
          />
        </Form.Field>
        <Form.Field>
          <label>Password</label>
          <input
            type="password"
            placeholder="Password"
            ref={e => (this.password1 = e)}
          />
        </Form.Field>
        <Form.Field>
          <label>Password (Repeat)</label>
          <input
            type="password"
            placeholder="Password (Repeat)"
            ref={e => (this.password2 = e)}
          />
        </Form.Field>
        <Form.Field>
          <label>First and last name</label>
          <input placeholder="First and last name" ref={e => (this.name = e)} />
        </Form.Field>
        <Form.Field>
          <label>
            Card details{' '}
            <span className="hint">(this will not initiate a charge)</span>
          </label>
          <div className="fieldWrapper">
            {IS_SERVER ? null : (
              <CardElement
                style={{
                  base: {
                    fontSize: '16px',
                    fontFamily:
                      "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif",
                    lineHeight: '24px',
                  },
                }}
              />
            )}
          </div>
        </Form.Field>
        <Button type="submit">Sign up</Button>
        <style jsx>{`
          .fieldWrapper {
            padding: 6px;
            border: 1px solid rgba(34, 36, 38, 0.15);
            border-radius: 0.28571429rem;
          }
          .hint {
            color: rgba(100, 100, 100, 0.4) !important;
            font-weight: normal;
          }
        `}</style>
      </Form>
    );
  }
}

if (!IS_SERVER) {
  const TempSignupForm = injectStripe(SignupForm);
  SignupForm = () => (
    <Elements>
      <TempSignupForm />
    </Elements>
  );
}

const render = () => (
  <HttpsRedirect>
    <Head>
      <link
        rel="stylesheet"
        href="//cdnjs.cloudflare.com/ajax/libs/semantic-ui/2.2.2/semantic.min.css"
      />
      <link rel="stylesheet" href="/static/css/global.css" />
      <script src="https://js.stripe.com/v3/" />
    </Head>
    <div className="holder">
      <Nav activeItem="signup" authAccount={null} />
      <Segment padded>
        <Header as="h1">Sign up</Header>
        <SignupForm />
      </Segment>
    </div>
  </HttpsRedirect>
);

export default () =>
  IS_SERVER ? (
    render()
  ) : (
    <StripeProvider apiKey="pk_live_cldTU8d2mloPzbJhYzb8dQF2">
      {render()}
    </StripeProvider>
  );
