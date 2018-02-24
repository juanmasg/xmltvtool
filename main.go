package main

import (
    "flag"
    "sort"
    "fmt"
    "time"
    "strings"
    "./xmltv"
    "./watcher"
    "./actions"
    "./lookup"
)

func main(){
    flag_join := flag.Bool("j", false, "Concatenate multiple XMLTV files")
    flag_read := flag.Bool("r", false, "Read multiple XMLTV files")
    flag_demo := flag.Bool("d", false, "Demo data")
    flag_gaps := flag.Bool("g", false, "Check data gaps")
    flag_over := flag.Bool("o", false, "Check data overlaps")

    flag_watch := flag.Bool("w", false, "Watch for program schedules")
    flag_action := flag.Bool("a", false, "Action to execute for matched program schedules")
    flag_now := flag.Bool("now", false, "Print airing now")

    flag_search := flag.Bool("s", false, "Search programme")
    flag_search_title := flag.String("title", "", "Search this title")
    flag_search_subtitle := flag.String("subtitle", "", "Search this subtitle")
    flag_search_season := flag.String("season", "", "Search this season")
    flag_search_episode := flag.String("episode", "", "Search this episode")

    flag.Parse()

    if *flag_join{
        tv := xmltv.NewXMLTVFile()

        existingc := make(map[string]bool)

        for _, path := range flag.Args(){
            //fmt.Println(path, flag.Args())
            thistv, err := xmltv.ReadFile(path); if err != nil{
                //fmt.Println(path, err)
                continue
            }

            for _, c := range thistv.Channel{
                _, exists := existingc[c.Name]; if exists{
                    fmt.Println("Channel already exists", c)
                    continue
                }
                existingc[c.Name] = true
                tv.Channel = append(tv.Channel, c)
            }

            for _, p := range thistv.Programme{
                tv.Programme = append(tv.Programme, p)
            }
        }

        data, err := xmltv.Marshal(tv); if err != nil{
            fmt.Println(err)
        }

        fmt.Println(string(data))
    }

    if *flag_read{
        for _, path := range flag.Args(){
            tv, err := xmltv.ReadFile(path); if err != nil{
                fmt.Println(err)
            }

            tv.Channel = tv.Channel[:2]
            tv.Programme = tv.Programme[:2]

            data, err := xmltv.Marshal(tv)
            fmt.Println(string(data))
        }
    }

    if *flag_demo{
        data, err := xmltv.Marshal(xmltv.Demo())
        fmt.Println(string(data), err)
    }

    if *flag_gaps{
        path := flag.Args()[0]
        tv, err := xmltv.ReadFile(path); if err != nil{
            fmt.Println(err)
        }

        for _, c := range tv.Channel{
            progs := make([]*xmltv.Programme, 0)
            for _, p := range tv.Programme{
                if p.Channel == c.Id{
                    progs = append(progs, p)
                }
            }
            fmt.Println("Channel", c, "has", len(progs), "programs")
            sort.Slice(progs, func(i, j int) bool{ return progs[i].Start < progs[j].Start })

            if len(progs) == 0{
                fmt.Println("No programme for channel", c)
                continue
            }

            prev := progs[0]

            for _, cur := range progs[1:]{

                stop, _ := xmltv.ParseTime(prev.Stop)
                start, _ := xmltv.ParseTime(cur.Start)

                if start.Sub(stop) > (5 * time.Minute){
                    fmt.Println("\tGAP!", start.Sub(stop), "from", stop, "until", start, prev.Title, " -> ", cur.Title)
                }

                prev = cur
            }
        }

    }

    if *flag_over{
    }

    if *flag_now{
        tv, err := xmltv.ReadFile(flag.Args()[0]); if err != nil{
            fmt.Println(tv, err)
            return
        }
        now := lookup.Now(tv, false)
        for _, ch := range tv.Channel{
            nowch, found := now[ch.Id]; if !found{
                printChannelNow(ch.Id, ch.Name, nil)
                continue
            }

            printChannelNow(ch.Id, ch.Name, nowch)
        }
    }

    if *flag_watch{

        err := watcher.Watch(flag.Args(), func(path string){
            fmt.Println("!", path)
            tv, err := xmltv.ReadFile(path); if err != nil{
                fmt.Println(tv, err)
            }
            reScheduleAll(tv.Programme)
        })
        fmt.Println(err)
    }

    if *flag_search{
        tv, err := xmltv.ReadFile(flag.Args()[0]); if err != nil{
            fmt.Println(tv, err)
            return
        }

        lookup.Search(tv, *flag_search_title, *flag_search_subtitle, *flag_search_season, *flag_search_episode, false)
    }
}

func reScheduleAll(progs []*xmltv.Programme){
    for _, prog := range progs{
        if strings.HasPrefix(prog.Title, "Friends"){
            actions.Log(prog)
        }
    }
}

func printChannelNow(id, name string, progs []*xmltv.Programme){

    fmt.Printf("% 4s % -30s", id, name)

    if progs == nil || len(progs) == 0{
        fmt.Println(" now   Unknown")
        return
    }

    if len(progs) == 1{
        printProgramme(" now  ", progs[0])
    }

    if len(progs) == 2{
        printProgramme(" next ", progs[1])
    }

    fmt.Println("")
}

func printProgramme(prefix string, prog *xmltv.Programme){
    start, _ := xmltv.ParseTime(prog.Start)
    stop, _ := xmltv.ParseTime(prog.Stop)

    fmt.Printf("%s %s to %s (% 8s left) - %s",
        prefix,
        start.Format("15:04"),
        stop.Format("15:04"),
        stop.Sub(time.Now()).Truncate(1 * time.Minute),
        prog.Title,
        )
}
