# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Added function to get the whole registry content as RegistryNode (#6)

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