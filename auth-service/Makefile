# install dependencies
install:
	npm install

# run dev mode
run:
	npm run start:dev

lint:
	npm run lint

# build project
build:
	npm run build

# test
test:
	npm run test

# lint
lint:
	npm run lint

# format
format:
	npm run format

clean:
	rm -rf dist node_modules package-lock.json


# reinstall fresh
reinstall: clean install

# update packages (minor/patch)
update:
	npm update

# upgrade all packages to latest
upgrade:
	npx npm-check-updates -u
	npm install


# reset full project
reset: clean install build