# IVS Linux Broadcast SDK for Go

`github.com/aws/amazon-ivs-client-go/realtime` is the Go binding of the IVS Linux Broadcast SDK.

Its package, `realtime`, joins an Amazon IVS Real-Time Stage as a headless server-side participant. It subscribes to remote participants' audio and delivers PCM to your program. It can publish PCM audio that your program pushes, and subtitle cues.

`realtime` is the media SDK that joins Stages. Creating Stages and participant tokens uses the AWS SDK for Go's IVS Real-Time service package, `github.com/aws/aws-sdk-go-v2/service/ivsrealtime`, through its `CreateParticipantToken` API.

## Requirements

- Linux on `amd64` or `arm64`.
- Go 1.24 or newer.
- cgo enabled, with a C compiler. The package does not build with `CGO_ENABLED=0`.
- `pkg-config` (or pkgconf), which locates the native SDK at build time.
- libcurl 7.68 or newer (`libcurl.so.4`) and glibc 2.34 or newer, at build time and at run time.
- A CA certificate bundle at run time.

## Getting started

Add the module to your project:

```sh
go get github.com/aws/amazon-ivs-client-go/realtime
```

Import the package. Its name is `realtime`:

```go
import "github.com/aws/amazon-ivs-client-go/realtime"
```

The module holds Go source only. To build, you also need the native SDK: the shared library, its C headers, and a pkg-config file. Download the SDK for your architecture (`go env GOARCH` prints `amd64` or `arm64`) from `https://broadcast.live-video.net/<version>/aws-ivs-realtime-<version>-linux-<arch>.tar.gz`. The same directory holds `SHA256SUMS`. Extract the SDK into a directory named after its version, such as `ivs-sdk/<version>` next to your project, and add that directory to `.gitignore`. These commands read the SDK version from the module version in your `go.mod`:

```sh
VERSION=$(go list -f '{{index .CgoPkgConfig 0}}' github.com/aws/amazon-ivs-client-go/realtime)
VERSION=${VERSION#aws-ivs-realtime-}
ARCH=$(go env GOARCH)
curl -fsSLO "https://broadcast.live-video.net/$VERSION/aws-ivs-realtime-$VERSION-linux-$ARCH.tar.gz"
curl -fsSL "https://broadcast.live-video.net/$VERSION/SHA256SUMS" | sha256sum --check --ignore-missing
mkdir -p "ivs-sdk/$VERSION"
tar -xzf "aws-ivs-realtime-$VERSION-linux-$ARCH.tar.gz" -C "ivs-sdk/$VERSION"
export PKG_CONFIG_PATH="$PWD/ivs-sdk/$VERSION/lib/pkgconfig"
go build
```

The SDK version must equal the module version in your `go.mod`. The package asks pkg-config for `aws-ivs-realtime-<version>`, the name the `go list` command above prints. If the build fails with `Package aws-ivs-realtime-<version> was not found in the pkg-config search path`, either `PKG_CONFIG_PATH` does not point at the SDK's `lib/pkgconfig`, or the SDK's version differs from the module version.

`go vet`, `go test`, and gopls compile the package too, so they need the same `PKG_CONFIG_PATH`. For gopls, set it in the `build.env` setting.

## Example

This program joins a Stage, subscribes to every remote participant's audio, and logs the PCM audio it receives. It publishes nothing. On Ctrl-C, it leaves the Stage, waits for the disconnect, and closes the Stage. It reads the participant token from the `IVS_STAGE_TOKEN` environment variable.

