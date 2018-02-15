package xmltv

import (
    "encoding/xml"
    "io/ioutil"
    "time"
)

const timefmt = "20060102150405 -0700"

type tv struct{
    XMLName     xml.Name            `xml:"tv"`
    Channel     []*Channel     `xml:"channel"`
    Programme   []*Programme   `xml:"programme"`
}

func NewXMLTVFile() *tv{
    xmltvf := &tv{}
    xmltvf.Channel = make([]*Channel, 0)
    xmltvf.Programme = make([]*Programme, 0)

    return xmltvf
}

type Channel struct{
    Id      string  `xml:"id,attr"`
    Name    string  `xml:"display-name"`
}

type Programme struct{
    Start       string  `xml:"start,attr"`
    Stop        string  `xml:"stop,attr"`
    Channel     string  `xml:"channel,attr"`
    Title       string  `xml:"title,omitempty"`
    SubTitle    string  `xml:"sub-title,omitempty"`
    Desc        string  `xml:"desc,omitempty"`
    Date        string  `xml:"date,omitempty"`
}

func ParseTime(t string) (time.Time, error){
    return time.Parse(timefmt, t)
}

func TimeString(t time.Time) string{
    return t.Format(timefmt)
}

func Marshal(v interface{}) ([]byte, error){
    data, err := xml.MarshalIndent(v, "", "  "); if err != nil{
        return data, err
    }

    data = append([]byte(xml.Header), data...)

    return data, err
}

func Unmarshal(data []byte, v interface{}) error{
    return xml.Unmarshal(data, v)
}

func ReadFile(path string) (*tv, error){
    data, err := ioutil.ReadFile(path); if err != nil{
        return nil, err
    }

    xmltvf := NewXMLTVFile()

    err = Unmarshal(data, xmltvf); if err != nil{
        return nil, err
    }

    return xmltvf, nil
}

func WriteFile(path string, data []byte) error{
    return ioutil.WriteFile(path, data, 0644)
}

func Demo() *tv{
    xmltvf := NewXMLTVFile()
    xmltvf.Channel = append(xmltvf.Channel, &Channel{Id: "0", Name: "AAAA" })
    xmltvf.Channel = append(xmltvf.Channel, &Channel{Id: "1", Name: "BBBB" })

    xmltvf.Programme = append(xmltvf.Programme, &Programme{
        Start: TimeString(time.Now()), Stop: TimeString(time.Now().Add(30 *time.Minute)), Channel: "0", Title: "aaaaa title", Date: "asdasd" })
    xmltvf.Programme = append(xmltvf.Programme, &Programme{
        Start: TimeString(time.Now()), Stop: TimeString(time.Now().Add(30 *time.Minute)), Channel: "1", Title: "bbbbb title", Date: "asdasd" })

    return xmltvf
}

