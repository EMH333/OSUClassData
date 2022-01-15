
const esbuild = require('esbuild');
const esbuildOptions = require('./esbuild.config');
const fs = require('fs');

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

  fs.copyFile('./dev/subject.html', './dist/subject.html', (err) => {
    if (err) throw err;
  });

  fs.copyFile('./dev/icon.ico', './dist/favicon.ico', (err) => {
    if (err) throw err;
  });
}

const directory = './dist';
let compileOptions = esbuildOptions;

if (process.argv[2] === "production") {
  compileOptions.pure = ['console.log'];
}

// remove build directory if building clean or for production
if (process.argv.length >= 2 && (process.argv[2] === "clean" || process.argv[2] === "production")) {
  if (fs.existsSync(directory)) {
    fs.rmSync(directory, { recursive: true });
  }
}


// make sure dist exists
if (!fs.existsSync(directory)) {
  fs.mkdirSync(directory);
}

//copy all the html files
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
  //allow for non-minified code
  if (process.argv.length >= 2 && process.argv[2] === "dev") { compileOptions.minify = false; compileOptions.watch = true; }

  //allow for non-minified code but no watching
  if (process.argv.length >= 2 && process.argv[2] === "ci") { compileOptions.minify = false;}

  esbuild.build(compileOptions)
    .catch((err) => { console.error(err); process.exit(1) });
}
