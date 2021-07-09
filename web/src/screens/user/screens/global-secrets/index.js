import React, { Component } from "react";

import {
  fetchGlobalSecretList,
  createGlobalSecret,
  deleteGlobalSecret,
} from "shared/utils/global-secrets";

import { branch } from "baobab-react/higher-order";
import { inject } from "config/client/inject";

import { List, Item, Form } from "screens/repo/screens/secrets/components";

import styles from "screens/repo/screens/secrets/index.less";

const binding = (props, context) => {
  return {
    loaded: ["globalSecrets", "loaded"],
    secrets: ["globalSecrets", "data"],
  };
};

@inject
@branch(binding)
export default class GlobalSecrets extends Component {
  constructor(props, context) {
    super(props, context);

    this.handleSave = this.handleSave.bind(this);
  }

  shouldComponentUpdate(nextProps, nextState) {
    return this.props.secrets !== nextProps.secrets;
  }

  componentWillMount() {
    this.props.dispatch(fetchGlobalSecretList, this.props.drone);
  }

  handleSave(e) {
    const { dispatch, drone } = this.props;
    const secret = {
      name: e.detail.name,
      value: e.detail.value,
      event: e.detail.event,
    };

    dispatch(createGlobalSecret, drone, secret);
  }

  handleDelete(secret) {
    const { dispatch, drone } = this.props;
    dispatch(deleteGlobalSecret, drone, secret.name);
  }

  render() {
    const { secrets, loaded } = this.props;

    if (!loaded) {
      return LOADING;
    }

    return (
      <div className={styles.root}>
        <div className={styles.left}>
          {Object.keys(secrets || {}).length === 0 ? EMPTY : undefined}
          <List>
            {Object.values(secrets || {}).map(renderGlobalSecret.bind(this))}
          </List>
        </div>
        <div className={styles.right}>
          <Form onsubmit={this.handleSave} />
        </div>
      </div>
    );
  }
}

function renderGlobalSecret(secret) {
  return (
    <Item
      name={secret.name}
      event={secret.event}
      ondelete={this.handleDelete.bind(this, secret)}
    />
  );
}

const LOADING = <div className={styles.loading}>Loading</div>;

const EMPTY = (
  <div className={styles.empty}>
    There are no global secrets for this repository.
  </div>
);
