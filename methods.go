package deplog

type Debugger interface {
	Debug(message string)
}

type ArgsDebugger interface {
	Debug(args ...interface{})
}

type FormatDebugger interface {
	Debugf(format string, args ...interface{})
}

type Infoer interface {
	Info(message string)
}

type ArgsInfoer interface {
	Info(args ...interface{})
}

type FormatInfoer interface {
	Infof(format string, args ...interface{})
}

type Warninger interface {
	Warning(message string)
}

type ArgsWarninger interface {
	Warning(args ...interface{})
}

type FormatWarninger interface {
	Warningf(format string, args ...interface{})
}

type Warner interface {
	Warn(message string)
}

type ArgsWarner interface {
	Warn(args ...interface{})
}

type FormatWarner interface {
	Warnf(format string, args ...interface{})
}

type Errorer interface {
	Error(message string)
}

type ArgsErrorer interface {
	Error(args ...interface{})
}

type FormatErrorer interface {
	Errorf(format string, args ...interface{})
}

type Criticaler interface {
	Critical(message string)
}

type ArgsCriticaler interface {
	Critical(args ...interface{})
}

type FormatCriticaler interface {
	Criticalf(format string, args ...interface{})
}

type Criticer interface {
	Critic(message string)
}

type ArgsCriticer interface {
	Critic(args ...interface{})
}

type FormatCriticer interface {
	Criticf(format string, args ...interface{})
}

type Criter interface {
	Crit(message string)
}

type ArgsCriter interface {
	Crit(args ...interface{})
}

type FormatCriter interface {
	Critf(format string, args ...interface{})
}

type Noticer interface {
	Notice(message string)
}

type ArgsNoticer interface {
	Notice(args ...interface{})
}

type FormatNoticer interface {
	Noticef(format string, args ...interface{})
}
