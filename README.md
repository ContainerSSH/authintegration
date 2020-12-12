[![ContainerSSH - Launch Containers on Demand](https://containerssh.github.io/images/logo-for-embedding.svg)](https://containerssh.github.io/)

<!--suppress HtmlDeprecatedAttribute -->
<h1 align="center">ContainerSSH Authentication Library</h1>

[![Go Report Card](https://goreportcard.com/badge/github.com/containerssh/library-template?style=for-the-badge)](https://goreportcard.com/report/github.com/containerssh/library-template)
[![LGTM Alerts](https://img.shields.io/lgtm/alerts/github/ContainerSSH/library-template?style=for-the-badge)](https://lgtm.com/projects/g/ContainerSSH/library-template/)

This library provides integration between the [sshserver library](https://github.com/containerssh/sshserver) and the [auth library](https://github.com/containerssh/auth)

<p align="center"><strong>Note: This is a developer documentation.</strong><br />The user documentation for ContainerSSH is located at <a href="https://containerssh.github.io">containerssh.github.io</a>.</p>

## Using this library

This library can be used to provide an authenticating overlay for ContainerSSH. It stacks well with other libraries. To use it you must first call the `authintegration.New()` method. This method has three parameters:

- `authClient` is an authentication client from the [auth library](https://github.com/containerssh/auth).
- `backend` is another implementation of the `Handler` interface from the [sshserver library](https://github.com/containerssh/sshserver).
- `behavior` influences when the backend is called for authentication purposes.
  - `BehaviorNoPassthrough` means that the backend will not be used for authentication, only for getting further handlers.
  - `BehaviorPassthroughOnFailure` will give the backend an additional chance to authenticate the user if the authentication server returns a failure.
  - `BehaviorPassthroughOnSuccess` passes the credentials to the backend for additional checks of an already verified successful authentication.
  - `BehaviorPassthroughOnUnavailable` passes the authentication to the backend as a fallback if the authentication server failed to return a valid response. 

For example:

```go
handler := authintegration.New(
    auth.ClientConfig{
        URL: "http://localhost:8080"
        Password: true,
        PubKey: false,
    },
    otherHandler,
    logger,
    authintegration.BehaviorNoPassthrough,
)
```

You can then use the handler to launch an SSH server:

```go
server, err := sshserver.New(
    cfg,
    handler,
    logger,
)
```
