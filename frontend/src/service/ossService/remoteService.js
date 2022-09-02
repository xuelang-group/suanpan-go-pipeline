import { fetch } from 'whatwg-fetch'

let s = {};

let processResponse = function (response, type) {
  switch (type) {
    case 'arrayBuffer':
      return response.arrayBuffer();
    case 'blob':
      return response.blob();
    case 'json':
      return response.json();
  }
  return response.text();
};

s.get = function (url, responseType) {
  return new Promise(function (resolve, reject) {
    fetch(url)
      .then(function (response) {
        return processResponse(response, responseType);
      })
      .then(function (result) {
        resolve(result)
      })
      .catch(function (ex) {
        reject(ex);
      })
  });
};

s.post = function (url, data, contentType, responseType) {
  return new Promise(function (resolve, reject) {
    fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': contentType || 'application/octet-stream;'
      },
      body: JSON.stringify(data)
    })
      .then(function (response) {
        return processResponse(response, responseType);
      })
      .then(function (result) {
        resolve(result)
      })
      .catch(function (ex) {
        reject(ex);
      })
  });
};

s.head = function (url, responseType) {
  return new Promise(function (resolve, reject) {
    fetch(url, {
      method: 'HEAD'
    })
      .then(function (response) {
        return processResponse(response, responseType);
      })
      .then(function (result) {
        resolve(result)
      })
      .catch(function (ex) {
        reject(ex);
      })
  });
};

s.delete = function (url, responseType) {
  return new Promise(function (resolve, reject) {
    fetch(url, {
      method: 'DELETE'
    })
      .then(function (response) {
        return processResponse(response, responseType);
      })
      .then(function (result) {
        resolve(result)
      })
      .catch(function (ex) {
        reject(ex);
      })
  });
};

s.put = function (url, data, contentType, responseType) {
  return new Promise(function (resolve, reject) {
    fetch(url, {
      method: 'PUT',
      headers: {
        'Content-Type': contentType || "application/octet-stream;"
      },
      body: data,
      transformRequest: []
    })
      .then(function (response) {
        return processResponse(response, responseType);
      })
      .then(function (result) {
        resolve(result)
      })
      .catch(function (ex) {
        reject(ex);
      })
  });
};

export default s;

