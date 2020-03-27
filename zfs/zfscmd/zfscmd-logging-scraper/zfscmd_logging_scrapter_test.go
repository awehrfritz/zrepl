package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseHumanFormatterNodate(t *testing.T) {

	type testCase struct {
		Name      string
		Input     string
		Expect    *RuntimeLine
		ExpectErr string
	}

	secs := func(s string) time.Duration {
		d, err := parseSecs(s)
		require.NoError(t, err)
		return d
	}

	tcs := []testCase{
		{
			Name:  "human-formatter-noerror",
			Input: `[jobname][zfs.cmd]: command exited without error usertime_s="0.008445" cmd="zfs list -H -p -o name -r -t filesystem,volume" systemtime_s="0.033783" invocation="84" total_time_s="0.037828619"`,
			Expect: &RuntimeLine{
				Cmd:       "zfs list -H -p -o name -r -t filesystem,volume",
				TotalTime: secs("0.037828619"),
				Usertime:  secs("0.008445"),
				Systime:   secs("0.033783"),
				Error:     "",
			},
		},
		{
			Name:  "human-formatter-witherror",
			Input: `[jobname][zfs.cmd]: command exited with error usertime_s="0.008445" cmd="zfs list -H -p -o name -r -t filesystem,volume" systemtime_s="0.033783" invocation="84" total_time_s="0.037828619" err="some error"`,
			Expect: &RuntimeLine{
				Cmd:       "zfs list -H -p -o name -r -t filesystem,volume",
				TotalTime: secs("0.037828619"),
				Usertime:  secs("0.008445"),
				Systime:   secs("0.033783"),
				Error:     "some error",
			},
		},
		{
			Name:  "from graylog",
			Input: `[csnas][zfs.cmd]:  command  exited  without  error  usertime_s="0"  cmd="zfs  send  -i  zroot/ezjail/synapse-12@zrepl_20200329_095518_000  zroot/ezjail/synapse-12@zrepl_20200329_102454_000"  total_time_s="0.101598591"  invocation="85"  systemtime_s="0.041581"`,
			Expect: &RuntimeLine{
				Cmd:       "zfs  send  -i  zroot/ezjail/synapse-12@zrepl_20200329_095518_000  zroot/ezjail/synapse-12@zrepl_20200329_102454_000",
				TotalTime: secs("0.101598591"),
				Systime:   secs("0.041581"),
				Usertime:  secs("0"),
				Error:     "",
			},
		},
	}

	for _, c := range tcs {
		t.Run(c.Name, func(t *testing.T) {
			l, err := parseHumanFormatterNodate(c.Input)
			t.Logf("l=%v", l)
			t.Logf("err=%T %v", err, err)
			if (c.Expect != nil && c.ExpectErr != "") || (c.Expect == nil && c.ExpectErr == "") {
				t.Fatal("bad test case", c)
			}
			if c.Expect != nil {
				require.Equal(t, *c.Expect, l)
			}
			if c.ExpectErr != "" {
				require.EqualError(t, err, c.ExpectErr)
			}
		})
	}

}
