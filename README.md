# collector-of-files

**This entire repository, including this README file except for this sentence, was created by AI with no human input.**

A simple tool to collect files from a directory and its subdirectories, excluding files and directories specified in a `.gitignore` file. It's to make feeding chunks of code to AI models easier.

## Usage

```
Usage: combine-code --path <start_path> --includeExtensions <comma_separated_extensions> [--output <output_file>] [--excludeDirs <comma_separated_dirs>]
Example: combine-code --path /path/to/start --includeExtensions go,js --excludeDirs node_modules,vendor
```

## Example

```
$ go run main.go --path ~/Projects/collector-of-files --includeExtensions go,md --output combined_code.md
Combined code has been written to combined_code.md
```
