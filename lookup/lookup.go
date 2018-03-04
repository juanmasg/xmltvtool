package lookup

import (
    "time"
    "sort"
    "golang.org/x/text/search"
    "golang.org/x/text/language"
    "../xmltv"
    "log"
)

func Now(tv *xmltv.Tv, next bool) map[string][]*xmltv.Programme{
    chs := make(map[string][]*xmltv.Programme)

    for _, ch := range tv.Channel{
        chs[ch.Id] = make([]*xmltv.Programme, 0)
    }

    now := time.Now()
    sort.Slice(tv.Programme, func(i, j int) bool{ return tv.Programme[i].Start < tv.Programme[j].Start})

    for _, prog := range tv.Programme{
        start, _ := xmltv.ParseTime(prog.Start)
        stop, _ := xmltv.ParseTime(prog.Stop)
        if start.Before(now) && stop.After(now){
            chs[prog.Channel] = append(chs[prog.Channel], prog)
        }
    }

    return chs
}

func Search(tv *xmltv.Tv, title, subtitle, season, episode string, exact bool) []*xmltv.Programme{
    matches := make([]*xmltv.Programme, 0)

    sort.Slice(tv.Programme, func(i, j int) bool{ return tv.Programme[i].Channel + tv.Programme[i].Start < tv.Programme[j].Channel + tv.Programme[j].Start})

    searcher := search.New(language.Spanish, search.IgnoreCase, search.IgnoreDiacritics)

    for _, prog := range tv.Programme{
        start, end := searcher.IndexString(prog.Title, title)
        if start == -1 && end == -1{
            continue
        }

        if start > 0{
            continue
        }

        chname := ""
        for _, ch := range tv.Channel{
            if ch.Id == prog.Channel{
                chname = ch.Name
            }
        }

        log.Println(chname, prog)
    }

    return matches
}
