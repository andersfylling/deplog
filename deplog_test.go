package deplog

import (
	"testing"
)

type LoggerA struct {
	counter int
}

func (l *LoggerA) Info(message string) {
	l.counter++
}

func TestBinding(t *testing.T) {
	var logger interface{} = &LoggerA{}
	if _, ok := logger.(Infoer); !ok {
		t.Fatal("does not implement Infoer")
	}

	_, err := NewDepLogger(logger)
	if err != nil {
		t.Error(err)
	}

	deplog := &DepLog{
		injectedLogger: logger,
		profile:        createProfile(logger),
		routes:         map[uint64](interface{}){},
	}
	SetupDefaultRoutes(deplog)

	if deplog.profile != FlagInfo {
		t.Errorf("profile was incorrect. Got %d, wants %d", deplog.profile, FlagInfo)
	}

	deplog.Info("")
	if logger.(*LoggerA).counter != 1 {
		t.Errorf("incorrect counter. Got %d, wants %d", logger.(*LoggerA).counter, 1)
	}
}

type LoggerB struct {
	a int
	b int
}

func (l *LoggerB) Info(message string) {
	l.a++
}

func (l *LoggerB) Debug(message string) {
	l.b++
}

func TestRerouting(t *testing.T) {
	instance := &LoggerB{}
	var logger interface{} = instance

	if _, ok := logger.(Infoer); !ok {
		t.Fatal("does not implement Infoer")
	}
	if _, ok := logger.(Debugger); !ok {
		t.Fatal("does not implement Debugger")
	}

	deplog := &DepLog{
		injectedLogger: logger,
		profile:        createProfile(logger),
		routes:         map[uint64](interface{}){},
	}
	SetupDefaultRoutes(deplog)

	if deplog.profile != (FlagInfo | FlagDebug) {
		t.Errorf("profile was incorrect. Got %d, wants %d", deplog.profile, (FlagInfo | FlagDebug))
	}

	deplog.Info("")
	if instance.a != 1 && instance.b != 0 {
		t.Errorf("incorrect counter. Got %d, wants %d", instance.a, 1)
	}

	deplog.Debug("")
	if instance.a != 1 && instance.b != 1 {
		t.Errorf("incorrect counter. Got %d, wants %d", instance.b, 1)
	}

	// route debug to info
	deplog.Route(deplog.profile, FlagDebug, FlagInfo)

	deplog.Debug("")
	deplog.Debug("")
	if instance.b != 1 {
		t.Errorf("incorrect counter. Got %d, wants %d", instance.b, 1)
	}
	if instance.a != 3 {
		t.Errorf("incorrect counter. Got %d, wants %d", instance.a, 3)
	}
}
