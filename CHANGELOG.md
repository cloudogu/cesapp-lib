# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
-[#58] make CheckDependencyVersion public
 - function can be used for dogu- and package-dependencies

## [v0.18.0] - 2025-01-15
### Changed
- [#55] Extract calculation of effective capabilities to be reusable

## [v0.17.0] - 2025-01-07
### Changed
- [#52] Adapt to changes in the CMS
  - Remove `Logo` and `BackgroundImage` from `MarketingDogu`
  - Add `Category` to `MarketingDogu`
### Fixed
- [#52] Add missing fields in `MarektingDogu` example

## [v0.16.0] - 2025-01-06
### Added
- [#53] Add Security field to `dogu_v2`, which can be to define security policies for the dogu.
  - This can be used for example in the pod security context on the kubernetes platform. 

## [v0.15.0] - 2024-11-13
### Added
- [#43] add a timestamp to `dogu_v1` and `dogu_v2` which represents the date and time when a dogu was created
- [#50] Add string representation to `core.Version`

## [v0.14.4] - 2024-11-06
### Fixed
- [#48] Map nginx dependency to `nginx-ingress` and `nginx-static`.
Both Dogus are required in the installation process in the dogu-operator.
Only mapping `nginx-ingress` can result in following installation order `ingress`, `cas`, `static` because the `nginx-ingress` has no dependency and the new dependency sorting algorithm is not deterministic.

## [v0.14.3] - 2024-10-30
### Fixed
- [#47] map nginx dependencies to k8s equivalent dogus 

## [v0.14.2] - 2024-10-18
### Fixed
- [#44] fix, that dogus with irrelevant optional dependencies were not included when sorting by dependency

## [v0.14.1] - 2024-10-17
### Changed
- use topological sorting to sort dogus by dependency

### Deprecated
- `SortDogusByDependency`: Use `SortDogusByDependencyWithError` instead for better error handling
- `SortDogusByInvertedDependency`: Use `SortDogusByInvertedDependencyWithError` instead for better error handling

## [v0.14.0] - 2024-09-18
### Changed
- Relicense to AGPL-3.0-only (#40)

## [v0.13.0] - 2024-09-16
### Added
- [#38] add a struct dedicated for dogu marketing data
   - Marketing information for dogus like description, deprecation state, URL to different translated release note, etc are bound to change independendly from an actual dogu release. 
   - This data can now reside in a different, independent structure.

### Fixed
- fix a typo regarding which key will be used for encryption

## [v0.12.2] - 2023-10-24
### Changed
- [#34] Reduce the wait time on failures while watching the etcd from 5 minutes to 10 seconds.

## [v0.12.1] - 2023-08-23
### Changed
- [#32] Change the log to error, if a dogu `GET` returns a `401` and the cache handles the request error.
  - In most cases user didn't recognize the permission error and just saw the following caching error.

## [v0.12.0] - 2023-03-24
### Added
- [#24] Add package `ssl` with functionality to generate selfsigned certificates.

### Removed
- Remove dogu-build-lib

### Changed
- Update ces-build-lib to 1.62.0

## [v0.11.0] - 2023-03-15
### Changed
- [#28] supplement core.Dogu documentation comments
  - these doc comments will result in the public [dogu documentation](https://github.com/cloudogu/dogu-development-docs)

## [v0.10.0] - 2023-03-03
### Added
- [#25] Extract the dogu configuration facility from cesapp to share them. 

## [v0.9.0] - 2022-11-17
### Added
- [#22] Support for extended volume definitions to add client-specific configurations.
- [#22] Support for extended service account definitions using a new `Kind` field.  
  This enables the use of service accounts for non-dogu services.

## [v0.8.0] - 2022-11-07
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