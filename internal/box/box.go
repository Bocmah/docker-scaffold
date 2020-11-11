//go:generate go run generator.go

package box

type embedBox struct {
	storage map[string][]byte
}

func newEmbedBox() *embedBox {
	return &embedBox{storage: make(map[string][]byte)}
}

func (e *embedBox) Add(file string, content []byte) {
	e.storage[file] = content
}

func (e *embedBox) Get(file string) []byte {
	if f, ok := e.storage[file]; ok {
		return f
	}

	return nil
}

var box = newEmbedBox()

func Add(file string, content []byte) {
	box.Add(file, content)
}

func Get(file string) []byte {
	return box.Get(file)
}
