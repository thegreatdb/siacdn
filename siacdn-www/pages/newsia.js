import React from 'react';
import Head from 'next/head';
import Link from 'next/link';
import cookies from 'next-cookies';
import {
  Segment,
  Grid,
  Item,
  Header,
  Button,
  List,
  Step,
  Form,
  Menu,
  Dropdown,
} from 'semantic-ui-react';
import Nav from '../components/nav';
import redirect from '../lib/redirect';
import Client from '../lib/client';

const siaCostOptions = [
  { key: 1, text: '5TB  ($16/mo)', value: 1 },
  { key: 2, text: '10TB ($22/mo)', value: 2 },
  { key: 3, text: '15TB ($28/mo)', value: 3 },
  { key: 4, text: '20TB ($34/mo)', value: 4 },
  { key: 5, text: '25TB ($40/mo)', value: 5 },
  { key: 6, text: '30TB ($46/mo)', value: 6 },
  { key: 7, text: '35TB ($52/mo)', value: 7 },
  { key: 8, text: '40TB ($58/mo)', value: 8 },
  { key: 9, text: '45TB ($64/mo)', value: 9 },
  { key: 10, text: '50TB ($70/mo)', value: 10 },
];

export default class NewSia extends React.Component {
  state = {
    stage: 'sia',
    selectedCost: 1,
  };

  static async getInitialProps(ctx) {
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
  }

  async handleSiaCapacityChange(ev, data) {
    await this.setState({ selectedCost: data.value });
  }

  render() {
    const { authAccount } = this.props;
    const { stage, selectedCost } = this.state;
    return (
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
            {stage === 'sia' ? (
              <Header as="h1">Let&rsquo;s start a new Sia full node</Header>
            ) : (
              <Header as="h1">
                Now let&rsquo;s start some Minio instances
              </Header>
            )}
          </Segment>
          <Step.Group ordered fluid size="small">
            <Step
              completed
              title="Sign up"
              description="Sign up for an account"
            />
            <Step
              completed={stage !== 'sia'}
              title="Sia Node"
              description="Configure Sia node"
            />
            <Step
              title="Minio Instances"
              description="Set up Minio instances"
            />
          </Step.Group>
          <Segment padded>
            <Header as="h3">Sia Node</Header>
            <Form>
              <Form.Field>
                <label>Sia node capacity</label>
                <Form.Select
                  options={siaCostOptions}
                  onChange={this.handleSiaCapacityChange.bind(this)}
                  placeholder="Sia node capacity"
                />
              </Form.Field>
              <Button>Start Sia node</Button>
            </Form>
          </Segment>
        </div>
      </div>
    );
  }
}
