package file_watcher

import (
	"github.com/fsnotify/fsnotify"
	"github.com/meowalien/go-meowalien-lib/errs"
	"io"
	"log"
	"path/filepath"
	"sync"
)

type FileWatcher interface {
	io.Closer
	Add(name string) error
	Remove(name string) error
	Start() error
}

type fileWatcher struct {
	lock            sync.RWMutex
	allWatchingFile map[string]watchingFile
	*fsnotify.Watcher
	onFileUpDate func(file string)
	onFileRemove func(file string)
}

type watchingFile struct {
	fileName string
	dir      string
	realFile string
}

func (f *fileWatcher) Remove(filename string) (err error) {
	f.lock.Lock()
	defer f.lock.Unlock()
	cleanFilePath := filepath.Clean(filename)
	wFile, exist := f.allWatchingFile[cleanFilePath]
	if !exist {
		return
	}
	delete(f.allWatchingFile, cleanFilePath)
	err = f.Watcher.Remove(wFile.dir)
	if err != nil {
		err = errs.New(err)
		return
	}
	return
}

func (f *fileWatcher) Add(filename string) (err error) {
	f.lock.Lock()
	defer f.lock.Unlock()
	cleanFilePath := filepath.Clean(filename)
	dir, _ := filepath.Split(cleanFilePath)
	realFilePath, _ := filepath.EvalSymlinks(filename)

	f.allWatchingFile[cleanFilePath] = watchingFile{
		fileName: cleanFilePath,
		dir:      dir,
		realFile: realFilePath,
	}

	err = f.Watcher.Add(dir)
	if err != nil {
		err = errs.New(err)
		return
	}
	return
}

func (f *fileWatcher) Start() (err error) {
	go func() {
		for {
			select {
			case event, ok := <-f.Watcher.Events:
				if !ok { // 'Events' channel is closed
					log.Printf("file watcher on closed\n")
					//eventsWG.Done()
					return
				}
				f.onEvent(event)

			case err, ok := <-f.Watcher.Errors:
				if ok { // 'Errors' channel is not closed
					log.Printf("watcher error: %v\n", err)
				}
			}
		}
	}()
	return
}

func (f *fileWatcher) onEvent(event fsnotify.Event) {
	f.lock.RLock()
	defer f.lock.RUnlock()
	for _, wFile := range f.allWatchingFile {
		currentConfigFile, _ := filepath.EvalSymlinks(wFile.fileName)
		// we only care about the config file with the following cases:
		// 1 - if the config file was modified or created
		// 2 - if the real path to the config file changed (eg: k8s ConfigMap replacement)
		const writeOrCreateMask = fsnotify.Write | fsnotify.Create
		if (filepath.Clean(event.Name) == wFile.fileName &&
			event.Op&writeOrCreateMask != 0) ||
			(currentConfigFile != "" && currentConfigFile != wFile.realFile) {
			wFile.realFile = currentConfigFile
			f.onFileUpDate(wFile.fileName)
		} else if filepath.Clean(event.Name) == wFile.fileName &&
			event.Op&fsnotify.Remove&fsnotify.Remove != 0 {
			f.onFileRemove(wFile.fileName)
			return
		}
	}
}

func New(onFileUpDate func(file string), onFileRemove func(file string)) (watcher FileWatcher, err error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		err = errs.New(err)
		return
	}
	watcher = &fileWatcher{
		Watcher:         w,
		allWatchingFile: map[string]watchingFile{},
		onFileUpDate:    onFileUpDate,
		onFileRemove:    onFileRemove,
	}
	return
}

//func NewWatcher(filename string, onFileUpDate func(file string), onFileRemove func(file string)) (newWatcher file_watcher.FileWatcher, err error) {
//	watcher, err := fsnotify.NewWatcher()
//	if err != nil {
//		err = errs.New(err)
//		return
//	}
//	configFile := filepath.Clean(filename)
//	configDir, _ := filepath.Split(configFile)
//	realConfigFile, _ := filepath.EvalSymlinks(filename)
//	go func() {
//		for {
//			select {
//			case event, ok := <-watcher.Events:
//				if !ok { // 'Events' channel is closed
//					log.Printf("file watcher on %s closed\n", filename)
//					//eventsWG.Done()
//					return
//				}
//				currentConfigFile, _ := filepath.EvalSymlinks(filename)
//				// we only care about the config file with the following cases:
//				// 1 - if the config file was modified or created
//				// 2 - if the real path to the config file changed (eg: k8s ConfigMap replacement)
//				const writeOrCreateMask = fsnotify.Write | fsnotify.Create
//				if (filepath.Clean(event.Name) == configFile &&
//					event.Op&writeOrCreateMask != 0) ||
//					(currentConfigFile != "" && currentConfigFile != realConfigFile) {
//					realConfigFile = currentConfigFile
//					onFileUpDate(configFile)
//				} else if filepath.Clean(event.Name) == configFile &&
//					event.Op&fsnotify.Remove&fsnotify.Remove != 0 {
//					onFileRemove(configFile)
//					return
//				}
//
//			case err, ok := <-watcher.Errors:
//				if ok { // 'Errors' channel is not closed
//					log.Printf("watcher error: %v\n", err)
//				}
//				log.Printf("file watcher on %s closed\n", filename)
//				return
//			}
//		}
//	}()
//	err = watcher.Add(configDir)
//	if err != nil {
//		err = errs.New(err)
//		return
//	}
//	newWatcher = &fileWatcher{Watcher: watcher}
//	return
//}
