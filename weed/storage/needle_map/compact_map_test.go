package needle_map

import (
	"fmt"
	. "github.com/chrislusf/seaweedfs/weed/storage/types"
	"testing"
)

func TestOverflow2(t *testing.T) {
	m := NewCompactMap()
	m.Set(NeedleId(150088), ToOffset(8), 3000073)
	m.Set(NeedleId(150073), ToOffset(8), 3000073)
	m.Set(NeedleId(150089), ToOffset(8), 3000073)
	m.Set(NeedleId(150076), ToOffset(8), 3000073)
	m.Set(NeedleId(150124), ToOffset(8), 3000073)
	m.Set(NeedleId(150137), ToOffset(8), 3000073)
	m.Set(NeedleId(150147), ToOffset(8), 3000073)
	m.Set(NeedleId(150145), ToOffset(8), 3000073)
	m.Set(NeedleId(150158), ToOffset(8), 3000073)
	m.Set(NeedleId(150162), ToOffset(8), 3000073)

	m.AscendingVisit(func(value NeedleValue) error {
		println("needle key:", value.Key)
		return nil
	})
}

func TestIssue52(t *testing.T) {
	m := NewCompactMap()
	m.Set(NeedleId(10002), ToOffset(10002), 10002)
	if element, ok := m.Get(NeedleId(10002)); ok {
		fmt.Printf("key %d ok %v %d, %v, %d\n", 10002, ok, element.Key, element.Offset, element.Size)
	}
	m.Set(NeedleId(10001), ToOffset(10001), 10001)
	if element, ok := m.Get(NeedleId(10002)); ok {
		fmt.Printf("key %d ok %v %d, %v, %d\n", 10002, ok, element.Key, element.Offset, element.Size)
	} else {
		t.Fatal("key 10002 missing after setting 10001")
	}
}

func TestCompactMap(t *testing.T) {
	m := NewCompactMap()
	for i := uint32(0); i < 100*batch; i += 2 {
		m.Set(NeedleId(i), ToOffset(int64(i)), i)
	}

	for i := uint32(0); i < 100*batch; i += 37 {
		m.Delete(NeedleId(i))
	}

	for i := uint32(0); i < 10*batch; i += 3 {
		m.Set(NeedleId(i), ToOffset(int64(i+11)), i+5)
	}

	//	for i := uint32(0); i < 100; i++ {
	//		if v := m.Get(Key(i)); v != nil {
	//			glog.V(4).Infoln(i, "=", v.Key, v.Offset, v.Size)
	//		}
	//	}

	for i := uint32(0); i < 10*batch; i++ {
		v, ok := m.Get(NeedleId(i))
		if i%3 == 0 {
			if !ok {
				t.Fatal("key", i, "missing!")
			}
			if v.Size != i+5 {
				t.Fatal("key", i, "size", v.Size)
			}
		} else if i%37 == 0 {
			if ok && v.Size != TombstoneFileSize {
				t.Fatal("key", i, "should have been deleted needle value", v)
			}
		} else if i%2 == 0 {
			if v.Size != i {
				t.Fatal("key", i, "size", v.Size)
			}
		}
	}

	for i := uint32(10 * batch); i < 100*batch; i++ {
		v, ok := m.Get(NeedleId(i))
		if i%37 == 0 {
			if ok && v.Size != TombstoneFileSize {
				t.Fatal("key", i, "should have been deleted needle value", v)
			}
		} else if i%2 == 0 {
			if v == nil {
				t.Fatal("key", i, "missing")
			}
			if v.Size != i {
				t.Fatal("key", i, "size", v.Size)
			}
		}
	}

}

func TestOverflow(t *testing.T) {
	cs := NewCompactSection(1)

	cs.setOverflowEntry(1, ToOffset(12), 12)
	cs.setOverflowEntry(2, ToOffset(12), 12)
	cs.setOverflowEntry(3, ToOffset(12), 12)
	cs.setOverflowEntry(4, ToOffset(12), 12)
	cs.setOverflowEntry(5, ToOffset(12), 12)

	if cs.overflow[2].Key != 3 {
		t.Fatalf("expecting o[2] has key 3: %+v", cs.overflow[2].Key)
	}

	cs.setOverflowEntry(3, ToOffset(24), 24)

	if cs.overflow[2].Key != 3 {
		t.Fatalf("expecting o[2] has key 3: %+v", cs.overflow[2].Key)
	}

	if cs.overflow[2].Size != 24 {
		t.Fatalf("expecting o[2] has size 24: %+v", cs.overflow[2].Size)
	}

	cs.deleteOverflowEntry(4)

	if len(cs.overflow) != 4 {
		t.Fatalf("expecting 4 entries now: %+v", cs.overflow)
	}

	_, x, _ := cs.findOverflowEntry(5)
	if x.Key != 5 {
		t.Fatalf("expecting entry 5 now: %+v", x)
	}

	for i, x := range cs.overflow {
		println("overflow[", i, "]:", x.Key)
	}
	println()

	cs.deleteOverflowEntry(1)

	for i, x := range cs.overflow {
		println("overflow[", i, "]:", x.Key)
	}
	println()

	cs.setOverflowEntry(4, ToOffset(44), 44)
	for i, x := range cs.overflow {
		println("overflow[", i, "]:", x.Key)
	}
	println()

	cs.setOverflowEntry(1, ToOffset(11), 11)

	for i, x := range cs.overflow {
		println("overflow[", i, "]:", x.Key)
	}
	println()

}
