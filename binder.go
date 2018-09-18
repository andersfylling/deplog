package deplog

const (
	FlagDebug uint64 = 1 << iota
	FlagDebuga
	FlagDebugf
	FlagInfo
	FlagInfoa
	FlagInfof
	FlagWarning
	FlagWarninga
	FlagWarningf
	FlagWarn
	FlagWarna
	FlagWarnf
	FlagError
	FlagErrora
	FlagErrorf
	FlagCritical
	FlagCriticala
	FlagCriticalf
	FlagCritic
	FlagCritica
	FlagCriticf
	FlagCrit
	FlagCrita
	FlagCritf
	FlagPrint
	FlagPrinta
	FlagPrintf
	FlagNotice
	FlagNoticea
	FlagNoticef
)

const (
	ProfileLogrus uint64 = 0
)

func createProfile(logger interface{}) (profile uint64) {
	var k bool
	if _, k = logger.(Debugger); k {
		profile |= FlagDebug
	}
	k = false
	if _, k = logger.(ArgsDebugger); k {
		profile |= FlagDebuga
	}
	k = false
	if _, k = logger.(FormatDebugger); k {
		profile |= FlagDebugf
	}
	k = false
	if _, k = logger.(Infoer); k {
		profile |= FlagInfo
	}
	k = false
	if _, k = logger.(ArgsInfoer); k {
		profile |= FlagInfoa
	}
	k = false
	if _, k = logger.(FormatInfoer); k {
		profile |= FlagInfof
	}
	k = false
	if _, k = logger.(Warninger); k {
		profile |= FlagWarning
	}
	k = false
	if _, k = logger.(ArgsWarninger); k {
		profile |= FlagWarninga
	}
	k = false
	if _, k = logger.(FormatWarninger); k {
		profile |= FlagWarningf
	}
	k = false
	if _, k = logger.(Warner); k {
		profile |= FlagWarn
	}
	k = false
	if _, k = logger.(ArgsWarner); k {
		profile |= FlagWarna
	}
	k = false
	if _, k = logger.(FormatWarner); k {
		profile |= FlagWarnf
	}
	k = false
	if _, k = logger.(Errorer); k {
		profile |= FlagError
	}
	k = false
	if _, k = logger.(ArgsErrorer); k {
		profile |= FlagErrora
	}
	k = false
	if _, k = logger.(FormatErrorer); k {
		profile |= FlagErrorf
	}
	k = false
	if _, k = logger.(Criticaler); k {
		profile |= FlagCritical
	}
	k = false
	if _, k = logger.(ArgsCriticaler); k {
		profile |= FlagCriticala
	}
	k = false
	if _, k = logger.(FormatCriticaler); k {
		profile |= FlagCriticalf
	}
	k = false
	if _, k = logger.(Criticer); k {
		profile |= FlagCritical
	}
	k = false
	if _, k = logger.(ArgsCriticer); k {
		profile |= FlagCriticala
	}
	k = false
	if _, k = logger.(FormatCriticer); k {
		profile |= FlagCriticalf
	}
	k = false
	if _, k = logger.(Criter); k {
		profile |= FlagCritical
	}
	k = false
	if _, k = logger.(ArgsCriter); k {
		profile |= FlagCriticala
	}
	k = false
	if _, k = logger.(FormatCriter); k {
		profile |= FlagCriticalf
	}
	k = false
	if _, k = logger.(Noticer); k {
		profile |= FlagNotice
	}
	k = false
	if _, k = logger.(ArgsNoticer); k {
		profile |= FlagNoticea
	}
	k = false
	if _, k = logger.(FormatNoticer); k {
		profile |= FlagNoticef
	}

	return
}

