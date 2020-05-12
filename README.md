# 6.824 Final Project

[![rootjalex](https://circleci.com/gh/rootjalex/smart-cache.svg?style=shield)](https://app.circleci.com/pipelines/github/rootjalex/smart-cache)

Final project for [Meia Alsup](https://www.linkedin.com/in/meiaalsup/), [Amir Farhat](https://github.com/amirfarhat), and [Alexander Root](https://rootjalex.github.io) for our [Distributed Systems](https://pdos.csail.mit.edu/6.824/) class. We are implementing a smart caching system, for use in large distributed databases, when there are patterns in the data accesses.

## Paper
Paper can be found [here](paper/824finalpaper.pdf).

## Running the tasks


To check out our benchmarks, from the src/ directory, run:

```
go run benchmarks.go
```

## Testing
We rely on Go's testing infrastructure. From the root of the repository, run:

```
cd src
go test ./...
```


## Built With

* [Go](https://golang.org) - go version go1.13.7 darwin/amd64

## Authors and Contributions

* [**Alexander Root**](https://rootjalex.github.io) - Data structures including
  Markov Chains and Caches, and Benchmarks.
* [**Meia Alsup**](https://meiaalsup.github.io) - Cache Master, Hashing, CI, Paper, System Diagram.
* [**Amir Farhat**](https://github.com/amirfarhat) - Clients, Workloads and Task Management.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Thanks to Professor [Morris](https://pdos.csail.mit.edu/~rtm/) and [Anish Athalye](https://www.anish.io) for their guidance and helpful conversations.
* Thanks to [Jacob Kahn](https://jacobkahn.me) for providing the original inspiration for this project.
