package main

import (
//    "log"
//    "os"
    "encoding/xml"
//    "io"
    "io/ioutil"
//    "errors"
    "time"
    "flag"
    "sort"
    "fmt"
)

const timefmt = "20060102150405 -0700"

type tv struct{
    XMLName     xml.Name            `xml:"tv"`
    Channel     []*XMLTVChannel     `xml:"channel"`
    Programme   []*XMLTVProgramme   `xml:"programme"`
}

func NewXMLTVFile() *tv{
    xmltvf := &tv{}
    xmltvf.Channel = make([]*XMLTVChannel, 0)
    xmltvf.Programme = make([]*XMLTVProgramme, 0)

    return xmltvf
}

type XMLTVChannel struct{
    Id      string  `xml:"id,attr"`
    Name    string  `xml:"display-name"`
}

type XMLTVProgramme struct{
    Start       string  `xml:"start,attr"`
    Stop        string  `xml:"stop,attr"`
    Channel     string  `xml:"channel,attr"`
    Title       string  `xml:"title,omitempty"`
    SubTitle    string  `xml:"sub-title,omitempty"`
    Desc        string  `xml:"desc,omitempty"`
    Date        string  `xml:"date,omitempty"`
}

func parseTime(t string) (time.Time, error){
    return time.Parse(timefmt, t)
}

func timeString(t time.Time) string{
    return t.Format(timefmt)
}

func marshal(v interface{}) ([]byte, error){
    return xml.MarshalIndent(v, "", "  ")
}

func unmarshal(data []byte, v interface{}) error{
    return xml.Unmarshal(data, v)
}

func ReadFile(path string) (*tv, error){
    data, err := ioutil.ReadFile(path); if err != nil{
        return nil, err
    }

    xmltvf := NewXMLTVFile()

    err = unmarshal(data, xmltvf); if err != nil{
        return nil, err
    }

    return xmltvf, nil
}

func WriteFile(path string, data []byte) error{
    return ioutil.WriteFile(path, data, 0644)
}

func demo() *tv{
    xmltvf := NewXMLTVFile()
    xmltvf.Channel = append(xmltvf.Channel, &XMLTVChannel{Id: "0", Name: "AAAA" })
    xmltvf.Channel = append(xmltvf.Channel, &XMLTVChannel{Id: "1", Name: "BBBB" })

    xmltvf.Programme = append(xmltvf.Programme, &XMLTVProgramme{
        Start: timeString(time.Now()), Stop: timeString(time.Now().Add(30 *time.Minute)), Channel: "0", Title: "aaaaa title", Date: "asdasd" })
    xmltvf.Programme = append(xmltvf.Programme, &XMLTVProgramme{
        Start: timeString(time.Now()), Stop: timeString(time.Now().Add(30 *time.Minute)), Channel: "1", Title: "bbbbb title", Date: "asdasd" })

    return xmltvf
}

func main(){
    flag_join := flag.Bool("j", false, "Concatenate multiple XMLTV files")
    flag_read := flag.Bool("r", false, "Read multiple XMLTV files")
    flag_demo := flag.Bool("d", false, "Demo data")
    flag_gaps := flag.Bool("g", false, "Check data gaps")
    flag_over := flag.Bool("o", false, "Check data overlaps")

    flag.Parse()

    if *flag_join{
        tv := NewXMLTVFile()
        for _, path := range flag.Args(){
            //fmt.Println(path, flag.Args())
            thistv, err := ReadFile(path); if err != nil{
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

        data, err := marshal(tv); if err != nil{
            fmt.Println(err)
        }

        fmt.Println(string(data))
    }

    if *flag_read{
        for _, path := range flag.Args(){
            tv, err := ReadFile(path); if err != nil{
                fmt.Println(err)
            }

            tv.Channel = tv.Channel[:2]
            tv.Programme = tv.Programme[:2]

            data, err := marshal(tv)
            fmt.Println(string(data))
        }
    }

    if *flag_demo{
        data, err := marshal(demo())
        fmt.Println(string(data), err)
    }

    if *flag_gaps{
        path := flag.Args()[0]
        tv, err := ReadFile(path); if err != nil{
            fmt.Println(err)
        }

        for _, c := range tv.Channel{
            progs := make([]*XMLTVProgramme, 0)
            for _, p := range tv.Programme{
                if p.Channel == c.Id{
                    progs = append(progs, p)
                }
            }
            fmt.Println("Channel", c, "has", len(progs), "programs")
            sort.Slice(progs, func(i, j int) bool{ return progs[i].Start < progs[j].Start })
            var prev *XMLTVProgramme
            for _, p := range progs{
                if prev == nil{
                    prev = p
                    continue
                }

                stop, _ := parseTime(prev.Stop)
                start, _ := parseTime(p.Start)

                if start.Sub(stop) > (5 * time.Minute){
                    fmt.Println("\tGAP!", start.Sub(stop), "from", stop, "until", start, prev.Title, " -> ", p.Title)
                }

                //fmt.Println(start.Sub(stop))

                //if prev.Stop != p.Start{
                //    fmt.Println(prev.Stop, p.Start, prev.Stop == p.Start)
                //    fmt.Printf("%+v\n", p)
                //}

                prev = p
            }
        }

    }

    if *flag_over{
    }
}
