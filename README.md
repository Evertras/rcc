# rcc

A simple remote code coverage tool for self-hosting.

## Why does this exist?

You may have some privately hosted repositories.  You may not want these
repositories to ever be known to the outside world, but you still want
to have some sort of code coverage badge for them.  This is intended
as a simple-as-possible self-hosted store for code coverage badges.

## What it can do

### Get/set latest code coverage

For a given key (usually the repository URL), set the latest code coverage.
You can also retrieve it as a raw value for any purpose.

### Generate badges

You can query the service for a badge based on the main branch's last reported
total coverage.  Custom thresholds can be specified for red/orange/green colors.

## What is stored

This very simply stores the main branch code coverage for a given repository
as the key.

For the very first ultra-simple version, this is just stored locally on disk.
In the future this will need to go to either a DB or S3 for better persistence,
preferably both as options depending on deployment.

## Auth

Currently none.  Will need to add in the future.

## Endpoints

Note: `<key>` can be any `[a-zA-Z0-9\./-]` value and is intended to take a value
such as `github.com/Evertras/rcc`.

### PUT /api/v1/coverage/value100/?repo=<key>

Sets the coverage value.  Accepts plain text of a number in the range of `[0, 100]`
to the nearest 2 decimal points for internal storage.  The % sign is optional.

### GET /api/v1/coverage/value100/?repo=<key>

Returns the coverage value as a plaintext value in the format `<0-100.0>%`.  Returns
rounded to the nearest decimal point.

### GET /api/v1/badge/coverage/?repo=<key>

Returns a code coverage badge that can be linked to in a readme.
