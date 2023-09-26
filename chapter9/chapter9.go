package chapter9 // Package chapter9 包声明，当前包被其他包引入时作为默认的标识符（包名）

import (
	"fmt"
	"image/jpeg"
	"os"

	/*
		1. 包名是导入路径的最后一段（包名不包含尾缀），即使两个相同的包名，但是路径不同，也是允许的，但是需要给其中一个定义别名（重命名导入）
		2. 如果导入的包在文件中没有被引用（为了执行引入包的init函数），会产生编译错误，使用`_`表示导入的内容为空白标识符
	*/
	_ "github.com/go-sql-driver/mysql"
	"image"
	"io"
)

// helloworld
func Main() {
	if err := toJPEG(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "jpeg: %v\n", err)
		os.Exit(1)
	}
}

func toJPEG(in io.Reader, out io.Writer) error {
	img, kind, err := image.Decode(in)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "Input format =", kind)
	return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
}
