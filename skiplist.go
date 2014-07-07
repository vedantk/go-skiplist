package main

import (
    "fmt"
    "bytes"
    "math/rand"
)

const kMaxHeight = 20
const kPrLevelBump = 0.5

type Key interface {}

type Value interface {}

type Comparator func(lhs, rhs Key) int

type Node struct {
    key Key
    value Value
    forward []*Node
}

type Skiplist struct {
    height int
    header Node
    cmp Comparator
}

func IntComparator(lhs, rhs Key) int {
    return lhs.(int) - rhs.(int)
}

func New() *Skiplist {
    sl := new(Skiplist)
    sl.header.forward = make([]*Node, kMaxHeight)
    sl.cmp = IntComparator
    return sl
}

func (sl *Skiplist) Search(search Key) Value {
    x := &sl.header
    for i := sl.height; i >= 0; i -= 1 {
        for x.forward[i] != nil && sl.cmp(x.forward[i].key, search) < 0 {
            x = x.forward[i]
        }
    }
    x = x.forward[0]
    if x != nil && sl.cmp(x.key, search) == 0 {
        return x.value
    } else {
        return nil
    }
}

func (sl *Skiplist) randomLevel() (lvl int) {
    for rand.Float32() < kPrLevelBump && lvl < (kMaxHeight - 1) {
        lvl += 1
        if lvl > sl.height {
            return lvl
        }
    }
    return
}

func (sl *Skiplist) Insert(search Key, data Value) {
    x := &sl.header
    preds := make([]*Node, sl.height + 2)
    for i := sl.height; i >= 0; i -= 1 {
        for x.forward[i] != nil && sl.cmp(x.forward[i].key, search) < 0 {
            x = x.forward[i]
        }
        preds[i] = x
    }
    x = x.forward[0]
    if x != nil && sl.cmp(x.key, search) == 0 {
        x.value = data
    } else {
        newLevel := sl.randomLevel()
        if newLevel > sl.height {
            preds[newLevel] = &sl.header
            sl.height = newLevel
        }

        x = new(Node)
        x.key = search
        x.value = data
        x.forward = make([]*Node, newLevel + 1)
        for i := 0; i <= newLevel; i += 1 {
            x.forward[i] = preds[i].forward[i]
            preds[i].forward[i] = x
        }
    }
}

func (sl *Skiplist) Delete(search Key) {
    x := &sl.header
    preds := make([]*Node, sl.height + 1)
    for i := sl.height; i >= 0; i -= 1 {
        for x.forward[i] != nil && sl.cmp(x.forward[i].key, search) < 0 {
            x = x.forward[i]
        }
        preds[i] = x
    }
    x = x.forward[0]
    if x != nil && sl.cmp(x.key, search) == 0 {
        for i := 0; i < sl.height; i += 1 {
            if preds[i].forward[i] != x {
                break
            }
            preds[i].forward[i] = x.forward[i]
        }
        if sl.header.forward[sl.height] == nil {
            sl.height -= 1
        }
    }
}

func (sl *Skiplist) String() string {
    var buf bytes.Buffer
    buf.WriteString(fmt.Sprintf("Skiplist(height: %d)\n", sl.height))
    for i := sl.height; i >= 0; i -= 1 {
        buf.WriteString(fmt.Sprintf("[%d]:", i))
        for x := sl.header.forward[i]; x != nil; x = x.forward[i] {
            buf.WriteString(fmt.Sprintf(" %d ", x.key))
        }
        buf.WriteString("\n")
    }
    return buf.String()
}

func (sl *Skiplist) checkInvariants() {
    // Check height.
    if sl.height >= kMaxHeight {
        panic(fmt.Sprintf("height: %d >= %d", sl.height, kMaxHeight))
    }
    if sl.height < 0 {
        panic(fmt.Sprintf("height: %d < 0", sl.height))
    }
    for i := 0; i < kMaxHeight; i += 1 {
        if i <= sl.height {
            if sl.header.forward[i] == nil {
                panic(fmt.Sprintf("level %d is sparse", i))
            }
        } else {
            if sl.header.forward[i] != nil {
                panic(fmt.Sprintf("level %d not sparse", i))
            }
        }
    }

    // Check ordering on each level.
    for i := 0; i <= sl.height; i += 1 {
        for x := sl.header.forward[i]; x != nil; x = x.forward[i] {
            if x.forward[i] != nil && sl.cmp(x.key, x.forward[i].key) >= 0 {
                panic(fmt.Sprintf("Level %d unsorted", i))
            }
        }
    }
}

func checkExpect(lhs, rhs Value) {
    if lhs != rhs {
        panic(fmt.Sprintf("checkExpect: %s != %s", lhs, rhs))
    }
}

func main() {
    sl := New()

    // Check basic inserts.
    hello := "Hello"
    world := "World"
    sl.Insert(0, hello)
    sl.Insert(1, world)
    checkExpect(sl.Search(0), hello)
    checkExpect(sl.Search(1), world)

    // Make sure searching for non-existent keys returns nil.
    for i := 2; i < 1000; i += 1 {
        if sl.Search(Key(i)) != nil {
            panic("Search failed to return <nil>")
        }
    }

    // Check random inserts/updates.
    expect := map[int]string{}
    for i := 0; i < 1000; i += 1 {
        k := rand.Int()
        kp := k % 100
        v := fmt.Sprintf("%d", k)
        sl.Insert(Key(kp), v)
        expect[kp] = v
    }
    for kp, v := range expect {
        checkExpect(sl.Search(Key(kp)), v)
    }

    sl.checkInvariants()

    fmt.Println(sl)

    // Try a lot of inserts.
    for i := 0; i < 1000000; i += 1 {
        k := Key(rand.Int())
        v := fmt.Sprintf("%d", k)
        sl.Insert(k, v)
        checkExpect(sl.Search(k), v)
    }

    sl.checkInvariants()

    fmt.Println("Final sl.height =", sl.height)
}
