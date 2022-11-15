# Remote: A very bad run script manager

- [Remote: A very bad run script manager](#remote-a-very-bad-run-script-manager)
- [What?](#what)
- [Why?](#why)
- [How?](#how)
  - [Installation](#installation)
    - [Dependencies](#dependencies)
    - [Install](#install)
  - [Usage](#usage)
- [License](#license)

# What?

Remote is a run script manager. I hope you have enough `sh`/`cmd` experience.

# Why?

I didn't want to type `python -s src.build_main` and then `python -s src.main` every time I want to run my project. I just wanted to do a simple `remote run` and be done with it.

# How?

## Installation

### Dependencies
- [Go](https://golang.org/dl/)

### Install
Step 1. Clone the repo
```sh
git clone https://github.com/kuroyuki-kaze/remote.git
```

Step 2. Install the binary
```sh
cd remote
go install
```

Easy, right?

## Usage

Check out the example `remote.toml` file in the repo. It's pretty self-explanatory.

To use `remote`, you need to create a `remote.toml` file in the root of your project. This file will contain all the run scripts you want to use.

All scripts must be under the `[scripts]` table. The key can be any of the two format:
- Array: `run = ["python", "-s", "src.main"]`
- Dictionary: `build = {command = ["python", "-s", "src.build_main"], next = "run"}`

In the dictionary form, the `next` key is optional.

Once you have defined your `remote.toml` file, you can run your script by calling `remote <script_name>`.

For example, if you have a script called `run`, you can run it by calling `remote run`.

# License

[Unlicensed. Use it for whatever you want.](https://unlicense.org/)