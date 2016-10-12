const path = require('path');
const webpack = require('webpack');
const config = require('./webpackBase');


config.debug = true;
config.devtool = 'eval-source-map';

module.exports = config;
