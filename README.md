# Tin
Tin is a collection of utility functions, accessible through CLI.

## Base64 Encoder / Decoder
```bash
sn b64 [-flags] input
```

Supported flag are:
- `-d` decode the input (encode is default)
- `-f` read a file as the input
- `-o` write result to output file
- `-u` use base64URL encoding and decoding
- `-h` help

### Examples
```bash
# base64 encode string
sn b64 "hello, world"

# base64 decode string, writing file to output
sn b64 -d -o output.txt aGVsbG8sIHdvcmxk

# base64 encode from file
sn b64 -f output.txt
```

## Diff Tool
```bash
sn diff [file1] [file2]
```


# What to add?
- [ ] JWT generator / decoder
- [ ] JSON parser / formatter
- [ ] XML parser / formatter

and more...