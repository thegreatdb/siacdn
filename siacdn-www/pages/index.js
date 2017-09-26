import Head from 'next/head';
import Link from 'next/link';
import cookies from 'next-cookies';
import {
  Segment,
  Step,
  Item,
  Header,
  Button,
  List,
  Message,
  Icon,
} from 'semantic-ui-react';
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
          Increasing overall network usage is the real best way to show your
          support. We believe that if it&rsquo;s easy to get started with Sia —
          if, within an hour of hearing about it, a developer who has used
          Amazon Web Services can use it, then overall network usage will
          skyrocket. We believe SiaCDN achieves this, and we hope you support us
          in supporting the Sia network.
        </Message>

        <Header as="h2">Price</Header>
        <p>
          We charge a flat fee for our services of <strong>$10/month</strong>{' '}
          per Sia full node, <strong>$10/month</strong> per Minio instance,{' '}
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
          for bandwidth out. This is more than we&rsquo;d like to charge, but
          it&rsquo;s because we currently have to proxy all traffic and pay
          those bandwidth costs. In the future we will reduce the costs
          substantially by offloading the majority of the work to the Sia
          network itself (once it has a few more features.)
        </p>
        <Message info>
          Note that{' '}
          <strong>
            you cannot buy virtual currency of any kind from us.
          </strong>{' '}
          You are paying for our internet hosting service, and we use the Sia
          network in part of providing that service to you.
        </Message>

        <Header as="h2">Steps</Header>
        <Step.Group size="small">
          <Step>
            <Icon name="mouse pointer" />
            <Step.Content
              title="Sign up"
              description="Register for an account"
            />
          </Step>
          <Step>
            <Icon name="line chart" />
            <Step.Content
              title="Scale up"
              description="Start Sia & Minio nodes"
            />
          </Step>
          <Step>
            <Icon name="birthday" />
            <Step.Content
              title="Enjoy"
              description="Visit dashboard & use S3 tools"
            />
          </Step>
        </Step.Group>

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
