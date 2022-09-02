import { io } from "socket.io-client";

class SocketIOError extends Error {
  constructor(options) {
    const data = {
      message: '',
      traceback: '',
      ...options,
    };
    super(data.message);
    this.data = data;
  }

  get message() {
    return this.data.message;
  }

  get traceback() {
    return this.data.traceback;
  }
}

class SocketIOClient {
  constructor() {
    // this.socket = io('https://spnext.xuelangyun.com', {path: '/proxr/shanglu/63567/d0bbf000a8b411ec946ef39f840f6aa2/8888/socket.io'});
    this.socket = io({ path: new URL('./socket.io', location.href).pathname });
    this.timeout = 30 * 1000;
  }

  async request(options) {
    const opts = { data: {}, ...options };
    return new Promise((resolve, reject) => {
      const timer = setTimeout(
        () => reject(new SocketIOError({ message: 'Request timed out' })),
        this.timeout,
      );
      this.socket.emit(opts.event, opts.data, result => {
        clearTimeout(timer);
        if (result.success) {
          resolve(result.data);
        } else {
          const error = new SocketIOError(result.error);
          reject(error);
        }
      });
    });
  }

  emit(options) {
    const opts = { data: {}, ...options };
    this.socket.emit(opts.event, opts.data);
  }

  on(event, fn) {
    this.socket.on(event, fn);
    return this;
  }

  off(event, fn) {
    this.socket.off(event, fn);
    return this;
  }
}

export default new SocketIOClient();