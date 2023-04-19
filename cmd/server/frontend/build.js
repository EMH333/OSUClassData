
const esbuild = require('esbuild');
const esbuildOptions = require('./esbuild.config');
const fs = require('fs');

function copyHTML() {
  const folders = [
    'img',
  ];

  for (let i = 0; i < folders.length; i++) {
    if (!fs.existsSync(`./dist/${folders[i]}`)) {
      fs.mkdirSync(`./dist/${folders[i]}`);
    }
  }

  const files = [
    'index.html',
    'about.html',
    'subject.html',
    'favicon.ico',
    'img/logo.png',
  ];

  for (let i = 0; i < files.length; i++) {
    fs.copyFileSync(`./dev/${files[i]}`, `./dist/${files[i]}`);
  }

  //copy the templates
  const templateFiles = [
    'class.html',
    'leaderboard.html',
    'sitemap.html',
  ];

  for (let i = 0; i < templateFiles.length; i++) {
    fs.copyFileSync(`./templates/${templateFiles[i]}`, `./distTemplates/${templateFiles[i]}`);
  }
}

const directories = ['./dist', './distSSR', './distTemplates'];
let compileOptions = esbuildOptions;

if (process.argv[2] === "production") {
  compileOptions.pure = ['console.log'];
}

// remove build directory if building clean or for production
if (process.argv.length >= 2 && (process.argv[2] === "clean" || process.argv[2] === "production")) {
  directories.forEach(directory => {
    if (fs.existsSync(directory)) {
      fs.rmSync(directory, { recursive: true });
    }
  });
}


// make sure dist exists
directories.forEach(directory => {
  if (!fs.existsSync(directory)) {
    fs.mkdirSync(directory);
  }
});

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
  ssr();
  //allow for non-minified code
  if (process.argv.length >= 2 && process.argv[2] === "dev") { compileOptions.minify = false; compileOptions.watch = true; }

  //allow for non-minified code but no watching
  if (process.argv.length >= 2 && process.argv[2] === "ci") { compileOptions.minify = false; }

  esbuild.build(compileOptions)
    .then(output => {
      //fs.writeFileSync('./dist/metafile.json', JSON.stringify(output.metafile));

      for (file in output.metafile.outputs) {
        let fileInfo = output.metafile.outputs[file];
        switch (file) {
          case "dist/index.js":
            insertPreload('./dist/index.html', fileInfo.imports);
            break;

          case "dist/class.js":
            //insertPreload('./dist/class.html', fileInfo.imports);
            insertPreload('./distTemplates/class.html', fileInfo.imports);
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
        //don't include map files
        if (file.endsWith(".map")) { continue; }
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
      header += `<link rel="preload" href="/${fileName}" as="script">`;
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

function ssr() {
  // TODO also render the class and subject pages
  const esbuildSvelte = require('esbuild-svelte');
  const svelteOptions = require("./svelte.config");
  esbuild
    .build({
      ...esbuildOptions,
      entryPoints: ["ssr.js"],
      outdir: "./distSSR",
      format: "cjs",
      splitting: false,
      plugins: [
        esbuildSvelte({
          ...svelteOptions,
          compilerOptions: {
            generate: "ssr",
          },
        }),
      ],
    })
    .then(() => {
      //now we can generate the html
      const fs = require("fs");
      const output = require("./distSSR/ssr");

      const initialHTML = fs.readFileSync("./dist/index.html");
      let rendered = output.render({
        target: "document.body"
      });
      if (rendered.head != "") {
        console.error("Head is not empty, this is not supported");
      }
      //console.log(rendered.html)
      let final = initialHTML.toString().replace("<!--ssr-html-->", rendered.html)

      fs.writeFileSync("./dist/index.html", final);
    })
    .catch((err) => { console.error(err); process.exit(1) });
}
