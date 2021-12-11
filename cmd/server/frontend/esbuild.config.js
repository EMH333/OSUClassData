const esbuildSvelte = require('esbuild-svelte');
const svelteConfig = require('./svelte.config');

module.exports = {
    entryPoints: ['./dev/index.ts'],
    format: 'esm',
    minify: true,
    bundle: true,
    splitting: true,
    outdir: './dist',
    sourcemap: 'external',
    plugins: [esbuildSvelte({
        preprocess: svelteConfig.preprocess,
    })],
};
