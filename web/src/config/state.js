import Baobab from "baobab";

const user = window.DRONE_USER;
const sync = window.DRONE_SYNC;

const state = {
  follow: false,
  language: "en-US",

  user: {
    data: user,
    error: undefined,
    loaded: true,
    syncing: sync,
  },

  feed: {
    loaded: false,
    error: undefined,
    data: {},
  },

  repos: {
    loaded: false,
    error: undefined,
    data: {},
  },

  globalSecrets: {
    loaded: false,
    error: undefined,
    data: {},
  },

  secrets: {
    loaded: false,
    error: undefined,
    data: {},
  },

  registry: {
    error: undefined,
    loaded: false,
    data: {},
  },

  builds: {
    loaded: false,
    error: undefined,
    data: {},
  },

  logs: {
    follow: false,
    loading: true,
    error: false,
    data: {},
  },

  token: {
    value: undefined,
    error: undefined,
    loading: false,
  },

  message: {
    show: false,
    text: undefined,
    error: false,
  },

  location: {
    protocol: window.location.protocol,
    host: window.location.host,
  },
};

const tree = new Baobab(state);

if (window) {
  window.tree = tree;
}

export default tree;
