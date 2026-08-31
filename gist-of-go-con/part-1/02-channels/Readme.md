This section builds on channel communication and focuses on **how goroutines signal completion, coordinate without deadlocking, control concurrency, and handle special channel states**.

### 1. End-of-Data Signaling

**Concept:** A receiver needs a way to know when a producer has finished sending values.

* Using a special value such as `"__EOF__"` can signal the end, but it is unsafe when that value can also be legitimate data.
* Go provides a proper mechanism: **closing the channel**.
* The writer closes the channel after sending all values.
* The reader can detect closure using the two-value receive:

```go
value, ok := <-channel
```

* `ok == true` → channel is open and a value was received.
* `ok == false` → channel is closed and no more values will arrive.
* Receiving from a closed channel is safe and repeatedly returns the zero value with `ok == false`.

**Scenario:** A goroutine produces comma-separated words while another goroutine consumes and filters them. The consumer must know when the producer has finished.

---

### 2. Closing Channels

**Concept:** Closing a channel is a **signal that no more values will be sent**.

Important rules:

* A channel can only be closed once.
* Sending to a closed channel causes a panic.
* Closing an already closed channel causes a panic.
* The **writer/producer** should close the channel.
* A receiver should not close a channel it does not own.
* If there are multiple writers, they must coordinate because one writer closing the channel can cause the others to panic.

**Important distinction:**

> Channels do not need to be closed for garbage collection. Close a channel when the receiver needs to know that the producer has finished.

**Scenario:** A producer sends strings to a consumer and closes the channel when all strings have been sent.

---

### 3. Channel Iteration

**Concept:** `for range` provides a simpler way to consume values until a channel is closed.

Instead of:

```go
for {
    value, ok := <-stream
    if !ok {
        break
    }
    // process value
}
```

you can write:

```go
for value := range stream {
    // process value
}
```

The loop automatically terminates when the channel is closed and all buffered values have been consumed.

**Scenario:** A producer sends a sequence of values, closes the channel, and a consumer processes every value using `for range`.

---

### 4. Done Channel

**Concept:** A **done channel** can be used to signal that a goroutine has completed its work.

The worker sends a value:

```go
done <- struct{}{}
```

The waiting goroutine receives it:

```go
<-done
```

If multiple goroutines are running, the receiver must know how many completion signals to wait for.

**Scenario:** Several independent "talker" goroutines process different phrases concurrently. `main()` must wait for every talker to finish without using `time.Sleep()` or `sync.WaitGroup`.

---

### 5. Preventing Deadlocks

**Concept:** Concurrent operations can deadlock when goroutines are waiting on each other in a circular dependency.

Example:

```text
main → waits for done
worker → waits for out
```

If `worker` cannot send to `out` until someone receives from it, while `main` refuses to receive from `out` until `done` arrives, neither can proceed.

**Solution:** Ensure the receiving and completion operations can proceed independently.

**Scenario:** A worker sends multiple results through `out` and signals completion through `done`. The program must consume the results while independently handling the completion signal so that the worker never becomes permanently blocked.

---

### 6. Buffered Channels

**Concept:** A buffered channel can store a fixed number of values without requiring an immediate receiver.

```go
stream := make(chan int, 3)
```

The buffer allows:

```text
Producer → [ 1 | 2 | 3 ] → Consumer
```

The producer can send until the buffer becomes full.

Once the buffer is full:

```go
stream <- value
```

blocks until space becomes available.

Useful properties:

```go
cap(stream) // buffer capacity
len(stream) // current number of buffered values
```

**Scenario:** A producer submits several orders while the consumer takes time to become ready. A buffered channel allows the producer to submit orders without immediately waiting for the consumer.

---

### 7. Semaphore Using a Buffered Channel

**Concept:** A buffered channel can act as a **semaphore** to limit the number of goroutines executing concurrently.

For a limit of `N`:

```go
sema := make(chan struct{}, N)
```

A goroutine acquires a slot before starting and releases it when finished.

This prevents launching unlimited concurrent work when the workload is large.

**Scenario:** A program has many images/files/API requests to process, but only **N operations should run concurrently**.

---

### 8. Inverse Semaphore

**Concept:** The same concurrency limit can be implemented using the opposite token flow.

Instead of initially filling the channel:

```text
[token] [token]
```

start with an empty channel.

Before starting work:

```go
sema <- struct{}{}
```

After finishing:

```go
<-sema
```

If the buffer is full, the next goroutine cannot start until another goroutine removes its token.

**Scenario:** Many independent requests need to run, but the system must allow at most two concurrent requests.

---

### 9. Fixed Worker Pool

**Concept:** A semaphore is not the only way to limit concurrency.

Instead of creating one short-lived goroutine per task, create a fixed number of **long-lived workers**:

```text
              ┌── Worker 1
Producer → pending ── Worker 2
              └── Worker 3
```

Workers repeatedly receive tasks from an input channel:

```go
for task := range pending {
    process(task)
}
```

The producer closes `pending` when there are no more tasks.

**Scenario:** A large collection of tasks needs processing, but exactly three worker goroutines should handle all of them.

**Key distinction:**

* **Semaphore:** many short-lived goroutines, concurrency is limited.
* **Worker pool:** fixed number of long-lived goroutines consume tasks.

---

### 10. Closing a Buffered Channel

**Concept:** Closing a buffered channel does **not discard values already stored in the buffer**.

If:

```go
stream := make(chan int, 2)

stream <- 1
stream <- 2

close(stream)
```

the receiver still gets:

```text
1, true
2, true
0, false
```

The buffered values are delivered first. Only after the buffer is empty does receiving from the closed channel produce the zero value with `ok == false`.

**Scenario:** A producer sends an entire batch into a buffered channel and closes it before the consumer starts reading. The consumer must still receive every buffered value.

---

### 11. Nil Channels

**Concept:** A channel's zero value is `nil`.

```go
var stream chan int
```

A nil channel behaves differently from a closed channel:

| Operation | Nil channel    |
| --------- | -------------- |
| Send      | Blocks forever |
| Receive   | Blocks forever |
| Close     | Panics         |

This makes nil channels particularly dangerous if they are used unintentionally.

**Scenario:** A program attempts to send, receive, and close a channel that was declared but never initialized. The exercise is to observe and understand the different behaviors.
