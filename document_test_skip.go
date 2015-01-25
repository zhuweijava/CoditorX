package main

import (
	"fmt"
	"testing"
)

func Test_ringbuff_test(t *testing.T) {
	rc, err := newRingCache(10)
	if err != nil {
		t.Error(err)
		return
	}

	for i := 0; i < 11; i++ {
		rc.Put(i)
	}

	vals, err := rc.Tail(5)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(vals)
}

func Test_document_test(t *testing.T) {
	doc, err := newDocument("Hi,Tom.", 0)
	if err != nil {
		t.Error(err)
		return
	}

	content := doc.getContents()

	content = "Hi,Peter."
	patchs, _, err := doc.merge(content, doc.getVersion(""))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(patchs)

	content += "change1!"
	patchs, _, err = doc.merge(content, doc.getVersion(""), "")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(patchs)

	content = "change2!" + content
	patchs, _, err = doc.merge(content, doc.getVersion(""), "")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(patchs)

	var patchStrs []string
	patchStrs, err = doc.tail(0, "")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(doc.getContents(""))
	fmt.Println(patchStrs)
}
