const path = require('path');
const webpack = require('webpack');

module.exports = {
    entry: {
	app: ["babel-polyfill", "./js/game.js"],
	vendor: ['protobufjs']
    },
    output: {
	filename: "[name].js"
    },
    module: {
	loaders: [
    {
	    test: /\.json$/,
	    loader: 'json'
	},
	{
	    test: /\.js$/,
	    exclude: /node_modules/,
	    loader: 'babel'
	},
	{
	    test: path.resolve(__dirname, 'node_modules', 'pixi.js'),
	    loader: 'ify'
	}
	]
    },
    plugins: [
	new webpack.optimize.CommonsChunkPlugin("vendor", "vendor.js")
    ],
    postLoaders: [
	{
	    test: /\.js$/,
	    include: path.resolve(__dirname, 'node_modules/pixi.js'),
	    loader: 'transform/cacheable?brfs'
	}
    ],
    // TODO: configure production build
    debug: true,
    devtool: 'eval-source-map'
};
