package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
	jsoniter "github.com/json-iterator/go"
	"github.com/snichols/election2020/pkg/states"
)

type row struct {
	Time       string  `csv:"time"`
	Eevp       int64   `csv:"pct"`
	VotesTotal int64   `csv:"votes"`
	ShareBiden float64 `csv:"share\nbiden"`
	ShareTrump float64 `csv:"share\ntrump"`
	ShareOther float64 `csv:"share\nother"`
	TotalBiden int64   `csv:"votes\nbiden"`
	TotalTrump int64   `csv:"votes\ntrump"`
	TotalOther int64   `csv:"votes\nother"`
	BatchVotes int64   `csv:"batch size"`
	BatchBiden int64   `csv:"batch\nbiden"`
	BatchTrump int64   `csv:"batch\ntrump"`
	BatchOther int64   `csv:"batch\nother"`
	Note       string  `csv:"note"`
}

func update(in string, out string) error {
	data, err := ioutil.ReadFile(in)
	if err != nil {
		return err
	}
	var est *time.Location
	if est, err = time.LoadLocation("EST"); err != nil {
		return err
	}
	rows := []row{}
	timeseries := jsoniter.Get(data, "data", "races", 0, "timeseries")
	lastTrumpVotes, lastBidenVotes, lastOtherVotes, lastTotalVotes := 0.0, 0.0, 0.0, 0.0
	for i := 0; i < timeseries.Size(); i++ {
		sample := timeseries.Get(i)

		// get the total votes
		totalVotes := sample.Get("votes").ToFloat64()
		if totalVotes == 0 {
			// ignore samples with no vote impact
			continue
		}

		// compute the maximum variation in vote total due to 0.001 precision
		maxVariation := int64(totalVotes * 0.001)

		// get the sample time in EST timezone
		var sampleTime time.Time
		if sampleTime, err = time.Parse("2006-01-02T15:04:05Z", sample.Get("timestamp").ToString()); err != nil {
			return err
		}
		sampleTime = sampleTime.In(est)

		// get shares for trump, biden, and "other"
		trumpShare := sample.Get("vote_shares", "trumpd").ToFloat64()
		bidenShare := sample.Get("vote_shares", "bidenj").ToFloat64()
		otherShare := 1.0 - (trumpShare + bidenShare)

		// truncate otherShare to 0.001 precision
		otherShare = float64(int64(otherShare*1000.0)) / 1000.0

		// compute vote totals for biden, trump, and "other"
		bidenVotes := totalVotes * bidenShare
		trumpVotes := totalVotes * trumpShare
		otherVotes := totalVotes * otherShare

		// compute vote batch size overall and for trump, biden, and "other"
		batchTotal := int64(totalVotes - lastTotalVotes)
		batchTrump := int64(trumpVotes - lastTrumpVotes)
		batchBiden := int64(bidenVotes - lastBidenVotes)
		batchOther := int64(otherVotes - lastOtherVotes)

		// detect anomalies
		notes := []string{}

		if batchTotal < 0 {
			notes = append(notes, fmt.Sprintf("Total %d", batchTotal))
		}

		if batchBiden < -maxVariation {
			notes = append(notes, fmt.Sprintf("Biden %d", batchBiden))
		}

		if batchTrump < -maxVariation {
			notes = append(notes, fmt.Sprintf("Trump %d", batchTrump))
		}

		if batchOther < -maxVariation {
			notes = append(notes, fmt.Sprintf("Other %d", batchOther))
		}

		// add row to CSV data
		rows = append(rows, row{
			Time:       sampleTime.Format("2006-01-02 15:04:05"),
			Eevp:       sample.Get("eevp").ToInt64(),
			VotesTotal: int64(totalVotes),
			ShareBiden: bidenShare,
			ShareTrump: trumpShare,
			ShareOther: otherShare,
			TotalBiden: int64(bidenVotes),
			TotalTrump: int64(trumpVotes),
			TotalOther: int64(otherVotes),
			BatchVotes: batchTotal,
			BatchBiden: batchBiden,
			BatchTrump: batchTrump,
			BatchOther: batchOther,
			Note:       strings.Join(notes, "\n"),
		})

		// update last values
		lastTrumpVotes = trumpVotes
		lastBidenVotes = bidenVotes
		lastOtherVotes = otherVotes
		lastTotalVotes = totalVotes
	}

	// save rows to CSV file
	{
		if data, err = gocsv.MarshalBytes(rows); err != nil {
			return err
		}
		if err = ioutil.WriteFile(out, data, os.ModePerm); err != nil {
			return err
		}
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
