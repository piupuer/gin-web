package utils

import (
	"fmt"
	"testing"
)

func TestZip(t *testing.T) {
	fmt.Println(Zip("/Users/piupuer/eric/images", "/Users/piupuer/eric/2/images.zip"))
	fmt.Println(Zip("/Users/piupuer/eric/images", "/Users/piupuer/eric/2/张三.zip"))
}

func TestUnZip(t *testing.T) {
	fmt.Println(UnZip("/Users/piupuer/eric/2/images.zip", "/Users/piupuer/eric/3/images"))
	fmt.Println(UnZip("/Users/piupuer/eric/2/张三.zip", "/Users/piupuer/eric/4/images"))
}
