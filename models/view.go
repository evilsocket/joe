package models

import (
	"fmt"
	"github.com/evilsocket/islazy/fs"
	"github.com/evilsocket/islazy/log"
	"github.com/wcharczuk/go-chart"
	"io"
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

func cleanAllBut(basePath, exclude, expr string) {
	var err error
	err = fs.Glob(basePath, expr, func(fileName string) error {
		if fileName, err = filepath.Abs(fileName); err != nil {
			return err
		} else if fileName != exclude {
			log.Debug("removing %s", fileName)
			if err = os.Remove(fileName); err != nil {
				log.Error("error removing %s: %v", fileName, err)
			}
		}
		return nil
	})
	if err != nil {
		log.Error("%v", err)
	}
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
		view.SourceFileName = viewFileName
		view.NativeFileName = path.Join(basePath, fmt.Sprintf(
			"%s_%s.so",
			queryName,
			viewName))

		// check if the file has already been compiled
		if fs.Exists(view.NativeFileName) == false {
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