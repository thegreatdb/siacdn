import Head from 'next/head';
import Link from 'next/link';
import Router from 'next/router';
import cookies from 'next-cookies';
import { Segment, Header, Button, List, Card, Icon } from 'semantic-ui-react';
import TimeAgo from 'timeago-react';
import HttpsRedirect from 'react-https-redirect';
import Client from '../lib/client';
import redirect from '../lib/redirect';
import Nav from '../components/nav';
import { displayStatus } from '../lib/fmt';

const Dashboard = ({ authAccount, siaNodes }) => (
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
      <Nav activeItem="dashboard" authAccount={authAccount} />
      <Segment padded>
        <Header as="h1">SiaCDN</Header>
        <Card.Group>
          {siaNodes ? siaNodes.map(siaNode => (
            <Card key={siaNode.shortcode} fluid href={siaNode.status === 'ready' ? '/sianode?id='+siaNode.id : '/newsia'} onClick={(ev) => {
              ev.preventDefault();
              ev.stopPropagation();
              Router.push(siaNode.status === 'ready' ? '/sianode?id='+siaNode.id : '/newsia');
            }}>
              <Card.Content header={'Sia full node: ' + siaNode.shortcode} />
              <Card.Content>
                <Card.Description>
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
                  </List>
                </Card.Description>
              </Card.Content>
              <Card.Content extra>
                <Icon name="cloud" />
                {siaNode.minio_instances_requested} Minio instance{siaNode.minio_instances_requested === 1 ? '' : 's'}
              </Card.Content>
            </Card>
          )) : null}
        </Card.Group>
      </Segment>
    </div>
  </HttpsRedirect>
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
  const siaNodes = await client.getSiaNodes();
  return { authAccount, siaNodes };
};

export default Dashboard;
