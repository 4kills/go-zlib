<h1 align="center"> <b>zlib</b> for go </h1>

<p align="center">

<a href="http://unlicense.org/">
<img src="https://img.shields.io/badge/license-Unlicense-blue.svg" alt="License: Unlicense">
</a>

</p>

Native golang zlib implementation using cgo and the original zlib library written in C by Jean-loup Gailly and Mark Adler. 

## Table of Contents

- [Features](#features)
- [Installation](#installation)
  - [Prerequisites](#prerequisites)
  - [Download and Installation](#download-and-installation)
  - [Import](#import)
- [Usage](#usage)
  - [sub3](#sub3)
  - [sub4](#sub4)
- [License](#license)
- [Links](#links)

# Features

- [x] zlib compression / decompression
- [x] Seamless interchangeability with the go standard zlib library 
- [x] Alternative, super fast convenience methods for compression / decompression
- [x] Benchmarks with comparisons to the go standard zlib library
- [ ] Custom, user-defined dictionaries

# Installation

## Prerequisites 

In order to use this library with your go source code, you must be able to use the go tool [cgo](https://golang.org/cmd/cgo/), which in turn requires a GCC compiler.

If you are on **Linux** or **MacOS**, you are already good to go.

If you are on **Windows**, you will need to install a GCC compiler. 
I can recommand [tdm-gcc](https://jmeubank.github.io/tdm-gcc/) which is based
off of WinGW. Please note that [cgo](https://golang.org/cmd/cgo/) requires the 64-bit version (as stated [here](https://github.com/golang/go/wiki/cgo#windows)). 

## Download and Installation

To get the most recent stable version just type: 

```shell script 
$ go get github.com/4kills/zlib
```

You may also use go modules (available since go 1.11) to get the version of a specific branch or tag if you want to try out or use experimental features. But beware that these versions are not necessarily guaranteed to be stable or thoroughly tested.

## Import

This library is designed in a way to make it easy to swap it out for the go standard zlib library. Therefore you should only need to change imports and not a single line of your written code. 

Just remove: 

~~import compress/zlib~~

and use instead: 
 
```go
import "github.com/4kills/zlib"
```

If there are any problems with your existing code after this step, please let me know. 

# Usage

This library can be used exactly like the [go standard zlib library](https://golang.org/pkg/compress/zlib/) but it also adds additional methods to make your life easier.

## Compress


# Links 

zlib by Jean-loup Gailly and Mark Adler: [github](https://github.com/madler/zlib) and [website](https://zlib.net/). 
