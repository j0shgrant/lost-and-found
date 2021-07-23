# lost-and-found
lost-and-found is a CLI utility for finding AWS resources that are untagged, non-compliant, or in some way misbehaving.

## Installation
To install `laf`:
```shell
$ make
```
To validate your installation, run:
```shell
$ laf version
```
## Usage
To list EC2 instances by region:
```shell
$ laf ec2 --region <comma,separated,regions>
```

## Contributors
- [Josh Grant](https://github.com/j0shgrant)