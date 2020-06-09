package main
import "fmt"
import "golang.org/x/net/html"
import "os"
import "strings"
import "strconv"
import "bufio"
import "io"
import "bytes"
import "log"
/*
Christian Lockley 2020
This is a simple HTML transcluder, to use simply add a directive to a HTML file and this program will print the contents of the
file spesifed in the directive and the file included as a command line argument to STDOUT. Not this program does not support
nested includes.
EXAMPLE:
<!--#include file="footer.html" -->
*/
func readFileWithReadLine(fn string) (err error) {
	//From https://stackoverflow.com/questions/8757389/reading-a-file-line-by-line-in-go/16615559#16615559
    file, err := os.Open(fn)
    defer file.Close()

    if err != nil {
        return err
    }
    // Start reading from the file with a reader.
    reader := bufio.NewReader(file)
    for {
        var buffer bytes.Buffer
        var l []byte
        var isPrefix bool
        for {
            l, isPrefix, err = reader.ReadLine()
            buffer.Write(l)

            // If we've reached the end of the line, stop reading.
            if !isPrefix {
                break
            }
            // If we're just at the EOF, break
            if err != nil {
                break
            }
        }
        if err == io.EOF {
            break
        }
		line := buffer.String()
		print(line)
    }
    if err != io.EOF {
        fmt.Printf(" > Failed!: %v\n", err)
    }
    return nil
}
func usage() {
	print("USAGE: precludehtml infile")
	os.Exit(1)
}
func main() {
	if len(os.Args) != 2 {
		usage()
	}
	file, err := os.Open(os.Args[1])
	defer file.Close()
	if (err != nil) {
		os.Exit(1)
	}
	parser := html.NewTokenizer(file)
	fmt.Println(os.Args[1])
	for {
		tmp := parser.Next()
		if tmp == html.ErrorToken {
			break
		}
		if tmp == html.CommentToken {
			inputSplitBySpace := strings.Fields(strings.ToLower(string(parser.Text())))
			if len(inputSplitBySpace) <= 1 {
				print(string(parser.Raw()))
				continue
			}
			if inputSplitBySpace[0] != "#include" {
				print(string(parser.Raw()))
			}
			if (strings.Split(inputSplitBySpace[1],"=")[0] != "file") {
				print(string(parser.Raw()))
				continue
			}
			s, _ := strconv.Unquote(strings.ReplaceAll(strings.Split(inputSplitBySpace[1],"=")[1],"'", "\"")) ;
			err := readFileWithReadLine(s)
			if (err != nil) {
				log.Fatal(err)
			}
			continue;			
		}
		print(string(parser.Raw()))
	}
}	