# rcc

<p>
  <img src='https://rcc.evertras.com/api/v0/badge/coverage?key=github.com/Evertras/rcc' alt='Coverage Status'/>
</p>

A simple remote code coverage tool for self-hosting.  The above badge
is using this service!

## Why does this exist?

You may have some privately hosted repositories.  You may not want these
repositories to ever be known to the outside world, but you still want
to have some sort of code coverage badge for them.  This is intended
as a simple-as-possible self-hosted store for code coverage badges.

## What it can do

### Get/set latest code coverage

For a given key (usually the repository URL), set the latest code coverage.
You can also retrieve it as a raw value for any purpose.

### Generate coverage badges

You can query the service for a badge based on the main branch's last reported
total coverage.  Custom thresholds can be specified for red/orange/green colors.

## What is stored

This very simply stores the main branch code coverage for a given repository
as the key.  A few different options are available for storage.

When using the standalone server binary, you can supply the storage type as a flag
or as an environment variable.

```bash
# Build the standalone server
make bin/rcc

# Run the standalone server with file storage mode as flag
./bin/rcc --storage-type file

# Run it with environment variable
RCC_STORAGE_TYPE=file ./bin/rcc
```

See `rcc --help` for more configuration options.

### In Memory

`--storage-type in-memory`

A simple in-memory repository is default for testing.  This should not actually
be used anywhere besides for quick testing purposes.

### Local File

`--storage-type file`

Stores all key/values in a local directory that can be supplied with config.
Useful for running in a non-AWS setup.

### DynamoDB

`--storage-type dynamodb`

Stores data in DynamoDB because it's cheap and easy.  The following schema is used.

| Attribute | Type | Details |
|-----------|------|---------|
| Key       | S    | The key for the coverage data.  Generally the full repository URL, such as `github.com/Evertras/rcc`. |
| Value1000 | N    | The value of the percent coverage for the key, multiplied by 10 and stored as an integer for precision. |
| LastUpdated | S  | The timestamp of when this coverage was last set. |

Note that this is only really supported in the lambda version for now.

## How it's deployed

Currently deployed to `https://rcc.evertras.com/` via the [terraform](./terraform)
code provided here.  Runs as an AWS Lambda writing to DynamoDB purely for
cost purposes.  This could be deployed as a containerized instance somewhere,
but lambdas are pretty quick/cheap for something that doesn't need to be running
all the time.

## Auth

**Currently none.  Will need to add in the future to avoid abuse.**

## Endpoints

Note: `[key]` can be any `[a-zA-Z0-9\./-]` value and is intended to take a value
such as `github.com/Evertras/rcc`.  This key is case sensitive with a max length
of 64 characters.

### PUT /api/v1/coverage?key=[key]&value100=[value100]

Sets the coverage value.  `value100` is a number in the range of `[0.0, 100.0]%`
to the nearest 1 decimal point for internal storage.  The % sign is optional.

### GET /api/v1/coverage?key=[key]

Returns the coverage value as a plaintext value in the format `<0-100.0>%`.  Returns
rounded to the nearest decimal point.

### GET /api/v1/badge/coverage?key=[key]

Returns a code coverage badge that can be linked to in a readme.

This README's badge is set as follows:

```html
<p>
  <img src='https://rcc.evertras.com/api/v0/badge/coverage?key=github.com/Evertras/rcc' alt='Coverage Status'/>
</p>
```
