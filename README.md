# pulls

Command-line access to open pull requests across a configurable set of Github
repos

## Install

```sh
$ go get github.com/adorableio/pulls
```

## Configure
You'll need to **create two config files** in order to use `pulls` in your
project(s).

### Github Configuration
In order to query pull requests on a list of repos you'll need to create an
access token in Github so it may be included on API requests.

Navigate to the [Personal Access Tokens](https://github.com/settings/tokens)
settings page and create a new token for `pulls` to use.

Now create a file named `.pulls.github.yml` in your `$HOME` directory and put
the new access token in it.

```yml
---
accessToken: ea1cb4ab3b5b2cdecb6e2f07e09a75116d24022a
```

### Repo Configuration
You may create a file named `repos.yml` at any directory level that makes sense
for your project structure. The utility will recursively look upwards in the
directory structure in order to find it.

*Suggestion*: Place `repos.yml` in your `$HOME` directory if there are a common
set of repos you always work with.

The contents of `repos.yml` are a list of github repos identified without the
leading `github.com` subpath.

```yml
---
- adorableio/pulls
- adorableio/avatars-api
- adorableio/hypeharvest
```

## Usage
Now you can simply run `pulls` in any subdirectory of where you placed your
`repos.yml` file and you will be presented with a list of all open PRs from
all of the repos listed.

```sh
$ pulls
```
