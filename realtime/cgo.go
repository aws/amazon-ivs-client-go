//go:build cgo

package realtime

/*
#cgo pkg-config: aws-ivs-realtime-1.49.0
#cgo LDFLAGS: -Wl,-rpath,$ORIGIN -Wl,-z,now
*/
import "C"
