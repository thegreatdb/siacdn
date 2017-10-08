import React from 'react';
import Head from 'next/head';

const PageHeader = ({ children }) => (
  <Head>
    <link
      rel="stylesheet"
      href="//cdnjs.cloudflare.com/ajax/libs/semantic-ui/2.2.2/semantic.min.css"
    />
    <link rel="stylesheet" href="/static/css/global.css" />
    <script src="https://js.stripe.com/v3/" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    {children}
  </Head>
);

export default PageHeader;
