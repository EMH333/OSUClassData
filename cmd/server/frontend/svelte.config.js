const typescript = require('svelte-preprocess');

module.exports = {
    preprocess: typescript({}),
    compilerOptions: {
        //hydratable: true,
    }
};
