# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Add general packages from the cesapp in order to use these components with other applications.
  - Packages: `config`, `credentials`, `dependencies`, `doguConf`, `keys`, `tasks` (#20)

### Changed
- Update ces-build-lib to 1.57.0
- Update dogu-build-lib to 1.10.0

## Fixed
- Fixed a bug where an error in optional dogu dependency check would overwrite mandatory dependency errors

## [v0.7.0] - 2022-11-03
### Changed
- Make the retry policy for registry and dogu backend calls configurable (#18)

## [v0.6.0] - 2022-09-16
### Changed
- Update go version from 1.18.1 to 1.18.6
- Update Makefiles to version 7.0.1
- Update archive package to handle archives in memory (#16)

### Fixed
- Fix date of files added to archive (#14)

## [v0.5.0] - 2022-09-12
### Added
- Function to get the whole registry content as RegistryNode (#6)

### Changed
- Set default make goal to 'compile'
- switch etcd client lib to `go.etcd.io/etcd/client/v2` (#6)

## [v0.4.0] - 2022-08-29
### Added
- [#7] Added general logging interface. See [Logger-Interface](core/logger.go) for more information.

### Changed
- [#9] Moved dogu printing facilities back to the originating `cesapp`
  - Printing message to the stdout stream does not belong into this library but in the calling client
  - This enables to reduce the size of the logging interface

## [v0.3.0] - 2022-08-16
### Added
- Added functions to pack files and logs into archives

## [v0.2.0] - 2022-06-08
### Added
- Watch context used to watch registry keys that notifies the client when they are are changed. 

## [v0.1.0] - 2022-05-10
### Added
- Initial release
   - Extracted `core` into the library providing access to the dogu struct and common functions.
   - Extracted `registry` into the library providing access to a dogu registry.