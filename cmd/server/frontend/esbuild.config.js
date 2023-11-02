import esbuildSvelte from "esbuild-svelte";
import svelteConfig from "./svelte.config.js";

const options = {
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
    loader: {
        '.woff2': 'file',
    },
    assetNames: '[dir]/[name]',
};

export default options;
