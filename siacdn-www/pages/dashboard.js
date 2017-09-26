import Head from 'next/head';
import Link from 'next/link';
import cookies from 'next-cookies';
import { Segment, Header, Button, List } from 'semantic-ui-react';
import Client from '../lib/client';
import redirect from '../lib/redirect';
import Nav from '../components/nav';

const Dashboard = ({ authAccount }) => (
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
      <Nav activeItem="dashboard" authAccount={authAccount} />
      <Segment padded>
        <Header as="h1">SiaCDN</Header>
        <p>This is your dashboard</p>
      </Segment>
    </div>
  </div>
);

Dashboard.getInitialProps = async ctx => {
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

export default Dashboard;
