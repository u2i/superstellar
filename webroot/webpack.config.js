var webpack = require('webpack');

module.exports = {
    entry: {
	app: ["babel-polyfill", "./js/game.js"]
    },
    output: {
	filename: "bundle.js"
    },
    module: {
	loaders: [
	    {
	    test: /\.js$/,
	    exclude: /node_modules/,
	    loader: 'babel'
	}
	]
    }
};
