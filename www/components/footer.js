import React from 'react';
import Link from 'next/link';
import { Button, Menu, Icon } from 'semantic-ui-react';

const Footer = ({ activeItem, authAccount }) => (
  <div className="ui vertical footer segment">
    <div className="ui center aligned container">
      <div className="ui horizontal small divided link list">
        <span className="item">Copyright &copy; {(new Date()).getFullYear()} Maxint, LLC.</span>
        <Link href="/"><a className="item">Home</a></Link>
        <Link href="/dashboard"><a className="item">Dashboard</a></Link>
        <a className="item" href="mailto:eric@maxint.co">Contact Us</a>
        <Link href="/tos"><a className="item" href="#">Terms and Conditions</a></Link>
        <Link href="/privacy"><a className="item" href="#">Privacy Policy</a></Link>
      </div>
    </div>
  </div>
);

export default Footer;
