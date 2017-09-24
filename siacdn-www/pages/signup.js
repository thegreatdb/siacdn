//import React from 'react';
import Head from 'next/head';
import Link from 'next/link';
import cookies from 'next-cookies';
import {
  Segment,
  Header,
  Button,
  Form,
} from 'semantic-ui-react';
import Client from '../lib/client';
import redirect from '../lib/redirect';

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
        <Header as="h1">Sign up</Header>
        <Form onSubmit={(ev) => {console.log(ev)}}>
          <Form.Field>
            <label>Username</label>
            <input placeholder='Username' />
          </Form.Field>
          <Form.Field>
            <label>Password</label>
            <input type="password" placeholder='Password' />
          </Form.Field>
          <Form.Field>
            <label>Password (Repeat)</label>
            <input type="password" placeholder='Password (Repeat)' />
          </Form.Field>
          <Button type='submit'>Sign up</Button>
        </Form>
      </Segment>
    </div>
  </div>
);
