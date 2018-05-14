const MiniCssExtractPlugin = require("mini-css-extract-plugin");
var FriendlyErrorsWebpackPlugin = require('friendly-errors-webpack-plugin');
var path = require('path');

module.exports = {
  mode: 'none',
  context: __dirname,
  entry: {
    initium : [
      './assets/js/app.js', './assets/stylesheets/app.scss'
    ]
  },
  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: '[name].js'
  },
  module: {
    rules: [
      {
        test: /\.html$/,
        use: [ 'html-loader' ]
      },
      {
        test: /(\.(png|jpe?g|gif)$|^((?!font).)*\.svg$)/,
        use: ['file-loader', 'img-loader']
      },
      {
        test: /(\.(woff2?|ttf|eot|otf)$|font.*\.svg$)/,
        use: ['file-loader']
      },
      {
        test: path.resolve(__dirname, 'assets/stylesheets/app.scss'),
        use: [
          MiniCssExtractPlugin.loader,
          'css-loader',
          {
            loader: 'postcss-loader',
            options: {
              plugins: [ require('autoprefixer') ]
            }
          },
          'resolve-url-loader',
          'sass-loader'
        ]
      },
      {
        test: /\.css$/,
        use: [ 'css-loader' ]
      },
      {
        test: /\.s[ac]ss$/,
        exclude: [ path.resolve(__dirname, 'assets/stylesheets/app.scss')],
        use: [ 'css-loader', 'sass-loader' ]
      },
      {
        test: /\.less$/,
        use: [ 'css-loader', 'less-loader' ]
      }
    ]
  },
  plugins: [
    new FriendlyErrorsWebpackPlugin(),
    new MiniCssExtractPlugin('/stylesheet/app.css')
  ],
  performance: { hints: false },
  devtool: false,
}