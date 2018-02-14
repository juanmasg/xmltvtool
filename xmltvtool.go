package main

import (
    "log"
//    "os"
    "encoding/xml"
//    "io"
    "io/ioutil"
//    "errors"
    "time"
    "flag"
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
    flag_demo := flag.Bool("d", false, "Generate demo data")

    flag.Parse()

    if *flag_join{
        tv := NewXMLTVFile()
        for _, path := range flag.Args(){
            log.Println(path, flag.Args())
            thistv, err := ReadFile(path); if err != nil{
                log.Println(path, err)
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
            log.Println(err)
        }

        log.Println(string(data))
    }

    if *flag_read{
        for _, path := range flag.Args(){
            tv, err := ReadFile(path); if err != nil{
                log.Println(err)
            }

            tv.Channel = tv.Channel[:2]
            tv.Programme = tv.Programme[:2]

            data, err := marshal(tv)
            log.Println(string(data))
        }
    }

    if *flag_demo{
        data, err := marshal(demo())
        log.Println(string(data), err)
    }
}

