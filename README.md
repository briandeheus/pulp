# Pulp

File Caching Proxy for files on Google Cloud Storage

## Introduction

Setting up Google Cloud Storage to serve static files can be a pain in the behind. Pulp should make this a bit easier.

### TODO
- [ ] Add tests...
- [ ] Implement `PULP_TOKEN` support
- [ ] Use Google PubSub to automatically refresh cache
- [ ] Cache 404s to prevent hitting GCS with useless requests
- [ ] Make caching directory configurable
- [ ] Add SSL support

### How it works

First Pulp checks the local cache. It uses the filesystem for this, and stores all cached files in `/tmp/{uuid}`. 
A UUID is used to prevent collisions with other applications.
This directory will be created on bootstrap.
If the file exists, Pulp serves the file straight from the file system.

If the file does not exist locally, we query your Google Cloud Storage bucket to see if the file is present there. 
If that's the case we download the file to the local cache directory and serve it. 

### Invalidating Cache
Invalidating Cache is Easy™. 

```bash
curl -H "Authorization: Bearer Foo" -X "DELETE" http://example.com/file.png
```

Use the `PULP_TOKEN` environment variable to set the bearer token.
Adding authentication and not allowing it to be disabled was a conscious decision.
You don't want someone spamming your bucket driving up `get` operations.

## Configuration

| Env Var | Required | Description | Example | Default |
| - | - | - | - | - |
| PULP_BUCKET | Yes | GCS Bucket name | `foo-bar`, `foo.bar.com` | |
| PULP_PREFIX | Yes | Prefix for looking up objects on GCS | `foo`, `foo/bar` | |
| PULP_TOKEN   | Yes | Authentication Token for purging cache | `longandsecret` | | 
| PULP_INDEX  | No | File for index | `index.html`, `index.htm`, `home.html` | `index.html` |
| PULP_ADDRESS | No | Address on which Pulp server listens | `:8000`, `0.0.0.0:1337`, `127.0.0.1:80` | `0.0.0.0:8000` |

## Google Cloud Authentication

Pulp uses the `cloud.google.com/go` module. Please set the `GOOGLE_APPLICATION_CREDENTIALS` environment variable to a
location where you've stored your Service Account JSON key.