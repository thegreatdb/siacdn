import Head from 'next/head';
import Link from 'next/link';
import cookies from 'next-cookies';
import { Segment, Grid, Item, Header, Button, List } from 'semantic-ui-react';
import Nav from '../components/nav';
import Client from '../lib/client';

const Index = ({ authAccount }) => (
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
      <Nav activeItem="index" authAccount={authAccount} />
      <Segment padded>
        <Header as="h1">SiaCDN</Header>
        <p>
          <strong>SiaCDN</strong> is the easiest way to get started with Sia in
          the cloud.
        </p>
        <p>
          We spin up a full Sia node for you, along with a specialized version
          of Minio which provides an S3-compatible frontend to your Sia full
          node.
        </p>

        <Header as="h2">Price</Header>
        <p>
          We charge a flat fee for our services of <strong>$30/month</strong>{' '}
          per full Sia node, <strong>$10/month</strong> per Minio instance, and{' '}
          <strong>$0.03/GB</strong> bandwidth out. This is more than we&rsquo;d
          like to charge, but it&rsquo;s because we currently have to proxy all
          traffic and pay those bandwidth costs. In the future we will be able
          to cut this down dramatically when the Sia network{' '}
          <strong>adds CDN functionality</strong>. Note that you cannot buy any
          SiaCoins from us - you are paying for our hosting services and we use
          the Sia network to provide those services to you.
        </p>

        <Header as="h2">Steps</Header>
        <Grid columns={3} divided>
          <Grid.Row>
            <Grid.Column>
              <Item.Group divided>
                <Item>
                  <Item.Image
                    size="tiny"
                    src="https://react.semantic-ui.com/assets/images/wireframe/image.png"
                  />
                  <Item.Content verticalAlign="middle">
                    Start a Sia full node in the cloud
                  </Item.Content>
                </Item>
              </Item.Group>
            </Grid.Column>
            <Grid.Column>
              <Item.Group divided>
                <Item>
                  <Item.Image
                    size="tiny"
                    src="https://react.semantic-ui.com/assets/images/wireframe/image.png"
                  />
                  <Item.Content verticalAlign="middle">
                    Choose your scaling options
                  </Item.Content>
                </Item>
              </Item.Group>
            </Grid.Column>
            <Grid.Column>
              <Item.Group divided>
                <Item>
                  <Item.Image
                    size="tiny"
                    src="https://react.semantic-ui.com/assets/images/wireframe/image.png"
                  />
                  <Item.Content verticalAlign="middle">
                    Connect to S3-compatible Minio frontend
                  </Item.Content>
                </Item>
              </Item.Group>
            </Grid.Column>
          </Grid.Row>
        </Grid>

        <Header as="h2">Let&rsquo;s get going</Header>
        <Link href="/dashboard">
          <Button primary>Go to your dashboard</Button>
        </Link>
      </Segment>
    </div>
  </div>
);

Index.getInitialProps = async ctx => {
  const { authTokenID } = cookies(ctx);
  const client = new Client(authTokenID);
  let authAccount = null;
  try {
    authAccount = await client.getAuthAccount();
  } catch (err) {}
  return { authAccount };
};

export default Index;
