Based strictly on the attached content, there are three explicit exercise markers:

1. **Exercise: Concurrent map + 1 more**
2. **Exercise: Spot the race**
3. **Exercise: Rate limiter**

And, as requested, I’ve added **one extra exercise specifically for the final `Shared nothing` section**, even though that section has no explicit exercise marker. The extra question is designed around the concepts in that section only. 

---

# 1. Exercise — Concurrent Map

### Goal

Build a concurrent-safe map that supports concurrent reads and modifications while applying:

* shared mutable state
* data races
* `sync.Mutex`
* critical sections
* atomic multi-step operations
* `CompareAndSet`
* detecting stale state
* protecting a complex operation with a mutex

### Question

Build a concurrent-safe account store:

```go
type Accounts struct {
    // internal fields
}

func NewAccounts(balances map[string]int) *Accounts

func (a *Accounts) Get(name string) int

func (a *Accounts) Set(name string, amount int)

func (a *Accounts) CompareAndSet(
    name string,
    old int,
    new int,
) bool
```

The account store should initially contain:

```text
alice → 100
bob   → 50
```

Now simulate concurrent money transfers.

Implement:

```go
func transfer(
    accounts *Accounts,
    from string,
    to string,
    amount int,
) bool
```

A transfer must:

1. Read the sender's balance.
2. Check whether sufficient funds exist.
3. Deduct the amount from the sender.
4. Add the amount to the receiver.
5. Return `true` only if the transfer succeeds.

Run multiple transfers concurrently, for example:

```text
Alice → Bob   30
Alice → Bob   40
Alice → Bob   50
Bob   → Alice 20
```

The implementation must ensure that concurrent operations cannot produce an invalid account state.

First implement the transfer using a **mutex protecting the entire operation**.

Then implement the balance update using **compare-and-set** rather than protecting the entire transfer with one mutex.

Use the race detector:

```bash
go run -race .
```

Verify that:

* no data race occurs;
* money is never created or lost;
* an account never spends money based on a stale balance.

---

# 2. Exercise — CAS With Retry

### Goal

Practice the complete **compare-and-set retry pattern**:

* read current state
* calculate the desired new state
* attempt CAS
* detect failure because another goroutine changed the state
* retry
* stop when the operation succeeds
* stop when the operation is no longer possible
* avoid an uncontrolled infinite retry loop

### Question

You are building a concurrent inventory system.

```go
type Inventory struct {
    stock int
    mu    sync.Mutex
}
```

Implement:

```go
func (i *Inventory) Get() int

func (i *Inventory) CompareAndSet(old, new int) bool
```

Initially:

```text
stock = 10
```

Now start **20 goroutines**, where each goroutine tries to reserve one item:

```go
func reserve(inv *Inventory) bool {
    // implement
}
```

A reservation should:

1. Read the current stock.
2. If stock is `0`, immediately fail.
3. Calculate `stock - 1`.
4. Attempt `CompareAndSet(old, new)`.
5. If CAS fails, another goroutine changed the stock, so read it again and retry.
6. Stop after successfully reserving an item.
7. Limit the number of retries so a goroutine cannot retry forever.

Use a `sync.WaitGroup` to wait for all reservation attempts.

At the end, print:

```text
successful reservations: ?
remaining stock: ?
```

Verify that:

```text
successful reservations + remaining stock = 10
```

regardless of the order in which goroutines execute.

---

# 3. Exercise — Spot the Race

### Goal

Distinguish between a **data race** and a **race condition**, and understand why the race detector cannot necessarily identify an incorrect concurrent state.

Concepts to apply:

* data race
* race condition
* atomicity
* concurrent state transitions
* `sync.Mutex`
* `CompareAndSet`
* Go race detector

### Question

You are given this concurrent banking operation:

```go
type Account struct {
    balance int
    mu      sync.Mutex
}

func (a *Account) Get() int {
    a.mu.Lock()
    defer a.mu.Unlock()
    return a.balance
}

func (a *Account) Set(balance int) {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.balance = balance
}
```

The withdrawal operation is:

```go
func withdraw(a *Account, amount int) bool {
    balance := a.Get()

    if balance < amount {
        return false
    }

    time.Sleep(10 * time.Millisecond)

    a.Set(balance - amount)
    return true
}
```

