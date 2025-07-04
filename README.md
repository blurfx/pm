# pm

pm is a unified package manager for javascript.
it provides a seamless cli for users, regardless of which package manager your project uses.

it's still a work in progress, but you can try it out.

## Installation

#### Install from Homebrew:

```sh
brew install blurfx/tap/pm
```

#### Build from source:

```sh
go build
```

## Usage

interactive mode:

```sh
pm
```

running package manager commands:

```sh
pm i <packages>
pm add <packages>
pm rm <packages>
pm uninstall <packages>
pm ci
# ...
```

running package scripts:

```sh
pm dev
pm build
pm start
# ...
```
