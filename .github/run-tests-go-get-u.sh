#!/bin/sh
# =============================================================================
#  `go get -u` test (Go 1.15)
# =============================================================================

cd /tmp || {
    echo >&2 "failed to move temp dir"
    exit 1
}

GO111MODULE="on" go get -u "github.com/Qithub-BOT/QiiTask/qiitask@latest" || {
    echo >&2 "failed to install qiitask command via 'go get -u'"
    exit 1
}

qiitask --version | grep "qiitask version" || {
    echo >&2 "failed to execute command. The output does not contain a valid version info."
    exit 1
}
