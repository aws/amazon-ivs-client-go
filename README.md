# Amazon IVS Client SDKs for Go

Go bindings for the Amazon Interactive Video Service (IVS) client SDKs.

Each product is a separate Go module in its own directory, with its own version. Add only the modules you use.

## Modules

| Directory | Import path | What it does | Platforms |
|---|---|---|---|
| [`realtime`](realtime/README.md) | `github.com/aws/amazon-ivs-client-go/realtime` | Joins an IVS Real-Time Stage as a headless server-side participant: subscribes to participants' audio and publishes PCM audio and subtitle cues. | Linux amd64, Linux arm64 |

## How the modules work

Each module is a cgo binding over a native IVS SDK. The module holds Go source only. You download the native SDK for your platform from the IVS CDN, and the build finds it through `pkg-config`. Each module's README lists the requirements and the download, build and deploy steps.

## Versions

Each module is tagged `<directory>/vX.Y.Z`. Add a module with `go get`:

```sh
go get github.com/aws/amazon-ivs-client-go/realtime
```

A module version builds only against the native SDK of the same version.

## Documentation

- Module docs: `go doc github.com/aws/amazon-ivs-client-go/realtime`
- Amazon IVS documentation: https://docs.aws.amazon.com/ivs/

## Security

See [CONTRIBUTING](CONTRIBUTING.md#security-issue-notifications) for more information.

## License

This project is licensed under the Apache-2.0 License.
