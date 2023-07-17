# About `.output` Directory

This directory is used for storing any output files generated from deployment & test commands specified in the [Makefile](../Makefile). This includes:

- go modules cache when running the `make run*` command
- server build result when running the `make run*` command
- test coverage reports when running the `make test` command