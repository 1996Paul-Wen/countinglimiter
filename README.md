# countinglimiter
This package provides a simple and lightweight implementation of a count-based rate limiter

# usage
This package provides a `Limiter` struct and only 5 associated methods: `NewLimiter`, `Start`, `Allow`, `AllowN` and `Stop`

- `NewLimiter` generates a Limiter with specified ratelimit and time interval. N(N = ratelimit) requests at most could go through in every time interval.

- `Start` runs the Limiter. If a Limiter has not Started, it can head off nothing. In other words, anything can go through the Limiter. **So start the Limiter before using it**. And it's not harmful to start a Limiter that is already started, for it does nothing and return. 

- `Allow` is the shorthand of `AllowN`.

- `AllowN` judges if n requests could go through.

- `Stop` makes Limiter stop and resets the Limiter. And it's not harmful to stop a Limiter that is already stoped, for it does nothing and return.

here is an example for simple usage:
```
// get a limiter l 
l := NewLimiter(1000, 1*time.Second)

// start l
l.Start()

// use l to head off requests out of limit
headOff := 0
for i := 0; i < 2000; i++ {
    if !l.Allow() {
        headOff += 1
    }
}
fmt.Printf("headOff: %d\n", headOff)  // around 1000

// stop l
l.Stop()
```