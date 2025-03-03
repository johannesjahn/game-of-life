package main

import (
	"flag"
	"os"
	"testing"
)

func TestFlags(t *testing.T) {
	tests := []struct {
		args     []string
		expected struct {
			height   int
			width    int
			interval int
			living   int
			seed     int
			factions int
		}
	}{
		{
			args: []string{"-height=10", "-width=20", "-interval=200", "-living=50", "-seed=123", "-factions=2"},
			expected: struct {
				height   int
				width    int
				interval int
				living   int
				seed     int
				factions int
			}{10, 20, 200, 50, 123, 2},
		},
		{
			args: []string{"-h=15", "-w=25", "-i=150", "-l=60", "-s=456", "-f=3"},
			expected: struct {
				height   int
				width    int
				interval int
				living   int
				seed     int
				factions int
			}{15, 25, 150, 60, 456, 3},
		},
		{
			args: []string{"-height=30", "-width=40", "-interval=300", "-living=100", "-seed=789", "-factions=4"},
			expected: struct {
				height   int
				width    int
				interval int
				living   int
				seed     int
				factions int
			}{30, 40, 300, 100, 789, 4},
		},
		{
			args: []string{"-h=5", "-w=10", "-i=50", "-l=20", "-s=321", "-f=1"},
			expected: struct {
				height   int
				width    int
				interval int
				living   int
				seed     int
				factions int
			}{5, 10, 50, 20, 321, 1},
		},
		{
			args: []string{"-height=50", "-width=60", "-interval=500", "-living=200", "-seed=654", "-factions=5"},
			expected: struct {
				height   int
				width    int
				interval int
				living   int
				seed     int
				factions int
			}{50, 60, 500, 200, 654, 5},
		},
		{
			args: []string{},
			expected: struct {
				height   int
				width    int
				interval int
				living   int
				seed     int
				factions int
			}{-1, -1, 100, -1, 0, 0},
		},
	}

	for _, test := range tests {
		os.Args = append([]string{"cmd"}, test.args...)
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		c := parseInput()

		if c.height != test.expected.height {
			t.Errorf("expected height %d, got %d", test.expected.height, c.height)
		}
		if c.width != test.expected.width {
			t.Errorf("expected width %d, got %d", test.expected.width, c.width)
		}
		if c.interval != test.expected.interval {
			t.Errorf("expected interval %d, got %d", test.expected.interval, c.interval)
		}
		if c.living != test.expected.living {
			t.Errorf("expected living %d, got %d", test.expected.living, c.living)
		}
		if c.seed != test.expected.seed {
			t.Errorf("expected seed %d, got %d", test.expected.seed, c.seed)
		}
		if c.factions != test.expected.factions {
			t.Errorf("expected factions %d, got %d", test.expected.factions, c.factions)
		}
	}
}
