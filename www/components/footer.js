import React from 'react';
import Link from 'next/link';
import { Button, Menu, Icon } from 'semantic-ui-react';

const Footer = ({ activeItem, authAccount }) => (
  <div className="ui vertical footer segment">
    <div className="ui center aligned container">
      <div className="ui horizontal small divided link list">
        <span className="item">
          Copyright &copy; {new Date().getFullYear()} Maxint, LLC.
        </span>
        <Link href="/">
          <a href="/" className="item">
            Home
          </a>
        </Link>
        <Link href="/dashboard">
          <a className="item">Dashboard</a>
        </Link>
        <a className="item" href="mailto:eric@maxint.co">
          Contact Us
        </a>
        <Link href="/tos">
          <a className="item" href="/tos">
            Terms
          </a>
        </Link>
        <Link href="/copyright">
          <a className="item" href="/copyright">
            Copyright
          </a>
        </Link>
        <Link href="/privacy">
          <a className="item" href="/privacy">
            Privacy Policy
          </a>
        </Link>
        <a className="item" href="https://github.com/thegreatdb/siacdn" target="_blank">
          <Icon name="github" />
          Open Source
        </a>
      </div>
    </div>
  </div>
);

export default Footer;
