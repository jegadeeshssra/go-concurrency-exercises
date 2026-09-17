Based strictly on the attached content, there are **four explicit exercise points**, so I’ll frame only those four. Each question is designed to make you implement the concepts explained before that exercise. 

---

## 1. Exercise — Canceling the Generator

### Goal

Implement a **cancelable generator using `context.Context`**, applying:

* `context.Background()`
* `context.WithCancel()`
* `ctx.Done()`
* `ctx.Err()`
* goroutines
* channel cancellation
* channel closing
* preventing a leaked generator goroutine
* parent/child context cancellation behavior

### Question

You already have a number generator:

```go
func generate(start, stop int) <-chan int {
    // implement
}
```

The generator sends every number from `start` to `stop` through a channel.

The problem is that if the consumer stops reading early, the generator can remain blocked forever trying to send another value.

Modify the generator so that it accepts a context:

```go
func generate(ctx context.Context, start, stop int) <-chan int {
    // implement
}
```

The generator must:

1. Run in a separate goroutine.
2. Send numbers from `start` to `stop`.
3. Use `select` when sending values so cancellation can interrupt a blocked send.
4. Stop immediately when `ctx.Done()` is closed.
5. Close the output channel when the generator exits.
6. Allow the caller to cancel the generator using `cancel()`.
7. Ensure no generator goroutine remains blocked if the consumer stops reading.

Write a client that:

* creates a cancelable context,
* generates numbers from `1` to `1,000,000`,
* consumes only the first 10 numbers,
* cancels the context,
* verifies that the generator terminates.

Then modify the program to create a **parent context and child context** and demonstrate that:

* canceling the child does not cancel the parent;
* canceling the parent also cancels the child.

---

# 2. Exercise — Pipeline with Context

### Goal

Build a **cancelable concurrent pipeline** using:

* context propagation
* `ctx.Done()`
* pipeline stages
* goroutines
* channels
* cancellation-aware sends
* cancellation-aware receives
* proper channel closing
* preventing goroutine leaks
* context timeouts/deadlines
* parent/child timeout behavior

### Question

Build a number-processing pipeline with three stages:

```text
Generator → Processor → Consumer
```

The generator produces numbers:

```text
1, 2, 3, ... 1000000
```

The processor performs an expensive operation:

```go
func process(n int) int {
    time.Sleep(10 * time.Millisecond)
    return n * 2
}
```

Define the stages as:

```go
func generate(ctx context.Context, start, stop int) <-chan int

func process(ctx context.Context, in <-chan int) <-chan int

func consume(ctx context.Context, in <-chan int)
```

Your implementation must satisfy the following:

1. Every stage receives the same `context.Context`.
2. The generator must stop when `ctx.Done()` is triggered.
3. The processor must stop when `ctx.Done()` is triggered.
4. Channel sends must not remain blocked when cancellation occurs.
5. Channels must be closed by the goroutine responsible for producing their values.
6. The consumer must stop reading when the context is canceled.
7. No goroutine should remain blocked after cancellation.

Run the pipeline using:

```go
ctx, cancel := context.WithTimeout(
    context.Background(),
    100*time.Millisecond,
)
defer cancel()
```

Observe how the pipeline terminates when the deadline is reached.

Then create:

```go
parentCtx, parentCancel := context.WithTimeout(...)
childCtx, childCancel := context.WithTimeout(parentCtx, ...)
```

Run the pipeline with the child context and demonstrate that the **shorter timeout controls cancellation**.

Finally, print `ctx.Err()` after the pipeline terminates and distinguish between normal completion and context cancellation.

---

# 3. Exercise — Cancelable Worker

### Goal

Build a worker that supports **context cancellation and cancellation-triggered cleanup**, applying:

* `context.Context`
* `ctx.Done()`
* manual cancellation
* timeout cancellation
* `context.AfterFunc()`
* cleanup functions
* registering cancellation callbacks
* detaching a callback with the returned stop function
* handling cancellation without leaking goroutines

### Question

You are building a worker that acquires a resource before starting a long-running operation:

```go
func acquireResource() {
    fmt.Println("resource acquired")
}

func releaseResource() {
    fmt.Println("resource released")
}

func worker(ctx context.Context) {
    // implement
}
```

The worker should simulate a 100 ms operation.

Implement `worker()` so that:

1. It acquires the resource.
2. It performs work while respecting `ctx.Done()`.
3. If the context is canceled before the work completes, the worker stops.
4. The resource is released when cancellation occurs.
5. The cleanup logic should be registered using `context.AfterFunc()` rather than being tightly coupled to the worker's cancellation branch.
6. The worker must not leak a goroutine.

Then create a client with:

```go
ctx, cancel := context.WithTimeout(
    context.Background(),
    50*time.Millisecond,
)
defer cancel()
```

Verify that:

```text
resource acquired
resource released
```

occurs when the timeout cancels the worker.

### Extension

Add a situation where cleanup is **no longer necessary**.

Register the cleanup function:

```go
stopCleanup := context.AfterFunc(ctx, releaseResource)
```

Then call the returned function before cancellation.

Determine whether cleanup executes and inspect the boolean returned by the stop function.

Finally, test the case where the context has **already been canceled** before registering the cleanup callback.

---

# 4. Exercises — Scheduler + 1 More

### Goal

Build a context-aware scheduler that combines the chapter's concepts:

* periodic execution with a ticker
* context cancellation
* context timeout/deadline
* `ctx.Done()`
* ticker cleanup with `Stop()`
* cancellation callbacks with `context.AfterFunc()`
* preventing goroutine/resource leaks
* passing request/context information through a concurrent operation

### Question

Build a **periodic job scheduler** that executes a job every 50 ms:

```go
func runScheduler(ctx context.Context, interval time.Duration) {
    // implement
}
```

The job should take approximately 100 ms:

```go
func work() {
    fmt.Println("work started")
    time.Sleep(100 * time.Millisecond)
    fmt.Println("work finished")
}
```

The scheduler must:

1. Create a ticker.
2. Execute work when a tick arrives.
3. Stop the ticker when the scheduler exits.
4. Listen for `ctx.Done()` and stop when the context is canceled.
5. Avoid leaving a goroutine running after cancellation.
6. Handle the fact that `work()` takes longer than the ticker interval.
7. Run for a limited amount of time using:

```go
ctx, cancel := context.WithTimeout(
    context.Background(),
    500*time.Millisecond,
)
defer cancel()
```

Then add a cleanup operation:

```go
func cleanup() {
    fmt.Println("scheduler cleanup")
}
```

Register it using `context.AfterFunc()` so that cleanup runs when the scheduler's context is canceled.

### Extension

Create a parent context with a **1-second timeout** and a child context with a **300 ms timeout**.

Run the scheduler using the child context.

Verify through the program's output that the scheduler stops according to the child's earlier deadline and that cleanup is triggered when the child context is canceled.

