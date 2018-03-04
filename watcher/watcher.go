package watcher

import (
    "github.com/howeyc/fsnotify"
)

func Watch(files []string, callback func(path string)) error{
    w, err := fsnotify.NewWatcher(); if err != nil{
        return err
    }

    defer w.Close()

    for _, file := range files{
        err = w.Watch(file); if err != nil{
            return err
        }
    }

    for{
        select{
        case e := <-w.Event:
            if e.IsModify() || e.IsCreate(){
                callback(e.Name)
            }
        case err := <-w.Error:
            return err
        }
    }
}
