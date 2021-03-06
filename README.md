# hass-run

Allow executing long running commands from Home-Assistant and update a defined entity's attributes with:

```
state: running/success/failure
--
output: 'text'
running: false
exit_code: 0
started_at: '2022-04-29T08:24:51.757144+02:00'
updated_at: '0001-01-01T00:00:00Z'
ended_at: '2022-04-29T21:04:57.142761+02:00'
duration: 10000
```

- `output` aggregates `stdout` and `stderr`
- `exit_code` is `0` if the command is still running
- Dates are set to `0001-01-01T00:00:00Z` if not relevant (`ended_at` when the command is still running for instance)

## Installation

### Using Docker

- Replace `__RELEASE_URL__` by a valid `tar.gz` release archive
- Add it to your Dockerfile 

```
RUN wget -O /tmp/hass-run-release.tar.gz __RELEASE_URL__.tar.gz && \
  mkdir /tmp/hass-run-release && \
  tar -xzf /tmp/hass-run-release.tar.gz -C /tmp/hass-run-release && \
  cp /tmp/hass-run-release/hass-run /usr/bin && \
  chmod +x /usr/bin/hass-run && \
  rm /tmp/hass-run-release.tar.gz && \
  rm -rf /tmp/hass-run-release
```

### Other

- Download [a release `tar.gz`](https://github.com/simon-watiau/hass-run/releases)
- Copy the `hass-run` binary in your `$PATH`
- Use `chmod +x` to make it executable

## Usage

- `hass-run` starts your command as a daemon and exits immediatly.
- A PID file is kept to optionnaly kill the command later on

### Configuring host and token

#### Using a configuration file
Add a `hass-run.yaml` in `.`, `$HOME` or `/etc` with your host and token as follow:

```
host: "https://my_host_assistant_url.com"
bearer: "XXXXXX"
```

##### Using command line flags

Add the flags `--host` for your Home-Assistant host and `--bearer` for your token.

### Examples

**Run a command with config file:**

`hass-run run shell.my_entity /tmp/my_command.pid -- my_command`

**Run a command without config file:**

`hass-run run --host https://my_hass_url.com --bearer XXXTOKENXXX shell.my_entity /tmp/my_command.pid -- my_command`

**Run multiple commands:**

`hass-run run --host https://my_hass_url.com --bearer XXXTOKENXXX shell.my_entity /tmp/my_command.pid -- bash -c "my_command_1 && my_command_2"`

**Set Home-Assistant configuration:**

```
shell_command:
  my_command: hass-run run shell.my_command ./my_command.pid -- bash -c "cd / && ls"
```

**Kill a running command:**

- `hass-run kill --host https://my_hass_url.com --bearer XXXTOKENXXX shell.my_entity /tmp/my_command.pid`
- `hass-run kill shell.my_entity /tmp/my_command.pid`


## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D
