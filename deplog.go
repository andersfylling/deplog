package deplog

import (
	"errors"
	"fmt"
	"math/bits"
	"sync"
)

// log converts the logger interface to one of the func types above
// this func assume that the routing is valid
func log(fp interface{}, str string, args ...interface{}) {
	if fp == nil {
		return
	}

	switch f := fp.(type) {
	case func(string):
		f(str)
	case func(...interface{}):
		newArgs := []interface{}{}
		newArgs = append(newArgs, str)
		newArgs = append(newArgs, args...)
		f(newArgs...)
	case func(string, ...interface{}):
		f(str, args...)
	}
}

func lsb(x uint64) int {
	return bits.TrailingZeros64(x)
}

// nextMethodProfile is a Next Least Significant Bit algorithm
func nextMethodProfile(x, ignore uint64) uint64 {
	x ^= ignore
	return uint64(1) << uint64(lsb(x))
}

func getLoggerFunc(routes map[uint64](interface{}), flag uint64) (f interface{}) {
	var methods uint64
	for methods, f = range routes {
		if (methods & flag) > 0 {
			break
		}
	}

	return
}

type Logger interface {
	Info(message string)
	Infof(format string, args ...interface{})
	Warn(message string)
	Warnf(format string, args ...interface{})
	Debug(message string)
	Debugf(format string, args ...interface{})
	Error(message string)
	Errorf(message string, args ...interface{})
	Crit(message string)
	Critf(message string, args ...interface{})
}

func NewDepLogger(injected interface{}) (logger Logger, err error) {
	if injected == nil {
		err = errors.New("injected logger can not be nil")
		return
	}

	profile := createProfile(injected)
	deplog := &DepLog{
		injectedLogger: injected,
		profile:        profile,
		routes:         map[uint64](interface{}){},
	}
	SetupDefaultRoutes(deplog)

	logger = deplog
	return
}

type DepLog struct {
	// injectedLogger instance
	injectedLogger interface{}

	// profile holds the profile of the injected logger
	profile uint64

	// routes holds relationship between deplog logging methods and the
	// methods which exists within the injected logger.
	// eg. calling DepLog.Debug will check the router to see if Debug is
	//     binded to a method of the injected logger
	routes map[uint64](interface{})

	mu sync.RWMutex
}

// Route overwrites existing default routes for s specific logger type such as logrus
// eg. DepLog.Route(profileLogrus, flagDebugf, flagFatalf)
// will bind the DepLog.Debugf to injected.Fatalf, if the injected logger has the matching profile
//
// Yes, another injected logger migth have the same profile as logrus, in that case
// the routes will still take affect. This only cares about the profile, not the GoLang type
func (l *DepLog) Route(loggerProfile uint64, depLogMethods uint64, injectedMethod uint64) *DepLog {
	l.mu.Lock()
	defer l.mu.Unlock()

	// ensure this profile regards the current logger isntance
	if (loggerProfile & l.profile) != loggerProfile {
		return l
	}

	// if this solution already exists, skip
	solution := injectedMethod | depLogMethods
	if _, exists := l.routes[solution]; exists {
		return l
	}

	// 1. find the route(key) for the injectedMethod's method pointer
	// 2. for each binary index in depLogMethods, remove the index from the
	//    current keys or routes (after: if key == 0 then delete)
	// 3. add the depLogMethods to the route found in step.1

	// 1. check that the flag actually has a method pointer
	solutionfp := getLoggerFunc(l.routes, injectedMethod)
	if solutionfp == nil {
		return l
	}
	fmt.Printf("\n::::%+v\n", l.routes)

	// 2.remove the old methods
	methods := map[uint64]uint64{}
	for route := range l.routes {
		flags := route & (depLogMethods | injectedMethod)
		if flags > 0 {
			methods[route] = flags ^ route
		}
	}
	fmt.Printf("\n::::%+v\n", l.routes)

	// 2. overwrite
	for old, new := range methods {
		fp := l.routes[old]
		delete(l.routes, old)
		l.routes[new] = fp
	}
	fmt.Printf("\n::::%+v\n", l.routes)

	// remove dead key
	if _, exists := l.routes[0]; exists {
		delete(l.routes, 0)
	}

	// 3. add the new flags to the desired route
	l.routes[solution] = solutionfp
	fmt.Printf("\n::::%+v\n", l.routes)

	return l
}

func (l *DepLog) Info(message string) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	f := getLoggerFunc(l.routes, flagInfo)
	log(f, message)
}

func (l *DepLog) Infof(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	f := getLoggerFunc(l.routes, flagInfof)
	log(f, format, args...)
}

func (l *DepLog) Warn(message string) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	f := getLoggerFunc(l.routes, flagWarn|flagWarning)
	log(f, message)
}

func (l *DepLog) Warnf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	f := getLoggerFunc(l.routes, flagWarnf|flagWarningf)
	log(f, format, args...)
}

func (l *DepLog) Debug(message string) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	f := getLoggerFunc(l.routes, flagDebug)
	log(f, message)
}

func (l *DepLog) Debugf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	f := getLoggerFunc(l.routes, flagDebugf)
	log(f, format, args...)
}
func (l *DepLog) Error(message string) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	f := getLoggerFunc(l.routes, flagError)
	log(f, message)
}
func (l *DepLog) Errorf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	f := getLoggerFunc(l.routes, flagErrorf)
	log(f, format, args...)
}
func (l *DepLog) Crit(message string) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	f := getLoggerFunc(l.routes, flagCrit)
	log(f, message)
}
func (l *DepLog) Critf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	f := getLoggerFunc(l.routes, flagCritf)
	log(f, format, args...)
}
