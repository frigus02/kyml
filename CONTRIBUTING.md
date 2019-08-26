# Contributing

From opening a bug report to creating a pull request: every contribution is appreciated and welcome. If you're planning to implement a new feature or change a command please create an issue first. This way we can ensure that your precious work is not in vain.

## Process

1. Setup prerequisites on your machine. Currently the following are required:

   - go v1.12.x

1. Fork, then clone the repository.

1. Try to run the tests to see if everything works correctly. If you're doing this for the first time, go should automatically download the necessary dependencies.

   ```sh
   go test ./...
   ```

1. Create a new branch based on master and start to make your changes.

   ```sh
   git checkout -b my-feature
   ```

1. Once you're happy, push your branch and [open a pull request](https://github.com/frigus02/kyml/compare) ([help](https://help.github.com/articles/creating-a-pull-request/)).

   ```sh
   git push origin -u my-feature
   ```

## Additional notes

Some things that will increase the chance that your pull request is accepted:

- Write tests
- Follow the existing coding style
- Write a [good commit message](https://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html)
