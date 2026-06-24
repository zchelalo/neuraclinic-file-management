#!/bin/sh
set -eu

should_run_migrations() {
	[ "${AUTO_RUN_MIGRATIONS:-false}" = "true" ] || return 1
	[ "$#" -gt 0 ] || return 1

	case "$1" in
		neuraclinic-file-management|/usr/local/bin/neuraclinic-file-management)
			return 0
			;;
	esac

	return 1
}

if should_run_migrations "$@"; then
	/usr/local/bin/run-migrations.sh
fi

exec "$@"

