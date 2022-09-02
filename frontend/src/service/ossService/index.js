import remoteService from './remoteService';
const OSS = require('ali-oss');
import { getStorageInfo } from '../index'

/**
 * @return {string}
 */
let Utf8ArrayToStr = function (array) {
  let out, i, len, c;
  let char2, char3;

  out = "";
  len = array.length;
  i = 0;
  while (i < len) {
    c = array[i++];
    switch (c >> 4) {
      case 0: case 1: case 2: case 3: case 4: case 5: case 6: case 7:
        // 0xxxxxxx
        out += String.fromCharCode(c);
        break;
      case 12: case 13:
        // 110x xxxx   10xx xxxx
        char2 = array[i++];
        out += String.fromCharCode(((c & 0x1F) << 6) | (char2 & 0x3F));
        break;
      case 14:
        // 1110 xxxx  10xx xxxx  10xx xxxx
        char2 = array[i++];
        char3 = array[i++];
        out += String.fromCharCode(((c & 0x0F) << 12) |
          ((char2 & 0x3F) << 6) |
          ((char3 & 0x3F) << 0));
        break;
    }
  }

  return out;
};


/*
* options = {
*   type: '',
*   tokenPath: '',
*   region: '',
*   bucket: ''
* }
*
* */
export default function init(options) {
  let _stsToken = {};
  let _oss;
  var loading = false;
  var deferredList = [];

  options.tokenServer = options.tokenServer || '';
  options.tokenPath = options.tokenPath || '/oss/token';
  let isClient = !!options.client

  let fetchToken = function(isClient, tokenUrl) {
    if(isClient) {
      return getStorageInfo()
    }else {
      return remoteService.get(tokenUrl)
    }
  }

  let getStsToken = function () {
    return new Promise(function (resolve, reject) {
      if (_stsToken
        && _stsToken.Credentials
        && _stsToken.Credentials.Expiration
        && (new Date(_stsToken.Credentials.Expiration).getTime() - (new Date()).getTime() > 1800000)) {
        resolve({
          token: _stsToken,
          refresh: false
        });
      }
      else {
        if (loading) {
          deferredList.push({
            resolve: resolve,
            reject: reject
          });
        }
        else {
          loading = true;
          fetchToken(isClient, options.tokenServer + options.tokenPath, 'json')
            .then(function (token) {
              _stsToken = token;
              resolve({
                token: token,
                refresh: true
              });
              while (deferredList.length) {
                var d = deferredList.pop();
                d.resolve({
                  token: token,
                  refresh: true
                });
              }
            })
            .catch(function (e) {
              reject(e);
              while (deferredList.length) {
                var d = deferredList.pop();
                d.reject(e);
              }
            })
            .finally(function () {
              loading = false;
            });
        }
      }
    });
  };

  let getOss = function () {
    return new Promise(function (resolve, reject) {
      getStsToken().then(function (result) {
        let token = result.token;
        if (!_oss || result.refresh) {
          _oss = new OSS({
            region: token.region,
            accessKeyId: token.Credentials.AccessKeyId,
            accessKeySecret: token.Credentials.AccessKeySecret,
            stsToken: token.Credentials.SecurityToken,
            bucket: token.bucket
          });

          resolve(_oss);
        }
        else {
          resolve(_oss);
        }
      }, function (err) {
        reject(err);
      });
    });
  };

  let upload = function (file, key, progress, error, complete) {
    return new Promise(function (resolve, reject) {
      getOss().then(function (oss) {
        let options = {
          progress,
          partSize: 500 * 1024,
          timeout: 60000
        };
        oss.multipartUpload(key, file, options).then((res) => {
          complete(res)
          // currentCheckpoint = null;
          // uploadFileClient = null;
        }).catch((err) => {
          if (oss.isCancel()) {
          }
          else {
            if (err.name.toLowerCase().indexOf('connectiontimeout') !== -1) {
              // timeout retry
              // if (retryCount < retryCountMax) {
              //   retryCount++;
              //   console.error(`retryCount : ${retryCount}`);
              //   upload('');
              // }
            }
          }
        });
        resolve();
      }, function (err) {
        reject(err);
      });
    });
  };

  let getObject = function (key, isBinary) {
    return new Promise(function (resolve, reject) {
      getOss().then(function (oss) {
        oss.get(key).then(function (res) {
          if (isBinary) {
            return resolve(res.content);
          }
          resolve(Utf8ArrayToStr(res.content));
        }).catch(function (err) {
          reject(err);
        })
      }, function (err) {
        reject(err);
      });
    });
  };

  let deleteObject = function (key) {
    return new Promise(function (resolve, reject) {
      getOss().then(function (oss) {
        oss.delete(key).then(function (res) {
          resolve(res)
        }, function (err) {
          reject(err)
        })
      }, function (err) {
        reject(err);
      });

    });
  };

  let listObject = function (key) {
    return new Promise(function (resolve, reject) {
      getOss().then(function (oss) {
        oss.list({
          'prefix': key,
          'max-keys': 1000
        }).then(function (result) {
          resolve(result)
        })
      }, function (err) {
        reject(err);
      });
    });
  };

  let getSignedUrl = function (key) {
    return new Promise(function (resolve, reject) {
      getOss().then(function (oss) {
        resolve(oss.signatureUrl(key, { expires: 3600 }));
      }, function (err) {
        reject(err);
      });
    });
  };

  let putObject = function (key, data) {
    return new Promise(function (resolve, reject) {
      getOss().then(function (oss) {
        oss.put(key, data)
          .then((result) => {
            resolve(result)
          })
          .catch(function (e) {
            reject(e);
          })
      }, function (err) {
        reject(err);
      });
    });
  };

  let getObject2 = function (key, isBinary = false) {
    return new Promise(function (resolve, reject) {
      remoteService.post(options.tokenServer + '/oss/object/get', { Key: key }, 'application/json', 'json')
        .then(function (res) {
          remoteService.get(res.data, isBinary ? 'arrayBuffer' : 'text')
            .then(function (r) {
              if (isBinary) {
                resolve(new Uint8Array(r));
                return;
              }
              resolve(r);
            }, function (err) {
              reject(err);
            });
        }, function (err) {
          reject(err);
        });
    });
  };

  let deleteObject2 = function (key) {
    return new Promise(function (resolve, reject) {
      remoteService.post(options.tokenServer + '/oss/object/delete', {
        Key: key
      }, 'application/json', 'json')
        .then(function (res) {
          remoteService.delete(res.data, 'text')
            .then(function (r) {
              resolve(r);
            }, function (err) {
              reject(err);
            });
        }, function (err) {
          reject(err);
        });
    });
  };

  let listObject2 = function (key) {
    return new Promise(function (resolve, reject) {
      remoteService.post(options.tokenServer + '/oss/object/list', {
        Key: key
      }, 'application/json', 'json')
        .then(function (res) {
          resolve(res.data);
        }, function (err) {
          reject(err);
        });
    });
  };

  let headObject2 = function (key) {
    return new Promise(function (resolve, reject) {
      remoteService.post(options.tokenServer + '/oss/object/head', {
        Key: key
      }, 'application/json', 'json')
        .then(function (res) {
          res.NextAppendPosition = res.ContentLength;
          resolve(res);
        }, function (err) {
          reject(err);
        });
    });
  };

  let getSignedUrl2 = function (key, isBinary) {
    return new Promise(function (resolve, reject) {
      remoteService.post(options.tokenServer + '/oss/object/get', {
        Key: key
      }, 'application/json', 'json')
        .then(function (res) {
          resolve(res.data);
        }, function (err) {
          reject(err);
        });
    });
  };

  let putObject2 = function (key, body) {
    return new Promise(function (resolve, reject) {
      remoteService.post(options.tokenServer + '/oss/object/put', {
        Key: key,
        Expires: 0,
        CacheControl: 'no-cache, must-revalidate'
      }, 'application/json', 'json')
        .then(function (res) {
          if (!res.success) {
            reject(res);
            return;
          }
          if (res.data) {
            remoteService.put(res.data, body)
              .then(function (res) {
                resolve(res);
              }, function (err) {
                reject(err);
              });
          }
        }, function (err) {
          reject(err);
        }).catch(err => {
          reject(err);
        });
    });
  };

  let upload2 = function (file, key, onprogress, onerror, oncomplete) {
    return new Promise(function (resolve, reject) {
      remoteService.post(options.tokenServer + '/oss/object/put', {
        Key: key
      }, 'application/json', 'json')
        .then(function (res) {
          let xhr = new XMLHttpRequest();
          xhr.onload = function () {
            if (xhr.status === 200) {
              oncomplete({
                res: {
                  status: 200,
                  statusCode: 200,
                  size: 0
                },
                bucket: "suanpan",
                name: key,
              });
            }
            else {
              onerror(xhr.statusText);
            }
          };
          xhr.onprogress = function (evt) {
            onprogress(evt.loaded / evt.total);
          };
          xhr.open('PUT', res.data, true);
          xhr.send(file);

          resolve(res);
        }, function (err) {
          reject(err);
        });
    });
  };

  let s = {
    getOss: getOss,
    upload: upload,
    getObject: getObject,
    putObject: putObject,
    deleteObject: deleteObject,
    listObject: listObject2,
    getSignedUrl: getSignedUrl
  };
  if (options.type !== 'oss') {
    s.getObject = getObject2;
    s.putObject = putObject2;
    s.headObject = headObject2;
    s.upload = upload2;
    s.deleteObject = deleteObject2;
    s.listObject = listObject2;
    s.getSignedUrl = getSignedUrl2;
  }

  return s;
};
