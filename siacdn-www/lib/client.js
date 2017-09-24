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

  createAccount(email, password) {
    return this._reg('/accounts', email, password);
  }

  loginAccount(email, password) {
    return this._reg('/auth', email, password);
  }

  /*
  async uploadFile(file) {
    const formData = new FormData();
    formData.append('file', file);
    const resp = await this.post('/files', null, formData);
    return resp.file;
  }

  async getFile(id) {
    const resp = await this.get('/files/id/' + id, null);
    return resp.file;
  }
  */

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

  async _reg(path, email, password) {
    const body = JSON.stringify({ email, password });
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
