package templates

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

var _doc_templates_doc_md = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x56\xfb\x6f\x1a\xb9\x13\xff\x7d\xff\x8a\xd1\xa6\xfa\x2a\xd1\x37\x31\x90\x34\x7d\x20\xf5\x87\xb4\xa1\x39\xaa\x04\x38\xa0\xbd\x47\x55\x11\xb3\x3b\x80\x9b\xc5\xde\xda\x5e\x52\x6e\xd9\xff\xfd\x34\xf6\xf2\x4c\x9b\xea\xda\x54\xf7\xcb\x49\x48\xac\x67\xc6\xf3\xf8\xcc\xcb\x7b\xf0\x46\x21\xfc\x9a\xa1\x16\x68\x82\xa0\x3f\x11\x06\x62\x15\x65\x53\x94\x16\x26\xdc\xc0\x10\x51\x02\xcf\xac\x9a\x72\x2b\x22\x9e\x24\x73\x18\xa3\x44\xcd\x2d\xc6\x70\x2b\xec\x04\xde\x7f\x54\x08\xb3\x1a\xab\xb2\xea\x87\xfd\x89\xb5\xa9\xa9\x57\x2a\x63\x61\x27\xd9\x90\x45\x6a\x5a\xc1\x99\x48\x8c\x8a\x6e\xd0\x56\x3e\x2a\x3c\x60\x41\xb0\xb7\x07\x57\x5c\x48\x32\x87\x06\x81\x6b\x04\x3b\x41\xb2\x32\x41\x49\x56\xac\x50\x12\xb8\x8c\x61\x26\x4c\xc6\x13\xf1\x97\xa7\x38\xc3\x22\x02\x94\x71\xaa\x84\xb4\xc6\xe9\xda\x83\x8b\x46\x7f\xd1\x69\xf7\xfa\x50\xe1\xa9\xa8\xcc\x6a\x15\xd2\x04\x10\x04\x67\x6b\x8d\x08\x56\x39\x2b\x67\x9d\xa6\xf7\x3b\x33\xa8\x25\x9f\xa2\x33\x94\x72\x63\x6e\x95\x8e\x41\x48\x50\x3a\x46\x4d\xe2\x63\xb4\xc0\xe1\xcd\x6f\x7d\xb0\xea\x06\x25\x0b\x82\x45\x87\x6b\x3e\x45\x8b\x7a\x71\x8e\x23\x9e\x25\x76\x11\x2c\x8e\x8e\xe8\x17\x2c\xe0\x9a\x54\x5e\xc3\x02\x06\x52\x49\x1c\x80\xa3\x91\xe6\x2d\x1a\xb9\xbc\x07\x5d\xfc\x94\xa1\xb1\x41\x00\x00\x10\x65\x3a\x81\xa3\xa3\x98\x5b\x0e\x21\x29\x79\xc1\xe3\xa9\x90\xff\xa3\xbb\xfe\x33\x04\x42\xb6\x5e\xa9\x24\x2a\xe2\xc9\x44\x19\x5b\x7f\x56\x7d\x56\xdd\x0c\x78\xa5\xd8\xa4\x4a\x1a\x0c\x82\xeb\xeb\xeb\x8f\x46\xc9\x20\x0f\x00\x42\x17\x41\x58\x87\x90\x31\x16\x06\x05\x31\x57\xe0\xad\x70\xfb\xe4\xeb\xa0\x12\x04\x17\x2e\xf4\x44\x18\x0b\x6a\xe4\x70\x2b\x79\x60\x27\xdc\xba\x94\x45\x99\xd6\x28\x6d\x32\x87\x44\xf1\x18\x1d\x74\x24\x68\xe6\xc6\xe2\x94\x7d\x3d\xce\x5f\x20\xa4\xc4\x28\x5d\xe6\xb5\x0e\x2f\x91\x6b\xd4\xc0\x98\xc7\x99\x7d\x23\xda\x4f\xcb\x72\xfd\x4a\xc0\xef\x03\x80\xdc\x19\x0c\x23\x8d\x54\xa9\x03\x6e\x29\xf6\xe3\x6a\xed\xf9\x51\xad\x76\x54\x7d\xde\xaf\x55\xeb\xa7\x4f\xeb\xd5\x1a\x7b\xfa\xec\xf1\xc9\x49\xed\xf8\xf4\xf4\xff\xd5\x5a\xbd\x5a\x0d\x0f\xfd\xc5\x2c\x8d\xbf\x7d\xf1\xf1\xf3\xc7\xc7\x27\x3b\x17\xa9\xa4\xe8\x8a\x55\xe9\x20\x4d\xf8\x1c\xb5\x59\xb2\x22\x1e\x4d\x88\xe7\x7d\xa3\xa4\xcc\x53\x3a\xd7\x0e\x97\x84\x1b\x9c\x9b\xb0\x0e\xef\xcb\x33\x40\x98\x88\xa9\xb0\x61\x79\xfe\xb0\x12\xb4\x36\x09\xeb\x70\xf2\xa4\x5a\x75\x94\xa2\xb4\x10\xa3\x89\xb4\x48\x09\x55\xf2\xa1\xaf\x52\x28\x7d\x80\xe1\x1c\x78\x14\xa1\x31\x50\xf6\x4e\xb8\xba\xe3\x0a\xd9\x6c\x3a\xe6\xad\xd6\xe1\xf8\x74\x4b\x3d\x01\x3f\x27\xc5\xbd\xc6\x65\xe3\x55\x1f\x20\x63\x6b\x9c\x80\x1b\xe0\x91\x15\x33\x1c\x70\x7b\x08\x19\x23\x24\xe8\x7f\x24\xe4\x18\x75\xaa\x85\x74\xe4\x48\x65\xd2\xea\xf9\x21\xbc\x6a\xbf\x6d\xf5\xf7\x39\x13\xf1\x01\x9c\xf5\x40\xa2\xbd\x55\xfa\xc6\xc0\xeb\x6e\xfb\x0a\x32\x29\xac\x81\x0c\x9a\xad\x56\xa3\x0b\x6f\xda\xcd\x56\xe9\xfe\xc0\xbb\x0f\x1c\xda\x2d\xc8\x98\x88\xe1\x05\x70\x46\xe2\x03\x11\xc3\x45\xb7\xfd\xb6\x03\x2f\xff\xf0\x9c\x76\xf7\xbc\xd1\xa5\xd3\x4a\xf7\x79\xa3\xf7\x0a\x2e\x9b\x57\xcd\x3e\xe4\x2e\xc8\x62\x09\xc3\x4c\xe0\xed\x16\x06\x43\xae\xcd\x4e\x22\x07\x44\x63\x63\x15\x6e\xa1\xe2\xfd\xda\x48\x5b\xe8\x5b\xd6\x9d\x3e\x04\xa5\x20\x63\xac\x6c\xbb\x3c\xd7\x5c\x8e\x11\x58\x39\x76\x8b\xe2\xcb\x9d\x38\xaf\xe4\x39\x6b\xf1\x29\x16\x45\x85\xbc\xa3\x89\xd6\x9b\xa8\x5b\x10\x72\xa4\xf4\xb4\x9c\x92\x43\x95\x59\xd7\x77\x2b\x61\xd7\xaa\xf3\x3b\x0d\x98\xe7\x62\x04\x8c\x3a\x8f\x48\x42\x63\x0c\x45\xf1\xd0\x6d\xb9\xeb\x73\x90\xe7\x98\x18\xdc\x34\xf4\x5d\x3a\x64\x5c\xa2\xf4\xc5\x01\x97\xd3\x1f\xb0\x62\x7b\xac\x6d\xef\x84\x1d\xcd\xfb\x8c\xee\x2c\x58\x64\x66\x07\x84\x0d\x3b\x5f\xb7\x0e\x99\xf2\x68\xad\xc6\xbd\x21\xac\xee\x9b\xfe\xcb\xa4\x3e\xf2\x55\xff\x28\xc6\x11\xd4\x5f\xec\x68\x58\xc0\x75\x9e\x3b\x89\xa2\xa0\x9d\xe0\x8c\x90\x68\x51\xe4\xf9\xf2\xdf\xe3\xe5\xd7\x45\x19\x39\x2c\x56\x18\xec\x62\xf1\x2f\x25\xd7\xa1\x97\xe7\xae\x84\xe7\x3d\xab\x85\x1c\x7b\xe7\x7e\x20\xd9\xac\x4c\xe6\xc3\x27\xbb\x82\x9f\xd3\xc4\xbd\x38\xba\x68\x33\x2d\x41\xa3\xa1\x91\x07\x23\xa5\x81\x4b\x68\xfc\xde\xb9\x3c\x6b\xb6\x40\xa5\xf4\xb0\xa1\xb6\x52\x72\xa7\xa7\xa6\x5c\xc8\x55\x63\xfd\x57\x1c\xf7\x77\x6d\x09\xf7\x03\xd7\xc7\x7d\x6a\xef\x2f\x99\x60\xb9\x7c\xcb\x7d\x2e\xb3\x24\xa1\x99\x1c\xe2\x67\x8c\x06\x56\xb8\x8d\x7d\xfc\xf4\xc9\xc9\xf1\x93\x13\x47\x97\xd9\x74\xa0\x31\x52\x3a\x36\xcb\xfd\x1c\xae\xcf\x7e\xce\xaf\x16\x45\xe3\xb3\xd5\x9c\x36\xc5\x9f\xa8\x15\xf8\xd5\xb9\x5a\xd4\x23\x91\x58\xd4\x18\x6f\x58\x75\x74\x11\xef\x2e\xfe\x5d\x89\x1b\x9c\x0f\x12\xf7\x66\xdb\x22\xa7\x5c\x5b\x41\x59\x32\x77\x38\xca\x18\x31\x4c\x70\x50\xbe\x22\xb6\x98\x1a\x47\x77\x48\xea\xf6\x8e\x98\xc1\x04\x23\x3b\x28\x1f\x26\x61\xaf\x79\xd5\xb9\x6c\xac\xa3\xb1\x7c\x98\xe0\xee\x9d\x52\x98\x68\x7e\x2f\x06\xb4\xf6\x7c\x7b\x42\x90\xe7\xf0\xc8\x25\x92\x92\xe8\xaa\xde\x7d\xb8\xb4\x95\x1c\xd7\x06\xe6\x6e\x47\xac\x25\x7c\xaa\x9d\xc4\x46\xea\xb7\x44\xb6\xca\x9b\x04\x77\xeb\x7d\xdd\x85\xb4\x54\x5a\xbe\x13\xe9\xf3\xb5\x48\xbc\x63\xef\xe8\x09\x40\xa2\x5f\xdb\xc6\xeb\x40\x8a\x82\x8e\x4b\x45\x45\xc1\xf6\x53\x39\x5e\x98\xd9\xf8\x60\x35\x63\x38\x74\x5a\x17\xa0\x34\xf4\xde\x5d\x80\xc6\x54\xa3\x41\x69\xcb\xf1\x32\x02\x0e\x5b\x0a\x20\x9a\x70\x6d\xdd\x40\xf2\x83\x67\xd3\xd6\xf6\xe0\xd9\x02\xed\xfb\x26\xcf\x8e\x8a\x9f\x3c\x7a\xbe\x90\xa0\x9f\x33\x7f\xee\xc9\x4f\x2a\xc7\x4b\xfe\x8f\x8e\xa3\x7f\x6e\xc5\x23\xb3\xfb\xf1\x77\x00\x00\x00\xff\xff\x59\x8d\x9d\x4e\x7c\x0f\x00\x00")

func doc_templates_doc_md() ([]byte, error) {
	return bindata_read(
		_doc_templates_doc_md,
		"doc/templates/doc.md",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, error){
	"doc/templates/doc.md": doc_templates_doc_md,
}
// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() ([]byte, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"doc": &_bintree_t{nil, map[string]*_bintree_t{
		"templates": &_bintree_t{nil, map[string]*_bintree_t{
			"doc.md": &_bintree_t{doc_templates_doc_md, map[string]*_bintree_t{
			}},
		}},
	}},
}}
