const webpack = require('webpack');

module.exports = function override(config, env) {
  // Add polyfills for Node.js modules
  config.resolve.fallback = {
    ...config.resolve.fallback,
    "assert": require.resolve("assert/"),
    "buffer": require.resolve("buffer/"),
    "process": require.resolve("process/browser"),
    "crypto": require.resolve("crypto-browserify"),
    "stream": require.resolve("stream-browserify"),
    "util": require.resolve("util/"),
    "url": require.resolve("url/"),
    "os": require.resolve("os-browserify/browser"),
    "path": require.resolve("path-browserify"),
    "fs": false,
    "net": false,
    "tls": false,
  };

  // Add plugins to provide global variables
  config.plugins.push(
    new webpack.ProvidePlugin({
      process: 'process/browser',
      Buffer: ['buffer', 'Buffer'],
    })
  );

  return config;
};
