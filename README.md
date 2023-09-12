# knex

Pluggable Certification.

## Building and Registering a Plugin

- Implement the [Plugin](./plugin.go#L43)
- Create an `init` function somewhere in your plugin codebase that calls the [Register](./plugin.go#L31) function
- Submit a PR to the repository adding a blank-initialization of your plugin code (ex. [here](./plugin/registration/add_plugins_here.go#L8))
- Ensure the go.mod value for your plugin points to your version. The repository
  encourages semantic versioning, and will represent plugin versions to users.
  You are encouraged to ensure version changes are associated with new behaviors.

## Writing Logs and Artifacts

- Knex will pass a logger and an artifact writer to your plugin via the
  `context`.
- For the logger, utilize the logr helper function
  [FromContextOrDiscard](https://pkg.go.dev/github.com/go-logr/logr#FromContextOrDiscard)
  (or equivalents).
- For ArtifactsWriter, utilize the helper function
  [WriterFromContext](https://pkg.go.dev/github.com/redhat-openshift-ecosystem/openshift-preflight/artifacts#WriterFromContext)
- Plugins are generally discouraged from reconfiguring the included logger or artifact writer.

## Binding Environment/Flags

- A plugin will be passed a
  [pflag.FlagSet](https://pkg.go.dev/github.com/spf13/pflag#FlagSet) to its
  `BindFlags` method. It should bind all flags necessary for the plugin to
  operate at this time.
- These flags are converted to environment variables using viper's
  [AutomaticEnv](https://pkg.go.dev/github.com/spf13/viper#AutomaticEnv). Users
  will need to prefix your environment variables with `PFLT_`.
    - Dashes are converted to hyphens.
    - Other non-env special characters are not supported (e.g. period)


## Review notes

- The old preflight code was moved to the 'openshift-preflight' directory. It was moved as-is.
- The knex repo was copied to the root of github.com/redhat-openshift-ecosystem/openshift-preflight
- The module, however, is now named github.com/redhat-openshift-ecosystem/preflight
- All of the existing plugins have been updated to point to this repo, but via a replace for now to bcrochet.
- The plugin package was moved to the root of the repo. Instead of plugin.Plugin, it is now preflight.Plugin.
- types was also moved to the root. It is now preflight.Check, preflight.Result, etc.
