# What
For concurrency models that I need to use in Go, currently only has a basic worker pool.
More for understanding purposes.
Has example(s) on how to use.

# Why
1. Needed some kind of worker pool for my other projects.
2. One of Go's selling points is its concurrency, so I will likely need this sort of construct a lot for future projects
as well.
3. Concurrency is notoriously difficult to implement/test, so I decided to create an isolated area to create and test these

# Requirements
1. Go 1.25.0

# How to run
1. run ./build.sh
2. run whatever examples such as:
<br>
`fib.exe 33`
to get the 33rd fibonacci number


# Possible Extensions
Multiplexer/Load balancer (Fan in, fan out)
Pipelining
Actor Model
Events
Timer
PubSub
Tests

# Possible More Examples
Theoretical examples:
Map reduce

Practical examples:
Web page downloader
TCP connection pool

