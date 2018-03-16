const path = require('path');
const webpack = require('webpack');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
    entry: {
        app: ["babel-polyfill", "./js/game.js"],
        vendor: ['protobufjs', 'pixi.js']
    },
    output: {
        filename: "[name].[chunkhash].js",
        chunkFilename: "[id].[chunkhash].js"
    },
    module: {
        preLoaders: [
            {
                test: /\.js$/,
                loader: "eslint-loader",
                exclude: /node_modules/
            }
        ],
        loaders: [
            {
                test: /\.json$/,
                loader: 'json'
            }, {
                test: /\.js$/,
                exclude: /node_modules/,
                loader: 'babel'
            }, {
                test: path.resolve(__dirname, 'node_modules', 'pixi.js'),
                loader: 'ify'
            }, {
                test: /\.scss$/,
                loader: "style!css!sass?outputStyle=expanded&includePaths[]=" + path.resolve(__dirname, "./node_modules/compass-mixins/lib")
            }]
    },
    plugins: [
        new webpack.optimize.CommonsChunkPlugin("vendor", "vendor.[hash].js"),
        new HtmlWebpackPlugin({
            title: 'Superstellar',
            template: 'index.ejs'
        }),
        new webpack.DefinePlugin({
            'BACKEND_HOST': JSON.stringify(process.env.BACKEND_HOST || 'localhost'),
            'BACKEND_PORT': JSON.stringify(process.env.BACKEND_PORT || '80'),
            '__DEBUG__': process.env.DEBUG

        })
    ],
    postLoaders: [{
        test: /\.js$/,
        include: path.resolve(__dirname, 'node_modules/pixi.js'),
        loader: 'transform/cacheable?brfs'
    }]
};
