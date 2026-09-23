Based strictly on the attached content, there are **four explicit exercise points**. I’ve framed only those exercises, and each one is designed to make you implement the concepts explained leading up to that exercise. 

---

# 1. Exercise — From Channel to Wait Group

### Goal

Practice replacing a **done channel** with a `sync.WaitGroup`, and understand:

* `WaitGroup.Add()`
* `WaitGroup.Done()`
* `WaitGroup.Wait()`
* `WaitGroup.Go()`
* the internal counter model
* waiting for multiple goroutines
* why a `WaitGroup` must be shared rather than copied

### Question

You are given a program that starts several independent tasks and uses a channel to wait until each task finishes:

```go
func runTask(id int, done chan<- struct{}) {
    time.Sleep(time.Duration(id*50) * time.Millisecond)
    fmt.Printf("task %d done\n", id)
    done <- struct{}{}
}
```

The client currently looks like:

```go
func main() {
    done := make(chan struct{}, 3)

    for i := 1; i <= 3; i++ {
        go runTask(i, done)
    }

    for i := 0; i < 3; i++ {
        <-done
    }

    fmt.Println("all tasks done")
}
```

Rewrite the program using `sync.WaitGroup`.

Your implementation must:

1. Start all three tasks concurrently.
2. Use `wg.Add()` before starting the goroutines.
3. Call `wg.Done()` when each task finishes.
4. Use `wg.Wait()` instead of receiving from the `done` channel.
5. Modify the program so the `WaitGroup` is passed correctly to the worker rather than copied.
6. Then rewrite the same solution using `WaitGroup.Go()` instead of `Add()` + `Done()`.

The final program should print:

```text
task 1 done
task 2 done
task 3 done
all tasks done
```

The order of the task-completion messages is not guaranteed.

---

# 2. Exercise — Concurrent Group

### Goal

Build an abstraction that hides `sync.WaitGroup` from its caller and practice **encapsulation of concurrency synchronization**.

Concepts to apply:

* `sync.WaitGroup`
* concurrent execution
* `Add()` / `Done()` / `Wait()`
* variadic functions
* wrapper functions
* wrapper types
* keeping synchronization details internal
* reusing a concurrency abstraction

### Question

You need to create a reusable component that allows callers to register functions and execute all of them concurrently.

Create:

```go
type ConcurrentGroup struct {
    // internal fields
}
```

and:

```go
func NewConcurrentGroup() *ConcurrentGroup
```

The type must provide:

```go
func (cg *ConcurrentGroup) Add(fn func())
func (cg *ConcurrentGroup) Run()
```

The behavior should be:

1. `Add()` registers a function without executing it.
2. `Run()` executes all registered functions concurrently.
3. `Run()` does not return until every function has completed.
4. The caller must **not** need to know that a `sync.WaitGroup` is being used.
5. The same group should be reusable, so calling `Run()` a second time should execute the registered functions again.
6. Synchronization state must remain encapsulated inside `ConcurrentGroup`.

Test it with:

```go
func work(id int) func() {
    return func() {
        time.Sleep(50 * time.Millisecond)
        fmt.Printf("work %d done\n", id)
    }
}
```

Create a group, add three functions, run them, and then run the same group again.

Measure the execution time of each `Run()` call and verify that the three 50 ms operations execute concurrently rather than sequentially.

---

# 3. Exercise — Waiting for Worker

### Goal

Practice using **multiple `Wait()` calls** from different goroutines and understand that:

* multiple goroutines can wait on the same `WaitGroup`
* all waiters unblock when the counter reaches zero
* the order in which waiters continue is not guaranteed
* `WaitGroup` can coordinate multiple observers of the same completion event

### Question

You are building a batch-processing system.

A single worker performs a 100 ms task:

```go
func worker(wg *sync.WaitGroup) {
    defer wg.Done()

    time.Sleep(100 * time.Millisecond)
    fmt.Println("worker finished")
}
```

Three different parts of the application need to be notified when the worker finishes:

```text
Worker
  │
  ├── Waiter 1
  ├── Waiter 2
  └── Waiter 3
```

Implement a program that:

1. Creates a `sync.WaitGroup`.
2. Adds one worker to the wait group.
3. Starts the worker.
4. Starts **three separate waiter goroutines**.
5. Each waiter calls `wg.Wait()`.
6. Each waiter prints its own completion message after `Wait()` returns.
7. `main()` must also call `wg.Wait()`.

For example:

```text
worker finished
waiter 1 done
waiter 2 done
waiter 3 done
main waiter done
```

The exact order of the waiter messages must **not** be relied upon.

Then modify the program so that the worker performs multiple pieces of work using the same wait group and verify that **none of the waiters proceed until the entire group has completed**.

---

# 4. Exercise — Concurrent Group With Panic Handling

### Goal

Build a concurrent execution abstraction that handles panics **inside the goroutines where they occur**, while still allowing the caller to know whether any worker panicked.

Concepts to apply:

* `sync.WaitGroup`
* concurrent goroutines
* `WaitGroup.Go()`
* `defer`
* `recover()`
* per-goroutine panic recovery
* synchronization of worker completion
* communicating the panic result safely
* understanding why `recover()` in `main()` cannot catch a panic from another goroutine

### Question

Create a reusable function:

```go
func RunSafe(funcs ...func()) (panicked bool)
```

It should execute all supplied functions concurrently and wait for all of them to finish.

Some functions may panic:

```go
func normalWork(id int) func() {
    return func() {
        time.Sleep(50 * time.Millisecond)
        fmt.Printf("work %d completed\n", id)
    }
}

func failingWork(id int) func() {
    return func() {
        time.Sleep(50 * time.Millisecond)
        panic(fmt.Sprintf("work %d failed", id))
    }
}
```

Your implementation must:

1. Run all functions concurrently.
2. Wait for every function to finish.
3. Recover panics **inside the goroutine that executes the function**.
4. Prevent a worker panic from terminating the entire program.
5. Report whether at least one worker panicked.
6. Ensure `RunSafe()` doesn't return until all workers have completed, including workers that panic.
7. Do not attempt to recover the worker panic from `main()`.

Test it with:

```go
RunSafe(
    normalWork(1),
    failingWork(2),
    normalWork(3),
    failingWork(4),
)
```

The program should complete all four goroutines and report that at least one panic occurred.

### Additional challenge

Instead of returning only:

```go
bool
```

modify the implementation to collect the panic values from the workers and return them to the caller:

```go
func RunSafe(funcs ...func()) []any
```

For example, the caller should be able to determine which workers produced panic values.

