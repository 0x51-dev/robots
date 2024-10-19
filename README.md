# Robots

This is a simple Go library for parsing and handling `robots.txt` files. The library allows developers to check whether
a given user-agent is allowed or disallowed from accessing specific URLs on a site.

## Some Features

- Full compliance with RFC 9309 (the latest specification for robots.txt).
- Parsing robots.txt files to extract rules for various user-agents.
- Determining whether a URL is allowed or disallowed based on the rules in robots.txt.
- Handling of the following directives: `User-agent`, `Allow`, `Disallow`, `Sitemap`, etc.
- Support for comments, case-insensitivity, and empty lines as per the specification.

## Installation

To use the library, install it using:
```
go install github.com/0x51-dev/robots
```

## References

- [RFC9309](https://datatracker.ietf.org/doc/html/rfc9309)
