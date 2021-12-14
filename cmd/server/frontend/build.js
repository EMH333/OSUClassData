
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

  fs.copyFile('./dev/class.html', './dist/class.html', (err) => {
    if (err) throw err;
  });

  fs.copyFile('./dev/about.html', './dist/about.html', (err) => {
    if (err) throw err;
  });
}

if (process.argv.length >= 2 && (process.argv[2] === "clean" || process.argv[2] === "production")) {
  const directory = './dist';

  if (fs.existsSync(directory)) {
    fs.rmSync(directory, { recursive: true });
  }

  if (!fs.existsSync(directory)) {
    fs.mkdirSync(directory);
  }
}

copyHTML();

if (process.argv.length >= 2 && process.argv[2] === "serve") {
  let serveOptions = esbuildOptions;
  serveOptions.minify = false;

  esbuild.serve({
    port: 3000,
    servedir: './dist',
  }, serveOptions).then(server => {
    // Call "stop" on the server when you're done
    //server.stop()
    //process.exit(0)
  })
} else {
  let compileOptions = esbuildOptions;
  //allow for non-minified code
  if (process.argv.length >= 2 && process.argv[2] === "dev") { compileOptions.minify = false; compileOptions.watch = true; }
  
  esbuild.build(compileOptions).then(() => {
    if (process.argv.length >= 2 && process.argv[2] === "production") {
      compressJSandCSS();
    }
  }).catch((err) => {console.error(err); process.exit(1)} );
}

function compressJSandCSS() {
  if (!fs.existsSync("./dist/precompressed")) {
    fs.mkdirSync("./dist/precompressed", { recursive: true });
  }

  fs.readdirSync("./dist").forEach(file => {
    if (file.endsWith(".js") || file.endsWith(".css") || file.endsWith(".html")) {
      compressFile("./dist/" + file, "./dist/precompressed/" + file +".gz");
    }
  });
}

function compressFile(input, output){
  const gzip = createGzip({ level: constants.Z_MAX_LEVEL });

  let source = createReadStream(input);
  let destination = createWriteStream(output);
  pipeline(source, gzip, destination, (err) => {
    if (err) {
      console.error('An error occurred:', err);
      process.exitCode = 1;
    }
  });
}
