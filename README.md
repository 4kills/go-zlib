# **zlib** for go

<a href="https://choosealicense.com/licenses/zlib/">
<img src="https://img.shields.io/badge/license-zlibLicense-blue.svg" alt="License: zlibLicense">
</a>

This ultra fast **go zlib library** wraps the original zlib library written in C by Jean-loup Gailly and Mark Adler using cgo. 

**It offers considerable performance benefits compared to the standard go zlib library**, as the [benchmarks](#benchmarks) show.

This library is designed to be completely and easily interchangeable with the go standard zlib library. <ins>You won't have to rewrite or modify a single line of code!</ins> Checking if this library works for you is as easy as changing [imports](#import)!

This library also offers <ins>blazing fast convenience methods</ins> that can be used as a clean, alternative interface to that provided by the go standard library. (See [usage](#usage)).

**WARNING:** As of now, this library does not support streaming and can thus only compress as many bytes as fit into a single buffer (as much as your memory can hold) in one go (no pun intended). 
You can still chunk the data and put it back together after decompressing. 

## Table of Contents

- [Features](#features)
- [Installation](#installation)
  - [Prerequisites](#prerequisites)
  - [Download and Installation](#download-and-installation)
  - [Import](#import)
- [Usage](#usage)
  - [Compress](#compress)
  - [Decompress](#decompress)
- [Notes](#notes)
- [Benchmarks](#benchmarks)
  - [Compression](#compression)
  - [Decompression](#decompression)
- [License](#license)
- [Links](#links)

# Features

- [x] zlib compression / decompression
- [x] A variety of different `compression strategies` and `compression levels` to choose from 
- [x] Seamless interchangeability with the go standard zlib library 
- [x] Alternative, super fast convenience methods for compression / decompression
- [x] Benchmarks with comparisons to the go standard zlib library
- [ ] Custom, user-defined dictionaries
- [ ] More customizable memory management 
- [ ] Support streaming of data to compress/decompress data. 

# Installation

## Prerequisites 

In order to use this library with your go source code, you must be able to use the go tool [cgo](https://golang.org/cmd/cgo/), which in turn requires a GCC compiler.

If you are on **Linux**, there is a good chance you already have GCC installed, otherwise just get it with your favorite package manager.

If you are on **MacOS**, Xcode - for instance - supplies the required tools.

If you are on **Windows**, you will need to install a GNU environment.
I can recommend [tdm-gcc](https://jmeubank.github.io/tdm-gcc/) which is based
off of MinGW. Please note that [cgo](https://golang.org/cmd/cgo/) requires the 64-bit version (as stated [here](https://github.com/golang/go/wiki/cgo#windows)). 

## Download and Installation

To get the most recent stable version just type: 

```shell script 
$ go get github.com/4kills/zlib
```

You may also use go modules (available since go 1.11) to get the version of a specific branch or tag if you want to try out or use experimental features. However, beware that these versions are not necessarily guaranteed to be stable or thoroughly tested.

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

### Like with the standard library: 

```go
var b bytes.Buffer              // use any writer
w := zlib.NewWriter(&b)         // create a new zlib.Writer, compressing to b
defer w.Close()                 // don't forget to close this
w.Write([]byte("uncompressed")) // put in any data as []byte  
```

### Or alternatively: 

```go 
w := zlib.NewWriter(nil)                     // requires no writer if WriteBytes is used
defer w.Close()                              // always close when you are done with it
c, _ := w.WriteBytes([]byte("uncompressed")) // compresses input & returns compressed []byte 
```

## Decompress

### Like with the standard library: 

```go
b := bytes.NewBuffer(compressed) // reader with compressed data
r, err := zlib.NewReader(&b)     // create a new zlib.Reader, decompressing from b 
defer r.Close()                  // don't forget to close this either
io.Copy(os.Stdout, r)            // read all the decompressed data and write it somewhere
// or:
// r.Read(someBuffer)            // can also be done directly
```

### Or alternatively: 

```go 
r := zlib.NewReader(nil)         // requires no reader if ReadBytes is used
defer r.Close()                  // always close or bad things will happen
dc, _ := r.ReadBytes(compressed) // decompresses input & returns decompressed []byte 
```

# Notes

- **Do NOT use the <ins>same</ins> Reader / Writer across multiple threads <ins>simultaneously</ins>.** You can do that if you **sync** the read/write operations, but you could also create as many readers/writers as you like - for each thread one, so to speak. This library is generally considered thread-safe.

- **Always `Close()` your Reader / Writer when you are done with it** - especially if you create a new reader/writer for each decompression/compression you undertake (which is generally discouraged anyway). As the C-part of this library is not subject to the go garbage collector, the memory allocated by it must be released manually (by a call to `Close()`) to avoid memory leakage.

- **`HuffmanOnly` does NOT work as with the standard library**. This is the only exception from the philosophy to make this library interchangeable with the standard library. If you want to use 
`HuffmanOnly`, refer to the `NewWriterLevelStrategy()` constructor function. However, your existing code won't break by leaving `HuffmanOnly` as argument to `NewWriterLevel()`, it will just use the default compression strategy and compression level 2.  

- Memory Usage: `Compressing` requires ~256KiB of additional memory during execution, while `Decompressing` requires ~39 KiB of additional memory during execution. 
So if you have 8 simultaneous `WriteBytes` working from 8 Writers across 8 threads, your memory footprint from that alone will be about ~2MiByte.

- You are strongly encouraged to use the same Reader / Writer for multiple Decompressions / Compressions as it is not required nor beneficial in any way, shape or form to create a new one every time. The contrary is true: It is more performant to reuse a reader/writer. Of course, if you use the same reader/writer multiple times, you do not need to close them until you are completely done with them (perhaps only at the very end of your program).  

# Benchmarks

These benchmarks were conducted with "real-life-type data" to ensure that these tests are most representative for an actual use case in a practical production environment.
As the zlib standard has been traditionally used for compressing smaller chunks of data, I have decided to follow suite by opting for Minecraft client-server communication packets, as they represent the optimal use case for this library. 

To that end, I have recorded 930 individual Minecraft packets, totalling 11,445,993 bytes in umcompressed data and 1,564,159 bytes in compressed data.
These packets represent actual client-server communication and were recorded using [this](https://github.com/haveachin/infrared) software.

The benchmarks were executed on different hardware and operating systems, including AMD and Intel processors, as well as all the supported operating systems (Windows, Linux, MacOS). All of the benchmarked functions/methods were executed hundreds of times and the numbers you are about to see are the averages over all these executions.

These benchmarks compare this library (blue) to the go standard library (yellow) and show that this library performs better in all cases. 

- <details>
  
    <summary> (A note regarding testing on your machine) </summary>
  
    Please note that you will need an Internet connection for some of the benchmarks to function. This is because these benchmarks will download the mc packets from [here](https://github.com/4kills/zlib_benchmark) and temporarily store them in memory for the duration of the benchmark tests, so this repository won't have to include the data in order save space on your machine and to make it a lightweight library.
  
  </details>

## Compression

![compression total](https://i.imgur.com/CPjYJQQ.png)

This chart shows how long it took for the methods of this library (blue) and the standard library (yellow) to compress **all** of the 930 packets (~11.5MB) on different systems in nanoseconds. Note that the two rightmost data points were tested on **exactly the same** hardware in a dual-boot setup and that Linux seems to generally perform better than Windows.

![compression relative](https://i.imgur.com/dK6i9Ij.png)

This chart shows the time it took for this library's `Write` (blue) to compress the data in nanoseconds, as well as the time it took for the standard library's `Write` (WriteStd, yellow) to compress the data in nanoseconds. The vertical axis shows percentages relative to the time needed by the standard library, thus you can see how much faster this library is. 

For example: This library only needed ~88% of the time required by the standard library to compress the packets on an Intel Core i5-6600K on Windows. 
That makes the standard library **~13.6% slower** than this library. 

## Decompression

![compression total](https://i.imgur.com/Ef3xM6Q.png)

This chart shows how long it took for the methods of this library (blue) and the standard library (yellow) to decompress **all** of the 930 packets (~1.5MB) on different systems in nanoseconds. Note that the two rightmost data points were tested on **exactly the same** hardware in a dual-boot setup and that Linux seems to generally perform better than Windows.

![dcompression relative](https://i.imgur.com/UQ7dKpA.png)

This chart shows the time it took for this library's `Read` (blue) to decompress the data in nanoseconds, as well as the time it took for the standard library's `Read` (ReadStd, Yellow) to decompress the data in nanoseconds. The vertical axis shows percentages relative to the time needed by the standard library, thus you can see how much faster this library is. 

For example: This library only needed a whopping ~65% of the time required by the standard library to decompress the packets on an Intel Core i5-6600K on Windows. 
That makes the standard library a substantial **~53.8% slower** than this library.

# License

```txt
  Copyright (c) 1995-2017 Jean-loup Gailly and Mark Adler
  Copyright (c) 2020 Dominik Ochs

This software is provided 'as-is', without any express or implied
warranty.  In no event will the authors be held liable for any damages
arising from the use of this software.

Permission is granted to anyone to use this software for any purpose,
including commercial applications, and to alter it and redistribute it
freely, subject to the following restrictions:

  1. The origin of this software must not be misrepresented; you must not
     claim that you wrote the original software. If you use this software
     in a product, an acknowledgment in the product documentation would be
     appreciated but is not required.
 
  2. Altered source versions must be plainly marked as such, and must not be
     misrepresented as being the original software.
  
  3. This notice may not be removed or altered from any source distribution.
```

# Links 

- Original zlib by Jean-loup Gailly and Mark Adler: 
    - [github](https://github.com/madler/zlib) 
    - [website](https://zlib.net/)
- Go standard zlib by the Go Authors: 
    - [github](https://github.com/golang/go/tree/master/src/compress/zlib)
