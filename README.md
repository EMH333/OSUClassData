# OSU Class Data Explorer

[See the deployed version here](https://osuclassdata.ethohampton.com)

This is a web application that allows you to explore the raw data provided by **\_\_\_** acquired through a FOIA request for Oregon State University classes. The included data is between Fall 2019 and Summer 2021. Hopefully this will help your class forecasting in some way!

This application is still in development. This application is under a AGPLv3 license. For all you CS majors out there (or anyone who writes code), PRs are welcome!

## Development

Devlopment is on a local MySQL server, esbuild JS compiler and Golang server.

To get started:
```bash
# install JS dependencies and generate frontend
cd cmd/server/frontend
npm install
node build.js

# go back to root directory and start DB and dev server
# NOTE: This takes some time b/c it also preloads the DB tables
cd -
./scripts/startDB.sh && sleep 30 && ./scripts/startDB.sh
./scripts/devServer.sh
```
Once you see `Connected to Database!`, development can begin.

### File Structure

- `cmd/classParser` - A tool for parsing the raw data from the FOIA into SQL
- `cmd/server` - The Golang server
- `cmd/server/frontend` - The Svelte frontend
- `internal` - Utilities used by the Golang code
- `scripts` - Scripts I use during development and deployment
- `sql` - Initial SQL database schema and other useful SQL

### Development Methodology

I'm lazy and I like efficient things that I can work with quickly. That is why I choose Golang and Svelte as my primary languages/frameworks. I also did my best to keep the amount of external dependencies to a minimum. This allows me to iterate quickly.

The frontend sends a variety of GET requests to the backend, which returns JSON that is parsed by the frontend and displayed to the user. The backend currently doesn't do any caching since the DB server is on the same machine as the backend, and it works fast enough. There is some room for optimization of DB queries as well as frontend code.

The backend is not stateful and could easily be static files with a bit of work. I choose to increase complexity here to give this more room to grow in the future without having to re-architect everything.

### TODO

See `TODO.txt` in the root of this repository
