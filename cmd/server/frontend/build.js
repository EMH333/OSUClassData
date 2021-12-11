
const esbuild = require('esbuild');
const esbuildOptions = require('./esbuild.config');
const fs = require('fs');
const { createGzip, createBrotliCompress, constants } = require('zlib');
const { pipeline } = require('stream');
const {
  createReadStream,
  createWriteStream
} = require('fs');

function copyHTML() {
  fs.copyFile('./dev/index.html', './dist/index.html', (err) => {
    if (err) throw err;
  });
}

if (process.argv.length >= 2 && (process.argv[2] === "clean" || process.argv[2] === "production")) {
  const directory = './dist';

  if (fs.existsSync(directory)) {
    fs.rmdirSync(directory, { recursive: true });
  }

  if (!fs.existsSync(directory)) {
    fs.mkdirSync(directory);
  }
}

if (process.argv.length >= 2 && process.argv[2] === "serve") {
  let serveOptions = esbuildOptions;
  serveOptions.minify = false;
  copyHTML();
  esbuild.serve({
    port: 3000,
    servedir: './dist',
  }, serveOptions).then(server => {
    // Call "stop" on the server when you're done
    //server.stop()
    //process.exit(0)
  })
} else {
  esbuild.build(esbuildOptions).then(() => {
    if (process.argv.length >= 2 && process.argv[2] === "production") {
      compressJSandCSS();
    }
  }).catch(() => process.exit(1))
  copyHTML();
}

function compressJSandCSS() {
  if (!fs.existsSync("./dist/precompressed")) {
    fs.mkdirSync("./dist/precompressed", { recursive: true });
  }

  const gzip = createGzip({ level: constants.Z_MAX_LEVEL });
  const brotli = createBrotliCompress({
    params: {
      [constants.BROTLI_PARAM_QUALITY]: constants.BROTLI_MAX_QUALITY,
    }
  });


  compressFile('./dist/index.js', './dist/precompressed/index.js.gz', gzip);
  compressFile('./dist/index.js', './dist/precompressed/index.js.br', brotli);

  const gzip2 = createGzip({ level: constants.Z_MAX_LEVEL });
  const brotli2 = createBrotliCompress({
    params: {
      [constants.BROTLI_PARAM_QUALITY]: constants.BROTLI_MAX_QUALITY,
    }
  });

  compressFile('./dist/index.css', './dist/precompressed/index.css.gz', gzip2);
  compressFile('./dist/index.css', './dist/precompressed/index.css.br', brotli2);

}

function compressFile(input, output, type){
  let source = createReadStream(input);
  let destination = createWriteStream(output);
  pipeline(source, type, destination, (err) => {
    if (err) {
      console.error('An error occurred:', err);
      process.exitCode = 1;
    }
  });
}
