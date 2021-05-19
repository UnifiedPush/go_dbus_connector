# Go UnifiedPush Connector for DBus systems

This library can be embeded into any progam that supports importing a C library (in addition to the native Go). Check out this link (TODO) for existing wrapper libraries in other languages.  

Currently, the api is unstable, though major changes are not expected.

## Docs

### C

The header files for the C library contains basic documentation, since it is only a wrapper around the Go lib that documentation will be very relevent.

It can be built from source using just the Go command (check <./Makefile>). The releases contain statically linkable .a files and dynamically loadable .so files.

### Go

The docs can be found here <https://pkg.go.dev/unifiedpush.org/go/dbus_connector>.
[![Go Reference](https://pkg.go.dev/badge/unifiedpush.org/go/dbus_connector.svg)](https://pkg.go.dev/unifiedpush.org/go/dbus_connector)
