import Head from 'next/head';
import Link from 'next/link';
import cookies from 'next-cookies';
import { Segment, Grid, Item, Header, Button, List } from 'semantic-ui-react';
import Nav from '../components/nav';
import redirect from '../lib/redirect';
import Client from '../lib/client';

const NewSia = ({ authAccount }) => (
  <div>
    <Head>
      <link
        rel="stylesheet"
        href="//cdnjs.cloudflare.com/ajax/libs/semantic-ui/2.2.2/semantic.min.css"
      />
      <link rel="stylesheet" href="/static/css/global.css" />
      <script src="https://js.stripe.com/v3/" />
    </Head>
    <div className="holder">
      <Nav activeItem="newsia" authAccount={authAccount} />
      <Segment padded>
        <Header as="h1">Let&rsquo;s start a new Sia full node</Header>
      </Segment>
      <Segment padded>
        <Header as="h3">Fund account</Header>
      </Segment>
      <Segment padded>
        <Header as="h3">Do something</Header>
      </Segment>
      <Segment padded>
        <Header as="h3">Think of third point</Header>
      </Segment>
    </div>
  </div>
);

NewSia.getInitialProps = async ctx => {
  const { authTokenID } = cookies(ctx);
  const client = new Client(authTokenID);
  let authAccount = null;
  try {
    authAccount = await client.getAuthAccount();
    if (!authAccount) {
      redirect(ctx, '/signup');
    }
  } catch (err) {
    redirect(ctx, '/signup');
  }
  return { authAccount };
};

export default NewSia;
