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


## Timezone Converter
```bash
sn tz [-flags] [time] [date]
```

Supported flags are:
- `-f` timezone name to convert from (Local is default)
- `-t` timezone name to convert to (UTC is default)
- `-h` help

Arguments:
`time` argument supports both 12-hour and 24-hour formats
`date` argument format is `dd/mm/yyyy`

### Examples
```bash
# default converts local current time to UTC
sn tz

# specify from and to flags
sn tz -f Asia/Singapore -t Canada/Pacific

# specify time to convert
sn tz 1234

# specify date and time
sn tz 12:34pm 19/11/2023
```

# What to add?
- [x] Timezone converter
- [ ] JWT generator / decoder
- [ ] JSON parser / formatter
- [ ] XML parser / formatter

and more...