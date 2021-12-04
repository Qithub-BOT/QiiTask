#!/bin/sh
# =============================================================================
#  `go install` test (Go 1.16+)
# =============================================================================

cd /tmp || {
    echo >&2 "failed to move temp dir"
    exit 1
}

go install "github.com/Qithub-BOT/QiiTask/qiitask@latest" || {
    echo >&2 "failed to install qiitask command via 'go install'"
    exit 1
}

qiitask --version | grep "qiitask version" || {
    echo >&2 "failed to execute command. The output does not contain a valid version info."
    exit 1
}
