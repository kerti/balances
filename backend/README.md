# Balances Backend

[![codecov](https://codecov.io/gh/kerti/balances/graph/badge.svg?token=FTH30BY4MN)](https://codecov.io/gh/kerti/balances)

## How To Test

### Test Pre-requisite

To generate mock objects for testing, you need to install mockgen:

```bash
go install github.com/golang/mock/mockgen@v1.6.0
```

After installing, check if the mockgen executable is accessible:

```bash
❯ which mockgen
/Users/admin/go/bin/mockgen
```

If the mockgen is not found, make sure to include to go bin path in your PATH.
If you're using zsh, edit your `~/.zshrc` file and add the following line:

```
export PATH=$(go env GOPATH)/bin:$PATH
```

After that, reload the session configuration:

```
source ~/.zshrc
```

### Generating Mocks

To generate the mock objects, run:

```
make genmock
```

You can see what genmock is doing behind the scenes by reading the `Makefile`.

## How To Run

1. Clone the repository.
2. Prepare the database.
   - Create a new database
   - Run the SQL scripts located in `/migrations`
   - If you would like to reset the data, run the SQL scripts in `/migrations/demo`
3. Prepare the server
   - Make a copy of `.env.example` and name it `.env`
   - Setup your environment variables to point to the database you've setup in step 2 above
4. Run the service.
   ```bash
   make local
   ```