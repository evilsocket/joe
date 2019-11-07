package models

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/log"
	"github.com/wcharczuk/go-chart"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"plugin"
	"strings"
)

type Chart interface {
	Render(rp chart.RendererProvider, w io.Writer) error
}

type View struct {
	Name           string
	SourceFileName string
	NativeFileName string

	plugin *plugin.Plugin
	cb     func(*Results) Chart
}

func PrepareView(queryName, viewName, viewFileName string) (view *View, err error) {
	view = &View{
		Name: viewName,
	}

	if viewFileName, err = filepath.Abs(viewFileName); err != nil {
		return
	}
	basePath := path.Dir(viewFileName)

	if strings.HasSuffix(viewFileName, ".go") {
		// hash the file
		var raw []byte
		raw, err = ioutil.ReadFile(viewFileName)
		if err != nil {
			return
		}

		hash := sha256.New()
		hash.Write(raw)
		hex.EncodeToString(hash.Sum(nil))

		view.SourceFileName = viewFileName
		view.NativeFileName = path.Join(basePath, fmt.Sprintf(
			"%s_%s_%s.so",
			queryName,
			viewName,
			hex.EncodeToString(hash.Sum(nil))))

		// check if the file has already been compiled
		if fs.Exists(view.NativeFileName) == false {
			// TODO: check for older versions to remove

			goPath, err := exec.LookPath("go")
			if err != nil {
				return nil, fmt.Errorf("go not found, can't compile %s", viewFileName)
			}

			log.Info("compiling %s ...", viewFileName)

			cmdLine := fmt.Sprintf("%s build -buildmode=plugin -o '%s' '%s'",
				goPath,
				view.NativeFileName,
				view.SourceFileName)

			log.Debug("%s", cmdLine)

			cmd := exec.Command("sh", "-c", cmdLine)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Env = os.Environ()

			if err := cmd.Run(); err != nil {
				return nil, err
			}
		}
	} else {
		view.SourceFileName = viewFileName
		view.NativeFileName = viewFileName
	}

	log.Info("loading view %s ...", view.NativeFileName)

	if view.plugin, err = plugin.Open(view.NativeFileName); err != nil {
		return nil, err
	}

	f, err := view.plugin.Lookup("View")
	if err != nil {
		return nil, err
	}

	var ok bool
	if view.cb, ok = f.(func(*Results) Chart); !ok {
		return nil, fmt.Errorf("can't cast %+v to func(*Results) Chart", f)
	}

	log.Debug("f = %+v", view.cb)

	return view, nil
}

func (v *View) Call(res *Results) Chart {
	return v.cb(res)
}