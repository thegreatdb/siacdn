import fetch from 'isomorphic-fetch';
import queryString from 'query-string';

export default class Client {
  constructor(authTokenID) {
    this.authTokenID = authTokenID;
    this.authAccount = null;
    this.base = 'http://localhost:9095';
  }

  setAuthTokenID(id) {
    if (!id) {
      throw new Error('Invalid auth token id: ' + id);
    }
    this.authTokenID = id;
    document.cookie = 'authTokenID=' + id;
  }

  removeAuthTokenID() {
    this.authTokenID = null;
    this.authAccount = null;
    document.cookie = 'authTokenID=;expires=Thu, 01 Jan 1970 00:00:01 GMT;';
  }

  async getAuthAccount() {
    if (this.authAccount) {
      return this.authAccount;
    }
    const resp = await this.get('/auth');
    this.authAccount = resp['account'];
    return this.authAccount;
  }

  createAccount(username, password, name, stripeToken) {
    return this._reg('/accounts', username, password, name, stripeToken);
  }

  loginAccount(username, password) {
    return this._reg('/auth', username, password);
  }

  // Supporting and utility functions follow

  headers() {
    const headers = { 'Content-Type': 'application/json' };
    if (this.authTokenID) {
      headers['X-Auth-Token-ID'] = this.authTokenID;
    }
    return headers;
  }

  get(path, params) {
    const url =
      this.base + path + (params ? queryString.stringify(params) : '');
    return fetchJSON(url, { headers: this.headers() });
  }

  post(path, params, body) {
    const url =
      this.base + path + (params ? queryString.stringify(params) : '');
    const opts = { headers: this.headers(), method: 'post' };
    if (body instanceof FormData) {
      delete opts.headers['Content-Type'];
    }
    if (body) {
      opts['body'] = body; // TODO: JSON.stringify if not string already
    }
    return fetchJSON(url, opts);
  }

  async _reg(path, username, password, name, stripeToken) {
    let dat = { username, password };
    if (name) {
      dat['name'] = name;
    }
    if (stripeToken) {
      dat['stripe_token'] = stripeToken;
    }
    const body = JSON.stringify(dat);
    const resp = await this.post(path, null, body);
    if (!resp.auth_token || !resp.auth_token.id) {
      throw new Error('Got no token id in response: ' + JSON.stringify(resp));
    }
    this.setAuthTokenID(resp.auth_token.id);
    return resp.account;
  }
}

const fetchJSON = async (url, opts) => {
  const res = await fetch(url, opts);
  const json = await res.json();
  if (json && json['error']) {
    throw new Error(json['error']);
  }
  return json;
};
