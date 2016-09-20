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
	    test: /\.js$/,
	    exclude: /node_modules/,
	    loader: 'babel'
	}
	]
    },
    plugins: [
	new webpack.optimize.CommonsChunkPlugin("vendor", "vendor.js")
    ]
};
