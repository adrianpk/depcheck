# depcheck

depcheck is a simple command-line interface (CLI) tool designed to provide a quick analysis of the dependencies in a Go project. 

## Why?

When working on a Go project, it's often necessary to understand the dependencies that your project relies on. This can be especially important when reviewing pull requests that introduce new dependencies. depcheck provides a quick and easy way to get insights about all the dependencies of your project.

## What?

When executed in the top folder of the module, depcheck scans the `go.mod` file and fetches information about each dependency from the GitHub API. It then prints statistics about each dependency, including:

- Full name of the repository
- Whether the repository is a fork
- Parent repo if it is a fork
- Number of stargazers
- Number of watchers
- Number of open issues
- License information
- Default branch

## Usage

Follow these steps to install and use `depcheck`:

Clone the repository:
```sh
$ git clone https://github.com/adrianpk/depcheck.git
```

Navigate to the `depcheck` directory:
```sh
$ cd depcheck
```

Install the application:
```sh
$ make install
```

Once installed, you can run `depcheck` on the root folder of any Go module to analyze its dependencies.
```sh
$ depcheck
```

Or you can specify a sort flag:
```sh
$ depcheck -sort=stars
```

In that case, the usage instructions would be slightly different. The user would run the compiled executable directly, rather than using `go run`. Here's how you can update the instructions:

```markdown
## Usage

Follow these steps to install and use `depcheck`:

Clone the repository:
```sh
$ git clone https://github.com/adrianpk/depcheck.git
```

Navigate to the `depcheck` directory:
```sh
$ cd depcheck
```

Install the application:
```sh
$ make install
```

Once installed, you can run `depcheck` on the root folder of any Go module to analyze its dependencies:

```sh
$ depcheck
```

Or you can specify a sort flag:

```sh
$ depcheck -sort=stars
```

The available sort flags are:

* `watchers`: Sorts the repositories by the number of watchers. This is the default sort flag.
* `stars`: Sorts the repositories by the number of stars.
* `forks`: Sorts the repositories by the number of forks.
* `issues`: Sorts the repositories by the number of open issues.

If no sort flag is passed, the application will sort the repositories by the number of watchers.

**Prerequisite:** Before running the code, make sure to set the `GITHUB_TOKEN` environment variable to your GitHub personal access token. This is necessary for the code to authenticate with the GitHub API. You can generate a personal access token in your GitHub settings.

```bash
export GITHUB_TOKEN=your_github_token
```

Replace `your_github_token` with your actual GitHub personal access token.

## Sample output
```markdown
GitHub token found, proceeding...
Scanning go.mod file...
                       Name|  IsFork|  Parent Repo|  Stargazers|  Watchers|  Open Issues|                                  License|Default Branch
     sagikazarmark/locafero|   false|             |           4|         4|            1|                              MIT License|main
    sagikazarmark/slog-shim|   false|             |           7|         7|            1|  BSD 3-Clause "New" or "Revised" License|main
  inconshreveable/mousetrap|   false|             |         232|       232|            0|                       Apache License 2.0|master
            subosito/gotenv|   false|             |         287|       287|            0|                              MIT License|master
      magiconair/properties|   false|             |         321|       321|           23|        BSD 2-Clause "Simplified" License|main
         pmezard/go-difflib|   false|             |         379|       379|            4|                                    Other|master
          pelletier/go-toml|   false|             |        1661|      1661|           22|                                    Other|v2
                spf13/pflag|    true|  ogier/pflag|        2328|      2328|          178|  BSD 3-Clause "New" or "Revised" License|master
                 spf13/cast|   false|             |        3376|      3376|           89|                              MIT License|master
              hashicorp/hcl|   false|             |        5127|      5127|          208|               Mozilla Public License 2.0|main
                spf13/afero|   false|             |        5752|      5752|          133|                       Apache License 2.0|master
            davecgh/go-spew|   false|             |        5946|      5946|           63|                              ISC License|master
     mitchellh/mapstructure|   false|             |        7763|      7763|           82|                              MIT License|main
           sourcegraph/conc|   false|             |        8562|      8562|           13|                              MIT License|main
          fsnotify/fsnotify|   false|             |        9250|      9250|           26|  BSD 3-Clause "New" or "Revised" License|main
           stretchr/testify|   false|             |       22304|     22304|          387|                              MIT License|master
                spf13/viper|   false|             |       26085|     26085|          499|                              MIT License|master
                spf13/cobra|   false|             |       36448|     36448|          274|                       Apache License 2.0|main
```
