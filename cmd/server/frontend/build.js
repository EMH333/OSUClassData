
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
  if (process.argv.length >= 2 && process.argv[2] === "ci") { compileOptions.minify = false; }

  esbuild.build(compileOptions)
    .then(output => {
      for (file in output.metafile.outputs) {
        let fileInfo = output.metafile.outputs[file];
        switch (file) {
          case "dist/index.js":
            insertPreload('./dist/index.html', fileInfo.imports);
            break;

          case "dist/class.js":
            insertPreload('./dist/class.html', fileInfo.imports);
            break;

          case "dist/subject.js":
            insertPreload('./dist/subject.html', fileInfo.imports);
            break;
          default:
            break;
        }
      }

      //do some quick bundle calculations
      let bundleSize = 0;
      for (file in output.metafile.outputs) {
        bundleSize += output.metafile.outputs[file].bytes;
      }
      console.log(`Bundle size: ${(bundleSize / 1024).toFixed(1)} kb`);
    })
    .catch((err) => { console.error(err); process.exit(1) });
}

function generateLinkHeader(imports) {
  let header = '';
  for (let i = 0; i < imports.length; i++) {
    let fileName = imports[i].path.replace('dist/', '');
    if (fileName.startsWith("chunk")) { // only preload chunks, not async imports
      header += `<link rel="preload" href="${fileName}" as="script">`;
    }
  }
  return header;
}

function insertPreload(htmlPath, imports) {
  let html = fs.readFileSync(htmlPath, 'utf8');
  let headers = generateLinkHeader(imports);
  html = html.replace('<link ref="preloadReplace">', headers);
  fs.writeFileSync(htmlPath, html);
}
