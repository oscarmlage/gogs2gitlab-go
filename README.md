# gogs2gitlab

Script that migrates GoGS repos to Gitlab. Needs tokens in both sides:

- GoGS token
  - GoGS user
  - GoGS pass
- Gitlab token


## Configure

Configuration runs in `{home}/.config/gogs2gitlab/gogs2gitlab.ini`, fill the
gaps:

```ini
gogs_proto = https://
gogs_host = gogs.host.com
gogs_token = whatever-gogs-token
gogs_user = user
gogs_pass = pass

gitlab_proto = https://
gitlab_host = gitlab.host.com
gitlab_token = whatever-gitlab-token
```

## Run

```sh
$ ./gogs2gitlab
```
