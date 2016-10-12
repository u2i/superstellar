const path = require('path');
const webpack = require('webpack');
const config = require('./webpackBase');

config.plugins.push(new webpack.optimize.DedupePlugin());
config.plugins.push(new webpack.optimize.UglifyJsPlugin());
config.plugins.push(new webpack.optimize.AggressiveMergingPlugin());

module.exports = config;
