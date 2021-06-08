/*
Package gigasecond implements a simple library for a gigasecond to a time.
*/
package gigasecond

// import path for the time package from the standard library
import (
"time"
)

const GIGASECOND = time.Second * 1000000000

// AddGigasecond return a given time with one gigasecond added to it
func AddGigasecond(t time.Time) time.Time {
    return t.Add(GIGASECOND)
}
