package time_stamp_logger

import (
	"encoding/csv"
	"fmt"
	"main/consts"
	"os"
	"regexp"
	"strconv"
	"time"
)

const tFormat = "2006-01-02T15:04:05.999999999Z07:00"

var (
	re = regexp.MustCompile("[:-TZ]")
	w  *csv.Writer
)

func Initial(library string) {
	now := time.Now()

	fileName := consts.LOG_DIR + library + re.ReplaceAllString(now.Format(tFormat), "") + ".tsv"
	file, err := os.Create(fileName)
	if err != nil {
		panic("failed to create log file")
	}
	w = csv.NewWriter(file)
	w.Comma = '\t'
	if err := w.Write([]string{"id", "in", "out", "delta(nano seconds)"}); err != nil {
		panic("failed to create header")
	}
}

func Write(id int, in, out time.Time) {
	delta := out.Sub(in)
	fmt.Println(in.Format(tFormat), out.Format(tFormat), strconv.FormatInt(delta.Nanoseconds(), 10))
	if err := w.Write([]string{strconv.Itoa(id), in.Format(tFormat), out.Format(tFormat), strconv.FormatInt(delta.Nanoseconds(), 10)}); err != nil {
		fmt.Printf("failed to write data.%v\n", err)
	}
}

func Stop() {
	w.Flush()
}
