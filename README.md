# local gce metadata

This tool runs a simple mirror of the GCE instance metadata server to convince
the `gcloud` command line tool it is running in GCE. You can use this tool to
pass an arbitrary bearer token to `gcloud` without too much trouble.

## Setup

There are a few small configuration changes you'll need to make before running the
metadata server. If you'd like, you can simply `source env.sh` before running
`gcloud`, but the individual changes are explained below.

1. `gcloud` caches its GCE detection in `~/.config/gcloud/gce`. You can override
   this by simply writing `True` to that file.

```bash
$ echo -n True > ~/.config/gcloud/gce
```

2. `gcloud` uses the `metadata.google.internal` hostname to access the metadata
   service. You can either set up DNS for this host, or simply set the
   `GCE_METADATA_ROOT` environment variable.

```bash
$ export GCE_METADATA_ROOT=127.0.0.1
```

3. If you want to avoid polluting your existing config, you can set the `gcloud`
   config directory to a new temporary directory. Alternatively, you can use the
   `--account local` flag in `gcloud` to specify this tool's service account.

```bash
$ export CLOUDSDK_CONFIG=`mktemp -d`
```

## Usage

```bash
$ docker build -t local-gce-metadata .
```

```bash
$ docker run --rm -p 80:80 local-gce-metadata -h
Usage of /local-gce-metadata:
  -account string
    	name of service account to advertise (default "local")
  -token string
    	service account bearer token
```

As shown above, in order to use this tool, you'll need a bearer token to serve. Assuming you have such a token, you can run the server as follows:

```bash
$ docker run --rm -p 80:80 local-gce-metadata -token ya29.xyz...
2020/04/02 20:37:57 172.17.0.1:48950 GET /computeMetadata/v1/instance/service-accounts/ - Python-urllib/2.7
2020/04/02 20:37:57 172.17.0.1:48954 GET /computeMetadata/v1/instance/service-accounts/local/?recursive=True - gcloud/272.0.0...
2020/04/02 20:37:57 172.17.0.1:48954 GET /computeMetadata/v1/instance/service-accounts/local/token - gcloud/272.0.0...
```

```bash
$ gcloud --account local projects list
PROJECT_ID      NAME            PROJECT_NUMBER
wizardly_knuth  wizardly_knuth  1234567890
```

