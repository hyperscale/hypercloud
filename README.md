Hyperpaas ![Last release](https://img.shields.io/github/release/hyperscale/hyperpaas.svg) 
=========

![HyperPaaS logo](https://cdn.rawgit.com/hyperscale/hyperpaas/master/_resources/hyperpaas.svg "HyperPaaS logo")


[![Go Report Card](https://goreportcard.com/badge/github.com/hyperscale/hyperpaas)](https://goreportcard.com/report/github.com/hyperscale/hyperpaas)

| Branch  | Status | Coverage |
|---------|--------|----------|
| master  | [![Build Status](https://img.shields.io/travis/hyperscale/hyperpaas/master.svg)](https://travis-ci.org/hyperscale/hyperpaas) | [![Coveralls](https://img.shields.io/coveralls/hyperscale/hyperpaas/master.svg)](https://coveralls.io/github/hyperscale/hyperpaas?branch=master) |
| develop | [![Build Status](https://img.shields.io/travis/hyperscale/hyperpaas/develop.svg)](https://travis-ci.org/hyperscale/hyperpaas) | [![Coveralls](https://img.shields.io/coveralls/hyperscale/hyperpaas/develop.svg)](https://coveralls.io/github/hyperscale/hyperpaas?branch=develop) |

HyperPaaS a Cloud Application Platform based on Docker Swarm.

Install
-------

### Docker

```shell
docker pull hyperscale/hyperpaas
```

### MacOS

Install dependencies with glide:
```shell
glide install
```

Build hyperpaas:
```shell
make build
```

Run hyperpaas
```shell
./hyperpaas
```

Documentation
-------------

[HperPaaS API Reference](https://hyperscale.github.io/hyperpaas/)

License
-------

HperPaaS is licensed under [the MIT license](LICENSE.md).
