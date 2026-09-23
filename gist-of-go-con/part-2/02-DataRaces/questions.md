Based strictly on the attached content, there are **four explicit exercise points**, so I’ll frame only those four. Each exercise is designed to make you implement the concepts explained before that exercise. 

---

# 1. Exercise — Spot the Race

### Goal

Learn to identify and diagnose **data races** before fixing them.

Concepts to apply:

* concurrent access to shared variables
* data races
* concurrent map modification
* why apparently simple operations such as `counter[word]++` are not atomic
* `sync.WaitGroup`
* Go's race detector
* safe communication through channels

### Question

You are given a concurrent word-frequency counter:

```go
func main() {
    in := generate(1000, 3)

    var wg sync.WaitGroup
    wg.Add(4)

    counter := map[string]int{}

    count := func() {
        defer wg.Done()

        for word := range in {
            counter[word]++
        }
    }

    for range 4 {
        go count()
    }

    wg.Wait()

    fmt.Println(counter)
}
```

The program sometimes crashes with:

```text
fatal error: concurrent map writes
```

Your task is to **investigate the program rather than immediately fixing it**.

1. Identify which variable is being accessed concurrently.
2. Identify the exact operation that creates the data race.
3. Explain why `counter[word]++` is not an atomic operation.
4. Run the program with Go's race detector:

```bash
go run -race .
```

5. Use the race detector output to identify the conflicting goroutines and operations.
6. Modify the program so that the shared map is no longer concurrently modified.
7. Verify the corrected implementation using:

```bash
go run -race .
```

The final implementation must produce the correct word frequencies without reporting a data race.

**Additional requirement:** Implement the corrected version once using a synchronization primitive and once using channel-based communication, then compare the two approaches.

---

# 2. Exercise — Concurrent-Safe Counter

### Goal

Build a **concurrent-safe shared counter** using a mutex and learn how to protect a critical section.

Concepts to apply:

* shared mutable state
* data races
* `sync.Mutex`
* `Lock()`
* `Unlock()`
* critical sections
* pointer-based synchronization
* `defer` for reliable unlocking
* race detector

### Question

Build a concurrent-safe counter type:

```go
type Counter struct {
    // internal fields
}
```

It should provide:

```go
func (c *Counter) Add(key string)
func (c *Counter) Get(key string) int
```

The counter internally stores:

```go
map[string]int
```

`Add()` should increment the count for a key, while `Get()` should return its current value.

Now create **10 goroutines**, and have each goroutine perform 10,000 increments:

```go
counter.Add("requests")
```

Wait for all goroutines using a `sync.WaitGroup`.

Your implementation must:

1. Allow all goroutines to share the same counter.
2. Prevent concurrent map access from causing a data race.
3. Protect the appropriate critical section using `sync.Mutex`.
4. Ensure the mutex is always unlocked, even if the protected operation later changes.
5. Pass the counter as a shared pointer rather than copying the synchronization state.
6. Produce the expected final count:

```text
100000
```

Run the program with:

```bash
go run -race .
```

There must be **no race detector warnings**.

### Additional challenge

Add:

```go
func (c *Counter) Reset()
```

and make sure `Reset()` can safely execute concurrently with `Add()` and `Get()`.

---

# 3. Exercise — Counter with RWMutex

### Goal

Implement a counter where **multiple readers can operate concurrently while writes remain exclusive**, using:

* `sync.RWMutex`
* `RLock()`
* `RUnlock()`
* `Lock()`
* `Unlock()`
* multiple concurrent readers
* exclusive writers
* critical sections
* measuring the effect of concurrent reads

### Question

Build a concurrent word-frequency store:

```go
type WordCounter struct {
    // internal fields
}
```

Provide:

```go
func (wc *WordCounter) Add(word string)
func (wc *WordCounter) Get(word string) int
```

The program will have:

* **1 writer**
* **4 readers**

The writer should perform 100 updates:

```go
word := randomWord(3)
wc.Add(word)
```

Each reader should perform 100 lookups:

```go
wc.Get(randomWord(3))
```

Each operation should simulate work by sleeping for:

```go
time.Sleep(time.Millisecond)
```

Your implementation must:

1. Protect writes using `Lock()` / `Unlock()`.
2. Protect reads using `RLock()` / `RUnlock()`.
3. Allow the four readers to execute concurrently with each other.
4. Prevent readers from accessing the map while a writer is modifying it.
5. Use a `sync.WaitGroup` to wait for all five goroutines.
6. Verify the implementation with:

```bash
go run -race .
```

Measure the total execution time.

Then create a second implementation using a regular `sync.Mutex` for both reads and writes.

Run both implementations with the same workload and compare their behavior.

Your goal is to observe why an `RWMutex` can allow concurrent readers while a regular mutex serializes them.

---

# 4. No Exercise Section After "Channel as Mutex"

There is **no explicit `Exercise:` / `Exercises:` marker after the Channel as Mutex section** in the attached content.

Therefore, according to your rule, I would **not create an additional exercise for this section**.

The last exercise remains:

> **Exercise: Counter with RWMutex**

The **Channel as Mutex** section is explanatory content showing that a buffered channel with capacity 1 can act as a mutual-exclusion mechanism, but since the source does not mark an exercise there, I’m not adding one. 
