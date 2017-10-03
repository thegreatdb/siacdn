import React from 'react';
import Router from 'next/router';
import cookies from 'next-cookies';
import HttpsRedirect from 'react-https-redirect';
import { Message } from 'semantic-ui-react';
import Client from '../lib/client';

export default class Logout extends React.Component {
  async getInitialProps(ctx) {
    const { authTokenID } = cookies(ctx);
    return { authTokenID };
  }

  componentDidMount() {
    const { authTokenID } = this.props;
    const client = new Client(authTokenID);
    client.removeAuthTokenID();
    Router.push('/');
  }

  render() {
    return (
      <HttpsRedirect>
        <Message header="Please wait..." content="Logging you out now." />
      </HttpsRedirect>
    );
  }
}
