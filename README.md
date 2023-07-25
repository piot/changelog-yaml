# Changelog generator from YAML


## Usage

```shell
changelog-yaml -t markdown < changelog.yaml > CHANGELOG.md
```

## Changelog Yaml format

### Supported change types

#### New Functionality

* **added**: new functionality was added. (`feat`)

#### Change

* **breaking**: the change needs the user of the library to modify their code.
* **fixed**: a bug was fixed.
* **workaround**: a bug was alleviated or temporarily bypassed with a workaround / "hack". The solution usually has low quality, is a short term remedy and must be fixed properly in upcoming versions.
* **changed**: a behaviour or code was changed.
* **improved**: code was changed to be of better quality and stability. (`enhancements`)
* **refactored**: the internal functions or code was improved or moved around, but might not be of any perceived value for the user of the library. (neither fixes a bug nor adds a feature)
* **performance**: code was changed in order to improve performance. (`perf`)

#### Removed or about to be removed

* **removed**: code was removed. it was usually marked as deprecated in previous releases.
* **deprecated**: code has been marked as deprecated. the code will be `removed` in a future release.

#### Other

* **docs**: updated documentation
* **tests**: changed or added tests
* **experimental**: code has been added, but not sure if it will work as intended, and it might not be supported in the future.
* **noted**. (known issues)

### Autolinks

#### Pull Request link

`#[number]` will be replaced with a link to that pull request for that repository, e.g. `#1`

#### Commit hash link

`$[hash]` gets replaced with a link to that specific github hash

#### Profile link

`@[GithubUsername]` will be replaced with a link to the user, e.g. `@piot` -> https://github.com/piot/

#### Admonition

`[ADMONITION]:[space] text`. Admonition types supported:

* NOTE
* IMPORTANT
* WARNING

Example:

```text
NOTE: This release requires latest firmware update
```

### Example

```yaml
repo: piot/nimble
repos:
  clog:
    name: CLog
    repo: piot/clog
    description: Basic logging
  secure-random-c:
    name: Secure Random
    repo: piot/secure-random-c
    description: Multi-platform Secure Random
  nimble-client-c:
    name: Nimble Client
    repo: piot/nimble-client-c
    description: Nimble Protocol Client
  nimble-server-lib:
    name: Nimble Server Library
    repo: piot/nimble-server-lib
    description: Nimble Server Library
releases:
  - name: 'v0.0.1-a06'
    date: '2023-06-22'
    notice: Minor compile fixes for emscripten.
    repos:
      clog:
        fixed:
          - use `tc_snprintf` instead of `sprintf` (#1)

      secure-random-c:
        workaround:
          - '`secureRandomUInt64()` on emscripten that only return 0'

  - name: 'v0.0.1-a05'
    date: '2023-06-14'
    notice: Hot fixes to alleviate skip ahead problems.
    repos:
      nimble-engine-client:
        changed:
          - Increase wait time between each skip ahead attempt

      nimble-server-lib:
        changed:
          - Skip increase forcedStepInRowCounter if client transport connection
            is downloading game state
```
