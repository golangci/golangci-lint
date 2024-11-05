# robustio

Extracted from go1.19.1/src/cmd/go/internal/robustio

There is only one modification:
- ERROR_SHARING_VIOLATION extracted from go1.19.1/src/internal/syscall/windows/syscall_windows.go to remove the dependencies to `internal/syscall/windows`

## History

- https://github.com/golangci/golangci-lint/pull/5100
    - Move package from `internal/robustio` to `internal/go/robustio`