Start two goroutines that simultaneously attempt:

```text
withdraw 70
withdraw 50
```

with an initial balance of:

```text
100
```

Your task is to investigate the program.

1. Run it repeatedly.
2. Run it with:

```bash
go run -race .
```

3. Determine whether the program contains a **data race**.
4. Determine whether it contains a **race condition**.
5. Explain why the race detector can report no data race even though the final balance can still be incorrect.
6. Fix the withdrawal operation so that the account can never approve withdrawals based on stale state.

Implement the fix in **two ways**:

* using a mutex around the complete check-and-update operation;
* using `CompareAndSet`.

The final implementation must never allow the account to spend more money than it contains.

---

# 4. Exercise — Rate Limiter

### Goal

Implement a concurrency limit where callers that arrive while the resource is occupied **fail immediately instead of blocking**.

Concepts to apply:

* `sync.Mutex`
* `TryLock`
* non-blocking lock acquisition
* immediate `"busy"` errors
* concurrent goroutines
* difference between `Lock()` and `TryLock()`
* why repeatedly calling `TryLock()` in a loop is undesirable

### Question

You are accessing a legacy external service that can process only **one request at a time**.

Create:

```go
type External struct {
    // internal fields
}

func (e *External) Call() error
```

Each successful call should simulate a 100 ms operation:

```go
time.Sleep(100 * time.Millisecond)
```

The requirements are:

1. Only one goroutine may call the external service at a time.
2. If the service is currently being used, another caller must **not wait**.
3. Instead, it should immediately receive:

```text
busy
```

4. A successful caller should receive `nil`.
5. The mutex must always be released after a successful lock.
6. Start 10 concurrent callers and wait for all of them.
7. Measure the total execution time.

Use `TryLock()` to implement the non-blocking behavior.

Then modify the program to implement the following **incorrect approach**:

```go
for {
    if mutex.TryLock() {
        // use resource
        mutex.Unlock()
        break
    }
}
```

Observe the behavior and explain why this creates a **busy-waiting loop** and can unnecessarily consume CPU.

Do **not** use this retry loop in the final implementation.

---

# 5. Extra Exercise — Shared-Nothing Purchase Processor

### Goal

Implement the **shared-nothing approach** from the final section by moving ownership of mutable account state into a single processor goroutine.

Concepts to apply:

* avoiding shared mutable state
* channels as communication boundaries
* single owner of state
* request/response channels
* concurrent callers
* sequential state modification inside one goroutine
* eliminating mutexes
* maintaining a single consistent version of account state

### Question

You are building a concurrent purchase system.

Define:

```go
type LegoSet struct {
    name  string
    price int
}

type Request struct {
    buyer string
    set   LegoSet
}

type Purchase struct {
    buyer   string
    set     LegoSet
    succeed bool
    balance int
}
```

Implement:

```go
func Processor(
    accounts map[string]int,
) (chan<- Request, <-chan Purchase)
```

The processor must:

1. Create an input channel for purchase requests.
2. Create an output channel for purchase results.
3. Make its own copy of the account state.
4. Run a dedicated goroutine that owns and modifies that account state.
5. Receive purchase requests from the input channel.
6. Check whether the buyer has enough money.
7. If sufficient funds exist:

   * deduct the price;
   * mark the purchase successful;
   * return the new balance.
8. Otherwise:

   * reject the purchase;
   * leave the balance unchanged.
9. Send the result through the output channel.
10. Ensure that **no mutex is required** for account access.

Then create multiple concurrent buyer goroutines:

```text
Alice → Castle   $20
Alice → Plants   $30
Alice → Robot    $40
Alice → Car      $10
```

with:

```text
Alice = $50
```

Each buyer goroutine should:

1. Send a `Request`.
2. Receive its corresponding `Purchase`.
3. Print whether the purchase succeeded and the resulting balance.

### Important requirement

The buyer goroutines must **never directly access or modify the account map**.

Only the processor goroutine may access the account state.

Finally, explain why this design avoids the race-condition problem from the earlier Alice/Lego example and why running **two processors over the same logical account state** would create a consistency problem. 
