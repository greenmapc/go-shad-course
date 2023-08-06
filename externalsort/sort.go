//go:build !solution

package externalsort

import (
	"container/heap"
	"io"
	"os"
	"sort"
	"strings"
)

const end = '\n'

type CustomLineReader struct {
	Reader io.Reader
	buffer []byte
}

type CustomLineWriter struct {
	Writer io.Writer
}

func NewReader(r io.Reader) *CustomLineReader {
	return &CustomLineReader{
		Reader: r,
		buffer: make([]byte, 0),
	}
}

func NewWriter(w io.Writer) CustomLineWriter {
	return CustomLineWriter{
		Writer: w,
	}
}

func (w CustomLineWriter) Write(l string) error {
	data := []byte(l + "\n")
	_, error := w.Writer.Write(data)
	return error
}

func (r *CustomLineReader) ReadLine() (string, error) {
	for {
		if i := strings.IndexByte(string(r.buffer), end); i != -1 {
			line := r.buffer[:i]
			r.buffer = r.buffer[i+1:]
			return string(line), nil
		}

		tmp := make([]byte, 1024) // choose the best buffer size

		bytes, error := r.Reader.Read(tmp)

		if error != nil {
			if error == io.EOF {
				line := append(r.buffer, tmp[:bytes]...)
				r.buffer = r.buffer[:0]
				return string(line), io.EOF
			}
			return "", error
		}

		r.buffer = append(r.buffer, tmp[:bytes]...)
	}
}

type Item struct {
	reader LineReader
	line   string
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].line < pq[j].line
}

func (pq PriorityQueue) Swap(i, j int) {
	tmp := pq[i]
	pq[i] = pq[j]
	pq[j] = tmp
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

func Merge(w LineWriter, readers ...LineReader) error {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	// init the pq
	for _, reader := range readers {
		line, err := reader.ReadLine()

		if err == io.EOF && line == "" {
			continue
		}
		if err != nil && err != io.EOF {
			return err
		}
		heap.Push(&pq, &Item{
			reader: reader,
			line:   line,
		})
	}

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		err := w.Write(item.line)
		if err != nil {
			return err
		}

		newOne, err := item.reader.ReadLine()

		if err != nil && err != io.EOF {
			return err
		}

		if err == io.EOF && newOne == "" {
			continue
		}

		heap.Push(&pq, &Item{
			reader: item.reader,
			line:   newOne,
		})
	}

	return nil
}

func Sort(w io.Writer, in ...string) error { // in is a file name.
	readers := make([]LineReader, 0)

	for _, fileName := range in {
		file, err := os.Open(fileName)

		if err != nil {
			return err
		}

		defer file.Close()

		reader := NewReader(file)

		fileLines := make([]string, 0)
		for {
			line, readErr := reader.ReadLine()

			if readErr == io.EOF && line == "" {
				break
			}

			if readErr != nil && readErr != io.EOF {
				return readErr
			}

			fileLines = append(fileLines, line)
		}
		sort.Strings(fileLines)

		file, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)

		if err != nil {
			return err
		}

		lw := NewWriter(file)

		for _, line := range fileLines {
			err = lw.Write(line)
			if err != nil {
				return err
			}
		}
		file.Close()
	}

	for _, fileName := range in {
		file, err := os.Open(fileName)

		if err != nil {
			return err
		}
		readers = append(readers, NewReader(file))
	}

	lineWriter := NewWriter(w)
	err := Merge(lineWriter, readers...)

	return err
}
