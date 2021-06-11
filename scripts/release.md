# Release instructions

To release a new version of kyml, follow these steps:

1. Choose a version number in the following format:

   ```sh
   VERSION=$(date -u +%Y%m%d)
   ```

   If you need to release a second version on the same day, add a dot followed by an incrementing number to the end, e.g.: 20181228.2, 20181228.3, ...

1. Create a branch.

   ```sh
   git checkout -b version-$VERSION
   ```

1. Update the [CHANGELOG.md](../CHANGELOG.md).

1. Update the version in the installation instructions in the [README.md](../README.md).

1. Commit and push the branch. Create a pull request for the branch. Then wait for the builds to succeed.

   ```sh
   git commit -m "Update changelog and readme for v$VERSION"
   ```

1. Tag the commit with `v$VERSION`.

   ```sh
   git tag v$VERSION
   ```

1. Push the tag. This will trigger another build, which creates a release on GitHub and uploads artifacts.

   ```sh
   git push --tags
   ```

1. Merge the pull request.

1. Update version and checksums in Homebrew formula [frigus02/tap/kyml](https://github.com/frigus02/homebrew-tap/blob/main/kyml.rb).