```go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/aws/amazon-ivs-client-go/realtime"
)

// audioSubscriber subscribes to every remote participant's audio and publishes
// nothing.
type audioSubscriber struct{}

func (audioSubscriber) SubscriptionTypes(realtime.ParticipantInfo) realtime.StreamType {
	return realtime.StreamTypeAudio
}

func (audioSubscriber) ShouldPublish(realtime.ParticipantInfo) bool {
	return false
}

func (audioSubscriber) StreamsToPublish(realtime.ParticipantInfo) realtime.PublishStreams {
	return realtime.PublishStreams{}
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	disconnected := make(chan struct{}, 1)
	hooks := realtime.StageHooks{
		OnConnectionStateChanged: func(event realtime.ConnectionStateChangedEvent) {
			if event.Err != nil {
				log.Printf("connection error: %v", event.Err)
			}
			log.Printf("connection: %s", event.State)
			if event.State == realtime.ConnectionStateDisconnected {
				select {
				case disconnected <- struct{}{}:
				default:
				}
			}
		},
	}

	stage, err := realtime.NewStage(realtime.StageConfig{
		Token:    os.Getenv("IVS_STAGE_TOKEN"),
		Strategy: audioSubscriber{},
		Hooks:    hooks,
	})
	if err != nil {
		return err
	}
	defer stage.Close()

	err = stage.SetPCMCallback(func(buf realtime.PCMBuffer) {
		log.Printf("%d frames at %d Hz", buf.NumFrames, buf.SamplesPerSec)
	})
	if err != nil {
		return err
	}

	if err := stage.Join(); err != nil {
		return err
	}

	<-ctx.Done()
	if err := stage.Leave(); err != nil {
		return err
	}
	<-disconnected
	return nil
}
```

The PCM callback runs on an audio thread, so return from it promptly: do not call SDK methods from it, and hand the buffer to a goroutine for anything more than trivial work. The callback owns the buffer it receives and may keep it. `run` returns its errors to `main`, because `log.Fatal` exits without running deferred calls such as `stage.Close`.

## Native library

An executable built with this package needs `libaws-ivs-realtime.so.1` at run time. It sets two rpath entries: the executable's own directory first, then the SDK's `lib/` directory. The second entry makes `go run` and `go test` work on the build machine without `LD_LIBRARY_PATH`; a deployed executable uses the copy beside it.

To deploy an executable to another machine or container image, copy `<dir>/lib/libaws-ivs-realtime.so.1` into the same directory as the executable. Use `cp` without `-a`: in the SDK, `libaws-ivs-realtime.so.1` is a symlink, and `cp -a` copies the symlink, which then points at nothing.

```sh
cp "ivs-sdk/$VERSION/lib/libaws-ivs-realtime.so.1" /opt/myserver/
```

```
/opt/myserver/
  myserver
  libaws-ivs-realtime.so.1
```

Alternatively, install the library on the loader's search path, or set `LD_LIBRARY_PATH` to a directory that holds it. Deploy the library from the SDK the executable was built with.

### Docker

This multi-stage Dockerfile builds the program on Amazon Linux 2023, downloads the module's SDK version for the build stage's architecture (`go env GOARCH`), and copies the executable and the library into an Amazon Linux 2023 runtime image, which ships libcurl and a CA bundle. The program runs as a non-root user.

```dockerfile
FROM public.ecr.aws/amazonlinux/amazonlinux:2023 AS build
RUN dnf install -y golang gcc pkgconf-pkg-config tar gzip && dnf clean all
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
RUN SDK=$(go list -f '{{index .CgoPkgConfig 0}}' github.com/aws/amazon-ivs-client-go/realtime) \
 && VERSION=${SDK#aws-ivs-realtime-} \
 && ARCH=$(go env GOARCH) \
 && curl -fsSLo /tmp/ivs-sdk.tar.gz "https://broadcast.live-video.net/${VERSION}/aws-ivs-realtime-${VERSION}-linux-${ARCH}.tar.gz" \
 && mkdir -p /opt/ivs-sdk \
 && tar -xzf /tmp/ivs-sdk.tar.gz -C /opt/ivs-sdk \
 && rm /tmp/ivs-sdk.tar.gz
ENV PKG_CONFIG_PATH=/opt/ivs-sdk/lib/pkgconfig
COPY . .
RUN go build -o /out/myserver .

FROM public.ecr.aws/amazonlinux/amazonlinux:2023
COPY --from=build /out/myserver /opt/myserver/myserver
COPY --from=build /opt/ivs-sdk/lib/libaws-ivs-realtime.so.1 /opt/myserver/
USER 10001
ENTRYPOINT ["/opt/myserver/myserver"]
```

`COPY --from` copies the file the symlink points to, so `/opt/myserver/libaws-ivs-realtime.so.1` is a regular file.

## Further reading

The package documentation covers threading, errors, and lifetimes:

```sh
go doc github.com/aws/amazon-ivs-client-go/realtime
```

Runnable examples are in `example_test.go`.
