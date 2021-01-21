# GDAL for Go

## About

The gdal.go package provides a go wrapper for GDAL, the Geospatial Data Abstraction Library. More information about GDAL can be found at http://www.gdal.org

## Installation
```sh
go get github.com/dewberry/gdal
```

## Compatibility

This software has been tested most recently on Alpine Linux, GDAL version 3.2.0.

## Examples
See the [examples directory](examples) for some starting points.

## Status 

GDAL 2.3.2 (3/08/2019):
 - The majority of GDAL functionality exposed by the C API is available, as well as much of the OGR API.
 - Most functionality is not covered by tests or benchmarks.

GDAL 3+ (10/21/2020):
 - GDAL functionality exposed by the C API that has been removed when GDAL 3 was released has also been removed from this Go API.