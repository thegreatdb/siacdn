import Head from 'next/head';
import Link from 'next/link';
import {
  Segment,
  Grid,
  Item,
  Header,
  Button,
  List,
} from 'semantic-ui-react';

export default () => (
  <div>
    <Head>
      <link
        rel="stylesheet"
        href="//cdnjs.cloudflare.com/ajax/libs/semantic-ui/2.2.2/semantic.min.css"
      />
      <link rel="stylesheet" href="/static/css/global.css" />
    </Head>
    <div className="holder">
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
                    Fund your account with SiaCoins
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
                    Set your budget
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
