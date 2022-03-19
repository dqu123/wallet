package extend

import "time"

func parseExtendTime(t string) *time.Time {
	// 2022-03-18T01:05:17.000+0000
	output, err := time.Parse("2006-01-02T03:04:05.000+0000", t)
	if err != nil {
		return nil
	}
	return &output
}
