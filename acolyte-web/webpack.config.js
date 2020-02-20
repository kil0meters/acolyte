const path = require('path');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const TerserJSPlugin = require('terser-webpack-plugin');
const OptimizeCSSAssetsPlugin = require('optimize-css-assets-webpack-plugin');
const BrotliPlugin = require('brotli-webpack-plugin');
const CompressionPlugin = require('compression-webpack-plugin');

const libraryName = 'acolyte';

module.exports = (env, argv) => ({
    optimization: {
        minimize: argv.mode === 'production',
        minimizer: [new TerserJSPlugin({}), new OptimizeCSSAssetsPlugin({})]
    },
    plugins: [
        new MiniCssExtractPlugin({
            filename: 'bundle.css',
            chunkFilename: 'chunk.css',
        }),
        new BrotliPlugin({
            asset: '[path].br[query]',
            test: /\.(js|css|html|svg|eot|woff|woff2|ttf)$/,
            threshold: 10240,
            minRatio: 0.8
        }),
        new CompressionPlugin({
            filename: '[path].gz[query]',
            test: /\.(js|css|html|svg|eot|woff|woff2|ttf)$/,
            threshold: 10240,
            minRatio: 0.8
        })
    ],
    mode: 'development',
    entry: './src/index.ts',
    output: {
        filename: 'bundle.js',
        library: libraryName,
        path: path.resolve(__dirname, 'dist'),
    },
    resolve: {
        extensions: ['.ts', '.js', '.json']
    },
    module: {
        rules: [
            {
                test: /\.ts$/,
                use: 'ts-loader',
                exclude: /node_modules/,
            },
            {
                test: /\.(eot|woff|woff2|svg|ttf)([?]?.*)$/,
                loader: 'file-loader',
                options: {
                    name(file) {
                        if (argv.mode !== 'production') {
                            return '[name].[ext]';
                        }

                        return '[contenthash].[ext]';
                    }
                }
            },
            {
                test: /\.css$/i,
                use: [MiniCssExtractPlugin.loader, 'css-loader'],
            },
            {
                test: /\.(png|jpe?g|gif)$/i,
                use: [
                    {
                        loader: 'responsive-loader',
                        options: {
                            name: (argv.mode !== 'production') ? '[name].[ext]' : '[contenthash].[ext]',
                            sizes: [128],
                            adapter: require('responsive-loader/sharp'),
                        }
                    },
                    {
                        loader: 'image-webpack-loader',
                        options: {
                            optipng: {
                                enabled: argv.mode === 'production',
                                optimizationLevel: 7,
                            },
                            pngquant: {
                                quality: [0.65, 0.90],
                                speed: (argv.mode === 'production') ? 1 : 11
                            },
                        },
                    }
                ],
            },
        ],
    }
});