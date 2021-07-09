(function(global, factory) {
  if (typeof define === "function" && define.amd) {
    define(["exports"], factory);
  } else if (typeof exports !== "undefined") {
    factory(exports);
  } else {
    var mod = {
      exports: {}
    };
    factory(mod.exports);
    global.index = mod.exports;
  }
})(this, function(exports) {
  "use strict";

  Object.defineProperty(exports, "__esModule", {
    value: true
  });

  function _classCallCheck(instance, Constructor) {
    if (!(instance instanceof Constructor)) {
      throw new TypeError("Cannot call a class as a function");
    }
  }

  var _createClass = (function() {
    function defineProperties(target, props) {
      for (var i = 0; i < props.length; i++) {
        var descriptor = props[i];
        descriptor.enumerable = descriptor.enumerable || false;
        descriptor.configurable = true;
        if ("value" in descriptor) descriptor.writable = true;
        Object.defineProperty(target, descriptor.key, descriptor);
      }
    }

    return function(Constructor, protoProps, staticProps) {
      if (protoProps) defineProperties(Constructor.prototype, protoProps);
      if (staticProps) defineProperties(Constructor, staticProps);
      return Constructor;
    };
  })();

  var DroneClient = (function() {
    function DroneClient(server, token, csrf) {
      _classCallCheck(this, DroneClient);

      this.server = server || "";
      this.token = token;
      this.csrf = csrf;
    }

    _createClass(
      DroneClient,
      [
        {
          key: "getRepoList",
          value: function getRepoList(opts) {
            var query = encodeQueryString(opts);
            return this._get("/api/user/repos?" + query);
          }
        },
        {
          key: "getRepo",
          value: function getRepo(owner, repo) {
            return this._get("/api/repos/" + owner + "/" + repo);
          }
        },
        {
          key: "activateRepo",
          value: function activateRepo(owner, repo) {
            return this._post("/api/repos/" + owner + "/" + repo);
          }
        },
        {
          key: "updateRepo",
          value: function updateRepo(owner, repo, data) {
            return this._patch("/api/repos/" + owner + "/" + repo, data);
          }
        },
        {
          key: "deleteRepo",
          value: function deleteRepo(owner, repo) {
            return this._delete("/api/repos/" + owner + "/" + repo);
          }
        },
        {
          key: "getBuildList",
          value: function getBuildList(owner, repo, opts) {
            var query = encodeQueryString(opts);
            return this._get(
              "/api/repos/" + owner + "/" + repo + "/builds?" + query
            );
          }
        },
        {
          key: "getBuild",
          value: function getBuild(owner, repo, number) {
            return this._get(
              "/api/repos/" + owner + "/" + repo + "/builds/" + number
            );
          }
        },
        {
          key: "getBuildFeed",
          value: function getBuildFeed(opts) {
            var query = encodeQueryString(opts);
            return this._get("/api/user/feed?" + query);
          }
        },
        {
          key: "cancelBuild",
          value: function cancelBuild(owner, repo, number, ppid) {
            return this._delete(
              "/api/repos/" +
                owner +
                "/" +
                repo +
                "/builds/" +
                number +
                "/" +
                ppid
            );
          }
        },
        {
          key: "approveBuild",
          value: function approveBuild(owner, repo, build) {
            return this._post(
              "/api/repos/" +
                owner +
                "/" +
                repo +
                "/builds/" +
                build +
                "/approve"
            );
          }
        },
        {
          key: "declineBuild",
          value: function declineBuild(owner, repo, build) {
            return this._post(
              "/api/repos/" +
                owner +
                "/" +
                repo +
                "/builds/" +
                build +
                "/decline"
            );
          }
        },
        {
          key: "restartBuild",
          value: function restartBuild(owner, repo, build, opts) {
            var query = encodeQueryString(opts);
            return this._post(
              "/api/repos/" +
                owner +
                "/" +
                repo +
                "/builds/" +
                build +
                "?" +
                query
            );
          }
        },
        {
          key: "getLogs",
          value: function getLogs(owner, repo, build, proc) {
            return this._get(
              "/api/repos/" + owner + "/" + repo + "/logs/" + build + "/" + proc
            );
          }
        },
        {
          key: "getArtifact",
          value: function getArtifact(owner, repo, build, proc, file) {
            return this._get(
              "/api/repos/" +
                owner +
                "/" +
                repo +
                "/files/" +
                build +
                "/" +
                proc +
                "/" +
                file +
                "?raw=true"
            );
          }
        },
        {
          key: "getArtifactList",
          value: function getArtifactList(owner, repo, build) {
            return this._get(
              "/api/repos/" + owner + "/" + repo + "/files/" + build
            );
          }
        },
        {
          key: "getGlobalSecretList",
          value: function getGlobalSecretList() {
            return this._get("/api/globalsecrets");
          }
        },
        {
          key: "createGlobalSecret",
          value: function createGlobalSecret(secret) {
            return this._post(
              "/api/globalsecrets",
              secret
            );
          }
        },
        {
          key: "deleteGlobalSecret",
          value: function deleteGlobalSecret(secret) {
            return this._delete(
              "/api/globalsecrets/" + secret
            );
          }
        },
        {
          key: "getSecretList",
          value: function getSecretList(owner, repo) {
            return this._get("/api/repos/" + owner + "/" + repo + "/secrets");
          }
        },
        {
          key: "createSecret",
          value: function createSecret(owner, repo, secret) {
            return this._post(
              "/api/repos/" + owner + "/" + repo + "/secrets",
              secret
            );
          }
        },
        {
          key: "deleteSecret",
          value: function deleteSecret(owner, repo, secret) {
            return this._delete(
              "/api/repos/" + owner + "/" + repo + "/secrets/" + secret
            );
          }
        },
        {
          key: "getRegistryList",
          value: function getRegistryList(owner, repo) {
            return this._get("/api/repos/" + owner + "/" + repo + "/registry");
          }
        },
        {
          key: "createRegistry",
          value: function createRegistry(owner, repo, registry) {
            return this._post(
              "/api/repos/" + owner + "/" + repo + "/registry",
              registry
            );
          }
        },
        {
          key: "deleteRegistry",
          value: function deleteRegistry(owner, repo, address) {
            return this._delete(
              "/api/repos/" + owner + "/" + repo + "/registry/" + address
            );
          }
        },
        {
          key: "getSelf",
          value: function getSelf() {
            return this._get("/api/user");
          }
        },
        {
          key: "getToken",
          value: function getToken() {
            return this._post("/api/user/token");
          }
        },
        {
          key: "on",
          value: function on(callback) {
            return this._subscribe("/stream/events", callback, {
              reconnect: true
            });
          }
        },
        {
          key: "stream",
          value: function stream(owner, repo, build, proc, callback) {
            return this._subscribe(
              "/stream/logs/" + owner + "/" + repo + "/" + build + "/" + proc,
              callback,
              {
                reconnect: false
              }
            );
          }
        },
        {
          key: "_get",
          value: function _get(path) {
            return this._request("GET", path, null);
          }
        },
        {
          key: "_post",
          value: function _post(path, data) {
            return this._request("POST", path, data);
          }
        },
        {
          key: "_patch",
          value: function _patch(path, data) {
            return this._request("PATCH", path, data);
          }
        },
        {
          key: "_delete",
          value: function _delete(path) {
            return this._request("DELETE", path, null);
          }
        },
        {
          key: "_subscribe",
          value: function _subscribe(path, callback, opts) {
            var query = encodeQueryString({
              access_token: this.token
            });
            path = this.server ? this.server + path : path;
            path = this.token ? path + "?" + query : path;

            var events = new EventSource(path);
            events.onmessage = function(event) {
              var data = JSON.parse(event.data);
              callback(data);
            };
            if (!opts.reconnect) {
              events.onerror = function(err) {
                if (err.data === "eof") {
                  events.close();
                }
              };
            }
            return events;
          }
        },
        {
          key: "_request",
          value: function _request(method, path, data) {
            var endpoint = [this.server, path].join("");
            var xhr = new XMLHttpRequest();
            xhr.open(method, endpoint, true);
            if (this.token) {
              xhr.setRequestHeader("Authorization", "Bearer " + this.token);
            }
            if (method !== "GET" && this.csrf) {
              xhr.setRequestHeader("X-CSRF-TOKEN", this.csrf);
            }
            return new Promise(
              function(resolve, reject) {
                xhr.onload = function() {
                  if (xhr.readyState === 4) {
                    if (xhr.status >= 300) {
                      var error = {
                        status: xhr.status,
                        message: xhr.response
                      };
                      if (this.onerror) {
                        this.onerror(error);
                      }
                      reject(error);
                      return;
                    }
                    var contentType = xhr.getResponseHeader("Content-Type");
                    if (
                      contentType &&
                      contentType.startsWith("application/json")
                    ) {
                      resolve(JSON.parse(xhr.response));
                    } else {
                      resolve(xhr.response);
                    }
                  }
                }.bind(this);
                xhr.onerror = function(e) {
                  reject(e);
                };
                if (data) {
                  xhr.setRequestHeader("Content-Type", "application/json");
                  xhr.send(JSON.stringify(data));
                } else {
                  xhr.send();
                }
              }.bind(this)
            );
          }
        }
      ],
      [
        {
          key: "fromEnviron",
          value: function fromEnviron() {
            return new DroneClient(
              process && process.env && process.env.DRONE_SERVER,
              process && process.env && process.env.DRONE_TOKEN,
              process && process.env && process.env.DRONE_CSRF
            );
          }
        },
        {
          key: "fromWindow",
          value: function fromWindow() {
            return new DroneClient(
              window && window.DRONE_SERVER,
              window && window.DRONE_TOKEN,
              window && window.DRONE_CSRF
            );
          }
        }
      ]
    );

    return DroneClient;
  })();

  exports.default = DroneClient;

  /**
   * Encodes the values into url encoded form sorted by key.
   *
   * @param {object} query parameters in key value object.
   * @return {string} query parameter string
   */
  var encodeQueryString = (exports.encodeQueryString = function encodeQueryString() {
    var params =
      arguments.length > 0 && arguments[0] !== undefined ? arguments[0] : {};

    return params
      ? Object.keys(params)
          .sort()
          .map(function(key) {
            var val = params[key];
            return encodeURIComponent(key) + "=" + encodeURIComponent(val);
          })
          .join("&")
      : "";
  });
});
