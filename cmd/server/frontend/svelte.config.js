const { typescript } = require('svelte-preprocess-esbuild');

module.exports = {
    preprocess: typescript({})
};
