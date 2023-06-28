package main

import (
	"github.com/lubingowan/pipeline"
)

type IpGroup struct {
	Ip     string
	Dns    string
	Weight int32
}

type IdcGroup struct {
	Id       string
	Weight   int32
	Ips      []*IpGroup
	Isp      string
	WithEdge uint32
}

type IdcPriorityGroup struct {
	Priority int64
	Idcs     []*IdcGroup
}

var IdcGroups = []*IdcPriorityGroup{
	{
		Priority: 100,
		Idcs: []*IdcGroup{
			{
				Id:       "1",
				Weight:   50,
				Ips:      []*IpGroup{{Ip: "1.1.1.1", Dns: "1.1.1.1", Weight: 1}},
				Isp:      "1",
				WithEdge: 1,
			},
			{
				Id:       "1",
				Weight:   50,
				Ips:      []*IpGroup{{Ip: "5.5.5.5", Dns: "1.1.1.1", Weight: 1}},
				Isp:      "1",
				WithEdge: 1,
			},
		},
	},
	{
		Priority: 50,
		Idcs: []*IdcGroup{
			{
				Id:       "1",
				Weight:   50,
				Ips:      []*IpGroup{{Ip: "2.2.2.2", Dns: "1.1.1.1", Weight: 1}},
				Isp:      "1",
				WithEdge: 1,
			},
		},
	},
	{
		Priority: 10,
		Idcs: []*IdcGroup{
			{
				Id:       "1",
				Weight:   100,
				Ips:      []*IpGroup{{Ip: "3.3.3.3", Dns: "1.1.1.1", Weight: 1}},
				Isp:      "1",
				WithEdge: 1,
			},
		},
	},
}

func main() {
	testNormal()

	testFlatAll()

	testCommonFlow()
}

func testFlatAll() {
	s := pipeline.NewStream()
	s.Add([][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12, 13}}).FLatAll().Print("flat all").
		Add([][]int{{996, 997, 998}, {111}, {111}, {111}}).FLatAll().Print("flat all").Sort(func(i, j int) bool {
		return s.Index(i).(int) < s.Index(j).(int)
	}).Print("sorted")
}

func testCommonFlow() {
	pipeline.NewStream(IdcGroups).Flat().Print("flat once").Map(func(i interface{}) interface{} {
		ipGroup := [][]string{}
		for _, idc := range i.(*IdcPriorityGroup).Idcs {
			ips := []string{}
			for _, ip := range idc.Ips {
				ips = append(ips, ip.Ip)
			}
			ipGroup = append(ipGroup, ips)
		}
		return ipGroup
	}).Print("map once").Flat().Print("flat twice").Limit(2).Print("limited")
}

func testNormal() {
	list := [][]int{{1, 2, 3}, {4, 5, 6, 7}, {8, 9, 10}}
	pipeline.NewStream(list).Flat().Flat().Filter(func(i interface{}) bool {
		return i.(int)%2 == 0
	}).Print("filter")
}
