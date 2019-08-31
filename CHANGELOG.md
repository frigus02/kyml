# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

## [20190831]

### Added

- Improved zsh completions.

## [20190103]

### Added

- Added hint about `--update` option to error messages of `kyml test` command, which is shown when the snapshot file does not exist or the snapshot doesn't match.

## [20181227]

### Added

- Cache resolved images in `kyml resolve`. If your manifests include the same image reference multiple times, kyml will only ask the registry once.
- The command `kyml tmpl` now errors if the manifest contains a template key, which was not specified in the command flags.

### Fixed

- Resource deduplication didn't correctly check for resource name. It does now.

## 20181226

- First release.

[unreleased]: https://github.com/frigus02/kyml/compare/v20190831...HEAD
[20190831]: https://github.com/frigus02/kyml/compare/v20190103...v20190831
[20190103]: https://github.com/frigus02/kyml/compare/v20181227...v20190103
[20181227]: https://github.com/frigus02/kyml/compare/v20181226...v20181227
