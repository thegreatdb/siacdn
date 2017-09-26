import React from 'react';
import Head from 'next/head';
import Router from 'next/router';
import cookies from 'next-cookies';
import {
  Segment,
  Header,
  Button,
  Form,
} from 'semantic-ui-react';
import Client from '../lib/client';
import Nav from '../components/nav';

export default class LoginForm extends React.Component {
  state = { error: null, submitting: false, authAccount: null };

  async getInitialProps(ctx) {
    const { authTokenID } = cookies(ctx);
    return { authTokenID };
  }

  handleSubmit = async (ev) => {
    ev.preventDefault();
    const username = this.username.value;
    const password = this.password.value;
    await this.setState({ error: null, submitting: true });
    const { authTokenID } = this.props;
    const client = new Client(authTokenID);
    try {
      const account = await client.loginAccount(username, password);
      await Router.push('/dashboard');
      // Updating state directly because component is unmounted
      this.state['authAccount'] = account;
      this.state['submitting'] = false;
    } catch (error) {
      await this.setState({ error, submitting: false });
    }
  }

  render() {
    return (
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
          <Nav activeItem="login" authAccount={null} />
          <Segment padded>
            <Header as="h1">Log in</Header>
            <Form onSubmit={this.handleSubmit}>
              <Form.Field>
                <label>Username</label>
                <input placeholder='Username' ref={e => (this.username = e)} />
              </Form.Field>
              <Form.Field>
                <label>Password</label>
                <input type="password" placeholder='Password' ref={e => (this.password = e)} />
              </Form.Field>
              <Button type='submit'>Log in</Button>
            </Form>
          </Segment>
        </div>
      </div>
    )
  }
}
