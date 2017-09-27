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
  Message,
} from 'semantic-ui-react';
import Nav from '../components/nav';
import redirect from '../lib/redirect';
import Client from '../lib/client';

const siaCostOptions = [
  { key: 5, text: ' 5TB - $6/mo', value: '5' },
  { key: 10, text: '10TB - $12/mo', value: '10' },
  { key: 15, text: '15TB - $18/mo', value: '15' },
  { key: 20, text: '20TB - $24/mo', value: '20' },
  { key: 25, text: '25TB - $30/mo', value: '25' },
  { key: 30, text: '30TB - $36/mo', value: '30' },
  { key: 35, text: '35TB - $42/mo', value: '35' },
  { key: 40, text: '40TB - $48/mo', value: '40' },
  { key: 45, text: '45TB - $54/mo', value: '45' },
  { key: 50, text: '50TB - $60/mo', value: '50' },
];

const displayStatus = {
  created: 'Sending specifications to the deployment server...',
  deployed: 'Waiting for resources from deployment server...',
  instanced: 'Initialized and snapshotting...',
  snapshotted:
    'Finished snapshotting, downloading latest blockchain updates...',
  synchronized: 'Initializing wallet...',
  initialized: 'Transferring funds based on your requested capacity...',
  funded: 'Setting allowance...',
  configured: 'Waiting for contracts to finalize...',
  ready: 'Sia node is up and running!',
  stopped: 'Stopped.',
  depleted: 'Insufficient funds to continue.',
  error: 'Errored.',
};

export default class NewSia extends React.Component {
  state = {
    selectedCost: -1,
    siaError: null,
    siaSubmitting: false,
    siaNode: null,
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
    const pendingSiaNode = await client.getPendingSiaNode();
    return { authTokenID, authAccount, pendingSiaNode };
  }

  handleSiaCapacityChange = async (ev, data) => {
    await this.setState({ selectedCost: data.value });
  };

  handleSiaSubmit = async ev => {
    ev.preventDefault();

    const { authTokenID } = this.props;
    const { selectedCost } = this.state;
    if (selectedCost < 0) {
      await this.setState({
        siaError: { message: 'You should select a capacity' },
      });
      return;
    }
    await this.setState({ siaError: null, siaSubmitting: true });
    try {
      const client = new Client(authTokenID);
      const siaNode = await client.createSiaNode(
        siaCostOptions[selectedCost].key
      );
      await this.setState({ siaSubmitting: false, siaNode: siaNode });
    } catch (error) {
      await this.setState({ siaError: error, siaSubmitting: false });
    }
  };

  checkSia = async () => {
    let { siaNode } = this.state;
    const { authTokenID } = this.props;

    if (siaNode && siaNode.state == 'ready') {
      console.log('Found ready sia node, shutting down Sia check');
      return;
    }

    await this.setState({ siaError: null });
    try {
      siaNode = await new Client(authTokenID).getPendingSiaNode();
    } catch (error) {
      await this.setState({ siaError: error });
      return;
    }
    await this.setState({ siaNode });

    this.timer = setTimeout(() => this.checkSia(), 1000);
  };

  componentDidMount = () => {
    this.timer = setTimeout(() => this.checkSia(), 0);
  };

  componentWillUnmount = () => {
    if (this.timer) {
      clearTimeout(this.timer);
      this.timer = null;
    }
  };

  render() {
    const { authAccount } = this.props;
    const { selectedCost, siaSubmitting, siaError } = this.state;
    const hasSiaError = Boolean(siaError);
    const siaNode = this.state.siaNode || this.props.pendingSiaNode;
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
            {siaNode ? (
              <Header as="h1">
                Now let&rsquo;s start some Minio instances
              </Header>
            ) : (
              <Header as="h1">Let&rsquo;s start a new Sia full node</Header>
            )}
          </Segment>
          <Step.Group ordered fluid size="small">
            <Step
              completed
              title="Sign up"
              description="Sign up for an account"
            />
            <Step
              completed={Boolean(siaNode)}
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
            {siaNode ? (
              <Message info>
                <Message.Header>Setting it up</Message.Header>
                <Message.Content>
                  Setting up your Sia full node<br />
                  Name: {siaNode.shortcode}<br />
                  Current Status:{' '}
                  <strong>{displayStatus[siaNode.status]}</strong>
                </Message.Content>
              </Message>
            ) : (
              <Form
                error={hasSiaError}
                loading={siaSubmitting}
                onSubmit={this.handleSiaSubmit}
              >
                {hasSiaError ? (
                  <Message header="Whoops!" content={siaError.message} error />
                ) : null}
                <Form.Field>
                  <label>Node base monthly price</label>
                  $10
                </Form.Field>
                <Form.Field>
                  <label>Sia node capacity</label>
                  <Form.Select
                    options={siaCostOptions}
                    onChange={this.handleSiaCapacityChange}
                    placeholder="Sia node capacity"
                  />
                </Form.Field>
                <Button>Start Sia node</Button>
              </Form>
            )}
          </Segment>
        </div>
      </div>
    );
  }
}
