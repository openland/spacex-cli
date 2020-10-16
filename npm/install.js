const os = require('os');
const fs = require('fs');

const type = os.type();
const arch = os.arch();

if (arch !== 'x64') {
    throw new Error(`Only x64 is supported`);
}

if (type === 'Darwin') {
    fs.copyFileSync(__dirname + '/bin/spacex-cli-macos', __dirname + '/bin/spacex-cli')
} else if (type === 'Linux') {
    fs.copyFileSync(__dirname + '/bin/spacex-cli-linux', __dirname + '/bin/spacex-cli')
} else {
    throw new Error(`Only macos or linux are supported`);
}