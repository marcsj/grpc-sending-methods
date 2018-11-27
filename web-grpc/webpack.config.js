const path = require("path");
const { BundleAnalyzerPlugin } = require("webpack-bundle-analyzer");

module.exports = {
  context: path.resolve(__dirname, "src"),
  entry: "./index",
  mode: "production",
  devServer: {
    compress: true
  },
  plugins: [new BundleAnalyzerPlugin()]
};
