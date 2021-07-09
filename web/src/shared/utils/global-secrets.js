import { displayMessage } from "./message";

/**
 * Get the global secret list
 * store the results in the state tree.
 *
 * @param {Object} tree - The drone state tree.
 * @param {Object} client - The drone client.
 */
export const fetchGlobalSecretList = (tree, client) => {
  tree.unset(["globalSecrets", "loaded"]);
  tree.unset(["globalSecrets", "error"]);

  client.getGlobalSecretList().then(results => {
    let list = {};
    results.map(secret => {
      list[secret.name] = secret;
    });
    tree.set(["globalSecrets", "data"], list);
    tree.set(["globalSecrets", "loaded"], true);
  });
};

/**
 * Create secret and if successful
 * store the result in the state tree.
 *
 * @param {Object} tree - The drone state tree.
 * @param {Object} client - The drone client.
 * @param {Object} secret - The secret object.
 */
export const createGlobalSecret = (tree, client, secret) => {
  client
    .createGlobalSecret(secret)
    .then(result => {
      tree.set(["globalSecrets", "data", secret.name], result);
      displayMessage(tree, "Successfully added the global secret");
    })
    .catch(() => {
      displayMessage(tree, "Failed to create the global secret");
    });
};

/**
 * Delete secret from the server and
 * remove from the state tree.
 *
 * @param {Object} tree - The drone state tree.
 * @param {Object} client - The drone client.
 * @param {string} secret - The secret name.
 */
export const deleteGlobalSecret = (tree, client, secret) => {
  client
    .deleteGlobalSecret(secret)
    .then(result => {
      tree.unset(["globalSecrets", "data", secret]);
      displayMessage(tree, "Successfully removed the global secret");
    })
    .catch(() => {
      displayMessage(tree, "Failed to remove the global secret");
    });
};
