package fastrlp

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestArena(t *testing.T) {
	t.Run("serial", func(t *testing.T) {
		var a Arena
		for i := 0; i < 10; i++ {
			if err := testArena(&a); err != nil {
				t.Fatal(err)
			}
			a.Reset()
		}
	})

	t.Run("concurrent", func(t *testing.T) {
		var ap ArenaPool
		workers := 4
		ch := make(chan error, workers)
		for i := 0; i < workers; i++ {
			go func() {
				a := ap.Get()
				defer ap.Put(a)
				var err error
				for i := 0; i < 10; i++ {
					if err = testArena(a); err != nil {
						break
					}
				}
				ch <- err
			}()
		}
		for i := 0; i < workers; i++ {
			select {
			case err := <-ch:
				if err != nil {
					t.Fatal(err)
				}
			case <-time.After(time.Second):
				t.Fatalf("timeout")
			}
		}
	})
}

func testArena(a *Arena) error {
	var v *Value

	// empty string
	v = a.NewString("")
	if err := validate(v, "0x80"); err != nil {
		return err
	}

	// bytestring00
	v = a.NewBytes([]byte{0x0})
	if err := validate(v, "0x00"); err != nil {
		return err
	}

	// bytestring01
	v = a.NewBytes([]byte{0x1})
	if err := validate(v, "0x01"); err != nil {
		return err
	}

	// bytestring7F
	v = a.NewBytes([]byte{0x7F})
	if err := validate(v, "0x7F"); err != nil {
		return err
	}

	// short string
	v = a.NewString("dog")
	if err := validate(v, "0x83646f67"); err != nil {
		return err
	}

	// short string2
	v = a.NewString("Lorem ipsum dolor sit amet, consectetur adipisicing eli")
	if err := validate(v, "0xb74c6f72656d20697073756d20646f6c6f722073697420616d65742c20636f6e7365637465747572206164697069736963696e6720656c69"); err != nil {
		return err
	}

	// long string
	v = a.NewString("Lorem ipsum dolor sit amet, consectetur adipisicing elit")
	if err := validate(v, "0xb8384c6f72656d20697073756d20646f6c6f722073697420616d65742c20636f6e7365637465747572206164697069736963696e6720656c6974"); err != nil {
		return err
	}

	// long string 2
	v = a.NewString("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur mauris magna, suscipit sed vehicula non, iaculis faucibus tortor. Proin suscipit ultricies malesuada. Duis tortor elit, dictum quis tristique eu, ultrices at risus. Morbi a est imperdiet mi ullamcorper aliquet suscipit nec lorem. Aenean quis leo mollis, vulputate elit varius, consequat enim. Nulla ultrices turpis justo, et posuere urna consectetur nec. Proin non convallis metus. Donec tempor ipsum in mauris congue sollicitudin. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia Curae; Suspendisse convallis sem vel massa faucibus, eget lacinia lacus tempor. Nulla quis ultricies purus. Proin auctor rhoncus nibh condimentum mollis. Aliquam consequat enim at metus luctus, a eleifend purus egestas. Curabitur at nibh metus. Nam bibendum, neque at auctor tristique, lorem libero aliquet arcu, non interdum tellus lectus sit amet eros. Cras rhoncus, metus ac ornare cursus, dolor justo ultrices metus, at ullamcorper volutpat")
	if err := validate(v, "0xb904004c6f72656d20697073756d20646f6c6f722073697420616d65742c20636f6e73656374657475722061646970697363696e6720656c69742e20437572616269747572206d6175726973206d61676e612c20737573636970697420736564207665686963756c61206e6f6e2c20696163756c697320666175636962757320746f72746f722e2050726f696e20737573636970697420756c74726963696573206d616c6573756164612e204475697320746f72746f7220656c69742c2064696374756d2071756973207472697374697175652065752c20756c7472696365732061742072697375732e204d6f72626920612065737420696d70657264696574206d6920756c6c616d636f7270657220616c6971756574207375736369706974206e6563206c6f72656d2e2041656e65616e2071756973206c656f206d6f6c6c69732c2076756c70757461746520656c6974207661726975732c20636f6e73657175617420656e696d2e204e756c6c6120756c74726963657320747572706973206a7573746f2c20657420706f73756572652075726e6120636f6e7365637465747572206e65632e2050726f696e206e6f6e20636f6e76616c6c6973206d657475732e20446f6e65632074656d706f7220697073756d20696e206d617572697320636f6e67756520736f6c6c696369747564696e2e20566573746962756c756d20616e746520697073756d207072696d697320696e206661756369627573206f726369206c756374757320657420756c74726963657320706f737565726520637562696c69612043757261653b2053757370656e646973736520636f6e76616c6c69732073656d2076656c206d617373612066617563696275732c2065676574206c6163696e6961206c616375732074656d706f722e204e756c6c61207175697320756c747269636965732070757275732e2050726f696e20617563746f722072686f6e637573206e69626820636f6e64696d656e74756d206d6f6c6c69732e20416c697175616d20636f6e73657175617420656e696d206174206d65747573206c75637475732c206120656c656966656e6420707572757320656765737461732e20437572616269747572206174206e696268206d657475732e204e616d20626962656e64756d2c206e6571756520617420617563746f72207472697374697175652c206c6f72656d206c696265726f20616c697175657420617263752c206e6f6e20696e74657264756d2074656c6c7573206c65637475732073697420616d65742065726f732e20437261732072686f6e6375732c206d65747573206163206f726e617265206375727375732c20646f6c6f72206a7573746f20756c747269636573206d657475732c20617420756c6c616d636f7270657220766f6c7574706174"); err != nil {
		return err
	}

	// zero
	v = a.NewUint(0)
	if err := validate(v, "0x80"); err != nil {
		return err
	}

	// smallint
	v = a.NewUint(1)
	if err := validate(v, "0x01"); err != nil {
		return err
	}

	// smallint2
	v = a.NewUint(16)
	if err := validate(v, "0x10"); err != nil {
		return err
	}

	// smallint3
	v = a.NewUint(79)
	if err := validate(v, "0x4f"); err != nil {
		return err
	}

	// smallint4
	v = a.NewUint(127)
	if err := validate(v, "0x7f"); err != nil {
		return err
	}

	// medium int

	// mediumint1
	v = a.NewUint(128)
	if err := validate(v, "0x8180"); err != nil {
		return err
	}

	// mediumint2
	v = a.NewUint(1000)
	if err := validate(v, "0x8203e8"); err != nil {
		return err
	}

	// mediumint3
	v = a.NewUint(100000)
	if err := validate(v, "0x830186a0"); err != nil {
		return err
	}

	// emptylist
	v = a.NewArray()
	if err := validate(v, "0xc0"); err != nil {
		return err
	}

	// stringlist
	v = a.NewArray()
	v.Set(a.NewString("dog"))
	v.Set(a.NewString("god"))
	v.Set(a.NewString("cat"))
	if err := validate(v, "0xcc83646f6783676f6483636174"); err != nil {
		return err
	}

	// multilist
	v = a.NewArray()
	v.Set(a.NewString("zw"))
	vv := a.NewArray()
	vv.Set(a.NewUint(4))
	v.Set(vv)
	v.Set(a.NewUint(1))
	if err := validate(v, "0xc6827a77c10401"); err != nil {
		return err
	}

	// TODO, build the rest of the tests
	return nil
}

func validate(v *Value, expected string) error {
	if strings.HasPrefix(expected, "0x") {
		expected = expected[2:]
	}
	buf, err := hex.DecodeString(expected)
	if err != nil {
		return err
	}
	dst := v.MarshalTo(nil)
	if !bytes.Equal(dst, buf) {
		return fmt.Errorf("bad")
	}
	return nil
}
