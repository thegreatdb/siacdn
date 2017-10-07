import React from 'react';
import clickRouter from '../lib/click-router';
import { Button, Menu, Icon } from 'semantic-ui-react';

const Nav = ({ activeItem, authAccount }) => (
  <Menu className="nav" stackable>
    <Menu.Item
      name="index"
      active={activeItem === 'index'}
      onClick={clickRouter('/')}
    >
      <Icon name="home" />
    </Menu.Item>
    <Menu.Item
      name="dashboard"
      active={activeItem === 'dashboard'}
      onClick={clickRouter('/dashboard')}
    />

    {activeItem === 'logout' ? null : (
      <Menu.Menu position="right">
        {authAccount ? <Menu.Item>{authAccount.name}</Menu.Item> : null}
        {authAccount ? (
          <Menu.Item>
            <Button
              onClick={clickRouter('/newsia')}
              basic
              primary
              disabled={activeItem === 'newsia'}
            >
              <Icon name="database" />New Sia node
            </Button>
          </Menu.Item>
        ) : null}
        <Menu.Item>
          {authAccount ? null :
            <Button onClick={clickRouter('/signup')} primary>
              Sign up
            </Button>}
          {authAccount ? null : <span className="sep">or</span>}
          {authAccount ?
            <Button onClick={clickRouter('/logout')}>Logout</Button> :
            <Button onClick={clickRouter('/login')} primary>
              Login
            </Button>}
        </Menu.Item>
      </Menu.Menu>
    )}
    <style jsx>{`
      .sep {
        padding: 0 .5em 0 .5em;
      }
    `}</style>
  </Menu>
);

export default Nav;
