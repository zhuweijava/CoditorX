package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func Test_create(t *testing.T) {
	_, err := os.Stat("test")

	exists := err == nil || os.IsExist(err)

	if !exists {
		os.Mkdir("test", 0755)
	}
}

func Test_write(t *testing.T) {
	data := make([]byte, 100)
	for i := 0; i < len(data); i++ {
		data[i] = 1
	}

	begin := time.Now().Unix()
	log, err := openBinLog("test/binLog.log")
	if err == nil {
		for i := 0; i < 100000; i++ {
			log.append(uint32(i), data)
		}
	}
	log.close()
	end := time.Now().Unix()
	t.Logf("write 100000 use second: %d.\n", end-begin)
}

func Test_read(t *testing.T) {
	blr, err := openBinLogReader("test/binLog.log")
	defer func() {
		blr.close()
	}()
	if err == nil {
		for i := 0; i < 10; i++ {
			version, data, err := blr.next()
			if err != nil {
				t.Logf("err is %v.\n", err)
				break
			} else {
				t.Log(version)
				t.Logf("length: %d, %v\n", len(data), data)
			}
		}
	}
}

func Test_walk(t *testing.T) {
	blr, err := openBinLogReader("test/binLog.log")
	defer func() {
		blr.close()
	}()
	begin := time.Now().Unix()
	if err == nil {
		for {
			_, _, err := blr.next()
			if err != nil {
				fmt.Println(err)
				break
			}
		}
	} else {
		fmt.Println(err)
	}
	end := time.Now().Unix()
	t.Logf("walk use second: %d.\n", end-begin)
}

func Test_del(t *testing.T) {
	err := os.Remove("test/binLog.log")
	if err != nil {
		t.Log(err)
	}
}
