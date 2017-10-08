import Head from 'next/head';
import Link from 'next/link';
import cookies from 'next-cookies';
import HttpsRedirect from 'react-https-redirect';
import {
  Segment,
  Step,
  Item,
  Header,
  Button,
  List,
  Message,
  Icon,
  Grid,
} from 'semantic-ui-react';
import Nav from '../components/nav';
import Footer from '../components/footer';
import Client from '../lib/client';

const Index = ({ authAccount }) => (
  <HttpsRedirect>
    <Head>
      <link
        rel="stylesheet"
        href="//cdnjs.cloudflare.com/ajax/libs/semantic-ui/2.2.2/semantic.min.css"
      />
      <link rel="stylesheet" href="/static/css/global.css" />
      <script src="https://js.stripe.com/v3/" />
      <meta name="viewport" content="width=device-width, initial-scale=1" />
    </Head>
    <div className="holder">
      <Nav activeItem="index" authAccount={authAccount} />
      <Segment padded>
        <Header as="h2">
          <strong>SiaCDN</strong> is the easiest way to get started with Sia in
          the cloud.
          <Header.Subheader>
            We host a Sia full node for you, along with a specialized version of
            Minio that provides an S3-compatible API into your Sia node.
          </Header.Subheader>
        </Header>

        <Message positive>
          <Message.Header>
            If you&rsquo;re interested in supporting distributed systems...
          </Message.Header>
          <Message.Content>
            <p>
              Increasing overall network usage is the real best way to show your
              support. We believe that if it&rsquo;s easy to get started with
              Sia — if, within a few hours of hearing about it, a developer who
              has used AWS S3 can use it, then overall network usage will
              skyrocket. We believe SiaCDN achieves this, and we hope you
              support us in supporting the Sia network.
            </p>
            <p>
              <Link href="/dashboard">
                <Button basic color="green">
                  Get started
                </Button>
              </Link>
            </p>
          </Message.Content>
        </Message>

        <Grid columns={2} doubling stackable padded>
          <Grid.Column>
            <Item>
              <Item.Content>
                <Item.Header as="h2">
                  <Icon name="credit card alternative" /> Price
                </Item.Header>
                <Item.Description>
                  We charge a flat fee for our services of{' '}
                  <strong>$10/month</strong> per Sia full node,{' '}
                  <strong>$10/month</strong> per Minio instance,{' '}
                  <strong>
                    <a
                      href="https://siastats.info/storage_pricing.html"
                      target="_blank"
                    >
                      $1.20/TB
                    </a>
                  </strong>{' '}
                  for Sia network storage capacity and{' '}
                  <strong>
                    <a
                      href="https://cloud.google.com/compute/pricing#internet_egress"
                      target="_blank"
                    >
                      $0.025/GB
                    </a>
                  </strong>{' '}
                  for bandwidth out. This is more than we&rsquo;d like to
                  charge, but it&rsquo;s because we currently have to proxy all
                  traffic and pay those bandwidth costs. In the future we will
                  reduce the costs substantially by offloading the majority of
                  the work to the Sia network itself (once it has a few more
                  features.)
                  <Message info>
                    Note that{' '}
                    <strong>
                      you cannot buy virtual currency of any kind from us.
                    </strong>{' '}
                    You are paying for our internet hosting service, and we use
                    the Sia network in part of providing that service to you.
                  </Message>
                </Item.Description>
              </Item.Content>
            </Item>
          </Grid.Column>

          <Grid.Column>
            <Item>
              <Item.Content>
                <Item.Header as="h2">
                  <Icon name="video play" /> Walkthrough
                </Item.Header>
                <Item.Description>
                  <iframe
                    src="https://www.youtube.com/embed/bfxXzcAo_J4?rel=0&amp;showinfo=0"
                    frameBorder="0"
                    allowFullScreen
                  />
                </Item.Description>
              </Item.Content>
            </Item>
          </Grid.Column>

          <Grid.Column width={16} textAlign="center" verticalAlign="middle">
            <Item>
              <Item.Content>
                <Item.Header as="h2">
                  <Icon name="line chart" /> Let&rsquo;s get going
                </Item.Header>
                <Item.Description>
                  <Link href="/dashboard">
                    <Button primary>Go to your dashboard</Button>
                  </Link>
                </Item.Description>
              </Item.Content>
            </Item>
          </Grid.Column>
        </Grid>
      </Segment>
      <Footer activeItem="index" authAccount={authAccount} />
    </div>
    <style jsx>{`
      iframe {
        width: 100%;
        min-height: 304px;
      }
    `}</style>
  </HttpsRedirect>
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
