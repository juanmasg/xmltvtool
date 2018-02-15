package main

import (
    "flag"
    "sort"
    "fmt"
    "time"
    "./xmltv"
)

func main(){
    flag_join := flag.Bool("j", false, "Concatenate multiple XMLTV files")
    flag_read := flag.Bool("r", false, "Read multiple XMLTV files")
    flag_demo := flag.Bool("d", false, "Demo data")
    flag_gaps := flag.Bool("g", false, "Check data gaps")
    flag_over := flag.Bool("o", false, "Check data overlaps")

    flag.Parse()

    if *flag_join{
        tv := xmltv.NewXMLTVFile()
        for _, path := range flag.Args(){
            //fmt.Println(path, flag.Args())
            thistv, err := xmltv.ReadFile(path); if err != nil{
                //fmt.Println(path, err)
                continue
            }

            for _, c := range thistv.Channel{
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
            progs := make([]*xmltv.XMLTVProgramme, 0)
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
}
