const path = require('path');

module.exports = {
    entry: './gateway.js', // Your Gateway class file path
    output: {
        filename: 'gateway.bundle.js',
        path: path.resolve(__dirname, 'dist'),
        library: 'Gateway',
        libraryTarget: 'var'
    },
    module: {
        rules: [
            {
                test: /\.js$/,
                exclude: /node_modules/
            }
        ]
    }
};