func getFuncPointerFromFlag(flag uint64, logger interface{}) (fp interface{}) {
	switch flag {
	case FlagDebug:
		fp = (logger.(Debugger)).Debug
	case FlagDebuga:
		fp = (logger.(ArgsDebugger)).Debug
	case FlagDebugf:
		fp = (logger.(FormatDebugger)).Debugf
	case FlagInfo:
		fp = (logger.(Infoer)).Info
	case FlagInfoa:
		fp = (logger.(ArgsInfoer)).Info
	case FlagInfof:
		fp = (logger.(FormatInfoer)).Infof
	case FlagWarning:
		fp = (logger.(Warninger)).Warning
	case FlagWarninga:
		fp = (logger.(ArgsWarninger)).Warning
	case FlagWarningf:
		fp = (logger.(FormatWarninger)).Warningf
	case FlagWarn:
		fp = (logger.(Warner)).Warn
	case FlagWarna:
		fp = (logger.(ArgsWarner)).Warn
	case FlagWarnf:
		fp = (logger.(FormatWarner)).Warnf
	case FlagError:
		fp = (logger.(Errorer)).Error
	case FlagErrora:
		fp = (logger.(ArgsErrorer)).Error
	case FlagErrorf:
		fp = (logger.(FormatErrorer)).Errorf
	case FlagCritical:
		fp = (logger.(Criticaler)).Critical
	case FlagCriticala:
		fp = (logger.(ArgsCriticaler)).Critical
	case FlagCriticalf:
		fp = (logger.(FormatCriticaler)).Criticalf
	case FlagCritic:
		fp = (logger.(Criticer)).Critic
	case FlagCritica:
		fp = (logger.(ArgsCriticer)).Critic
	case FlagCriticf:
		fp = (logger.(FormatCriticer)).Criticf
	case FlagCrit:
		fp = (logger.(Criter)).Crit
	case FlagCrita:
		fp = (logger.(ArgsCriter)).Crit
	case FlagCritf:
		fp = (logger.(FormatCriter)).Critf
	case FlagNotice:
		fp = (logger.(Noticer)).Notice
	case FlagNoticea:
		fp = (logger.(ArgsNoticer)).Notice
	case FlagNoticef:
		fp = (logger.(FormatNoticer)).Noticef
	}

	return
}

// BindRoutes will bind the existing logger methods in the injected logger with
// the given map. This map tells which of the DepLog methods should be routed
// to which injected logger methods, or ignore if they don't exist.
// see SetupDefaultRoutes(...) as an example
func BindRoutes(l *DepLog, profile uint64, levels map[uint64]uint64) {
	if l.profile != profile {
		return
	}

	l.routes = map[uint64](interface{}){} // clear routes

	// link the diff log methods if they exist
	links := map[uint64]uint64{}
	for lvl, methods := range levels {
		link := methods & l.profile
		if link > 0 {
			links[lvl] = link
		}
	}

	// binds the function pointer target to the desired level
	// if there exists more than one Flag in the link, choose the first one
	for lvl, fpt := range links {
		Flag := nextMethodProfile(fpt, 0)
		fp := getFuncPointerFromFlag(Flag, l.injectedLogger)
		l.routes[lvl] = fp
	}
}

// SetupDefaultRoutes confiure the DepLog to only use logger methods that
// exists in the injected logger. If they do not exist, logging is ignored.
func SetupDefaultRoutes(l *DepLog) {
	// each key or index, represents an method in DepLog
	// order does not matter, as each method uses their own Flag to find
	// which function pointer they can call in the routes map
	levels := map[uint64]uint64{
		FlagInfo:   FlagInfo | FlagInfoa | FlagPrint | FlagPrinta,
		FlagInfof:  FlagInfof | FlagPrintf,
		FlagWarn:   FlagWarn | FlagWarning | FlagWarna | FlagWarninga,
		FlagWarnf:  FlagWarnf | FlagWarningf,
		FlagDebug:  FlagDebug | FlagDebuga,
		FlagDebugf: FlagDebugf,
		FlagError:  FlagError | FlagErrora,
		FlagErrorf: FlagErrorf,
		FlagCrit:   FlagCrit | FlagCritic | FlagCritical | FlagCrita | FlagCritica | FlagCriticala,
		FlagCritf:  FlagCritf | FlagCriticf | FlagCriticalf,
	}

	BindRoutes(l, l.profile, levels)
}
