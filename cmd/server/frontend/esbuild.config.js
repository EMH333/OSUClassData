const esbuildSvelte = require('esbuild-svelte');
const svelteConfig = require('./svelte.config');

module.exports = {
    entryPoints: ['./dev/index.ts', './dev/class.ts', './dev/subject.ts'],
    format: 'esm',
    minify: true,
    bundle: true,
    splitting: true,
    metafile: true,
    outdir: './dist',
    sourcemap: 'external',
    mainFields: ["svelte", "browser", "module", "main"],
    plugins: [esbuildSvelte(svelteConfig)],
};
