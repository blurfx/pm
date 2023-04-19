# pm

pm is a unified package manager for javascript.
it provides a seamless cli for users, regardless of which package manager your project uses.

it's still a work in progress, but you can try it out.

## Installation

Build from source:

```sh
go build
```

move the binary to somewhere in your `PATH`.

## Usage

install dependencies:
```sh
pm install <packages> 
# or 
pm i <packages>
# or 
pm add <packages>
```

install all dependencies:
```sh
pm install
# or
pm i
```

clean install:
```sh
pm ci
# or 
pm install --frozen-lockfile
```

... and more