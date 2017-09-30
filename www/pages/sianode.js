import Head from 'next/head';
import Link from 'next/link';
import Router from 'next/router';
import cookies from 'next-cookies';
import { Segment, Header, Button, List, Card, Icon } from 'semantic-ui-react';
import TimeAgo from 'timeago-react';
import Client from '../lib/client';
import redirect from '../lib/redirect';
import Nav from '../components/nav';
import { displayStatus } from '../lib/fmt';

const SiaNode = ({ authAccount, siaNode }) => (
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
      <Nav activeItem="sianode" authAccount={authAccount} />
      <Segment padded>
        <Header as="h1">
          Sia full node: {siaNode.shortcode}
        </Header>
        <List>
          <List.Item>
            <List.Content>
              <Icon name="tag" />{' '}
              <strong>ID</strong>: {siaNode.id}
            </List.Content>
          </List.Item>
          <List.Item>
            <List.Content>
              <Icon name="signal" />{' '}
              <strong>Status</strong>: {displayStatus[siaNode.status]}
            </List.Content>
          </List.Item>
          <List.Item>
            <List.Content>
              <Icon name="time" />{' '}
              <strong>Created</strong>: <TimeAgo datetime={siaNode.created_time} />
            </List.Content>
          </List.Item>
          <List.Item>
            <List.Content>
              <Icon name="cloud" />{' '}
              {siaNode.minio_instances_requested} Minio instance{siaNode.minio_instances_requested === 1 ? '' : 's'}
            </List.Content>
          </List.Item>
        </List>
      </Segment>
    </div>
  </div>
);

SiaNode.getInitialProps = async ctx => {
  const { authTokenID } = cookies(ctx);
  const { query: { id } } = ctx;
  const client = new Client(authTokenID);
  let authAccount = null;
  try {
    authAccount = await client.getAuthAccount();
    if (!authAccount) {
      redirect(ctx, '/login');
    }
  } catch (err) {
    redirect(ctx, '/login');
  }
  const siaNode = await client.getSiaNode(id);
  return { authAccount, siaNode };
};

export default SiaNode;
