# Caching Proxy CLI 🚀

## Overview

This project implements a caching proxy server in Go, allowing users to forward requests to an origin server while caching the responses for efficient data retrieval. This project is inspired by the [roadmap.sh](https://roadmap.sh/projects/caching-server).

## Features

- Start a caching proxy server.
- Forward requests to a specified origin server.
- Cache responses for faster subsequent retrievals.
- Clear the cache with a simple command.
- CLI interface built using Cobra.

## Requirements

- Go (version 1.16 or later)

## Installation

1. Clone the repository:

   ```bash
   git clone <repository-url>
   cd <repository-directory>
2. Build the project:
   ```bash
   make build
   ```
3. Run the proxy server:
   ```bash
   cacheproxy --port <port-number> --origin <origin-url>
   ```
4. Clear Cache
   ```bash
   cacheproxy --clear-cache
   ```

## Usage

Send requests to your caching proxy server:
```bash
curl -v http://localhost:<port>/products
```
This will fetch the products from the origin server and cache the response. Subsequent requests for the same endpoint will return the cached response, improving performance.
```
< HTTP/1.1 200 OK
< X-Cache: HIT ( If the response is cached you get HIT and if the response is fetched you get MISS )
< Date: Sun, 20 Oct 2024 09:13:50 GMT
< Content-Type: text/plain; charset=utf-8
< Transfer-Encoding: chunked
```

I'm storing all the responses or cache in the file system. I'm sure there are many more optimizations that can be made to the application. This is the most rudimentary solution and the very first approach I could think of. I'll test this implementation and see where it falls short, and eventually, I may build a better version. 🚀✨
Feel free to raise issues or contribute to building a better approach! 🙌💡
