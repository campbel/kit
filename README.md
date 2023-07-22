# kit

`kit` is a pre-processor for [task](https://taskfile.dev/). 
Kit is used as a top-level replacement to task, but depends on task being installed and in the PATH.

Kits main purpose is to support including remote taskfiles.

```yaml
version: "3"

includes:
  hello-world: github.com/campbel/kit/tasks/hello-world
```

It does this by fetching the remote file with [hashicorp/go-getter](https://github.com/hashicorp/go-getter) storing it in a temporary local directory `.kit` and then compiling a modified taskfile to run.

## Why?

Utilizing `task` on a team or across multiple projects can be frustraing because you're often writing many similar tasks. Typically for any given set of technologies teams will utilize the same set of tasks. Kit enables you to easily share these tasks between projects without relying on copy paste or git submodules.

## Example

For a taskfile with a remote includes.

```yaml
# ./Taskfile.yaml
version: "3"

includes:
  hello-world: github.com/campbel/kit/tasks/hello-world
```

```sh
$ kit --list
task: Available tasks for this project:
* hello-world:run:       Run the project
```

```sh
$ kit hello-world:run
task: [hello-world:os:run] echo "Hello, darwin!"
Hello, darwin!
```