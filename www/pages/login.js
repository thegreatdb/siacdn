import React from 'react';
import Router from 'next/router';
import cookies from 'next-cookies';
import { Segment, Header, Button, Form, Message } from 'semantic-ui-react';
import Client from '../lib/client';
import Nav from '../components/nav';
import Footer from '../components/footer';
import PageHeader from '../components/pageheader';

export default class LoginForm extends React.Component {
  state = { error: null, submitting: false, authAccount: null };

  async getInitialProps(ctx) {
    const { authTokenID } = cookies(ctx);
    return { authTokenID };
  }

  handleSubmit = async ev => {
    ev.preventDefault();
    const email = this.email.value;
    const password = this.password.value;
    await this.setState({ error: null, submitting: true });
    const { authTokenID } = this.props;
    const client = new Client(authTokenID);
    try {
      const account = await client.loginAccount(email, password);
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
      <div>
        <PageHeader />
        <div className="holder">
          <Nav activeItem="login" authAccount={null} />
          <Segment padded>
            <Header as="h1">Log in</Header>
            <Form
              error={hasError}
              loading={submitting}
              onSubmit={this.handleSubmit}
            >
              {hasError ? (
                <Message header="Whoops!" content={error.message} error />
              ) : null}
              <Form.Field>
                <label>E-Mail</label>
                <input placeholder="E-Mail" ref={e => (this.email = e)} />
              </Form.Field>
              <Form.Field>
                <label>Password</label>
                <input
                  type="password"
                  placeholder="Password"
                  ref={e => (this.password = e)}
                />
              </Form.Field>
              <Button type="submit">Log in</Button>
            </Form>
          </Segment>
          <Footer activeItem="login" authAccount={null} />
        </div>
      </div>
    );
  }
}
