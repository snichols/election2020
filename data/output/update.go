package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/gocarina/gocsv"
	jsoniter "github.com/json-iterator/go"
	"github.com/snichols/election2020/pkg/states"
)

type row struct {
	Time       time.Time `csv:"time"`
	Eevp       int64     `csv:"eevp"`
	VotesTotal int64     `csv:"votes_total"`
	VotesDelta int64     `csv:"votes_delta"`
	ShareBiden float64   `csv:"share_biden"`
	ShareTrump float64   `csv:"share_trump"`
	ShareOther float64   `csv:"share_other"`
	DeltaBiden int64     `csv:"delta_biden"`
	DeltaTrump int64     `csv:"delta_trump"`
	DeltaOther int64     `csv:"delta_other"`
}

func update(in string, out string) error {
	data, err := ioutil.ReadFile(in)
	if err != nil {
		return err
	}
	rows := []row{}
	timeseries := jsoniter.Get(data, "data", "races", 0, "timeseries")
	trumpVotes, bidenVotes, otherVotes, totalVotes := 0.0, 0.0, 0.0, 0.0
	for i := 0; i < timeseries.Size(); i++ {
		s := timeseries.Get(i)
		var t time.Time
		if t, err = time.Parse("2006-01-02T15:04:05Z", s.Get("timestamp").ToString()); err != nil {
			return err
		}
		v := s.Get("votes").ToFloat64()
		ts := s.Get("vote_shares", "trumpd").ToFloat64()
		bs := s.Get("vote_shares", "bidenj").ToFloat64()
		os := float64(int64((1.0-(ts+bs))*1000)) / 1000.0
		tv := v * ts
		bv := v * bs
		ov := v * os
		rows = append(rows, row{
			Time:       t,
			Eevp:       s.Get("eevp").ToInt64(),
			VotesTotal: int64(v),
			VotesDelta: int64(v - totalVotes),
			ShareBiden: bs,
			ShareTrump: ts,
			ShareOther: os,
			DeltaBiden: int64(bv - bidenVotes),
			DeltaTrump: int64(tv - trumpVotes),
			DeltaOther: int64(ov - otherVotes),
		})
		trumpVotes = tv
		bidenVotes = bv
		otherVotes = ov
		totalVotes = v
	}
	if data, err = gocsv.MarshalBytes(rows); err != nil {
		return err
	}
	if err = ioutil.WriteFile(out, data, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func main() {
	for _, n := range states.Name {
		in := fmt.Sprintf("../input/%s.json", n)
		out := fmt.Sprintf("%s.csv", n)
		fmt.Println("update:", out)
		if err := update(in, out); err != nil {
			panic(err)
		}
	}
}
