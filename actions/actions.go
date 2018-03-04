package actions

import (
    "log"
    "../xmltv"
)

func Log(prog *xmltv.Programme){
    log.Printf("%+v\n", prog)
}
