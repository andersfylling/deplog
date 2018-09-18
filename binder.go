package deplog

const (
	flagDebug uint64 = 1 << iota
	flagDebuga
	flagDebugf
	flagInfo
	flagInfoa
	flagInfof
	flagWarning
	flagWarninga
	flagWarningf
	flagWarn
	flagWarna
	flagWarnf
	flagError
	flagErrora
	flagErrorf
	flagCritical
	flagCriticala
	flagCriticalf
	flagCritic
	flagCritica
	flagCriticf
	flagCrit
	flagCrita
	flagCritf
	flagPrint
	flagPrinta
	flagPrintf
	flagNotice
	flagNoticea
	flagNoticef
)

const (
	ProfileLogrus uint64 = 0
)

func createProfile(logger interface{}) (profile uint64) {
	var k bool
	if _, k = logger.(Debugger); k {
		profile |= flagDebug
	}
	k = false
	if _, k = logger.(ArgsDebugger); k {
		profile |= flagDebuga
	}
	k = false
	if _, k = logger.(FormatDebugger); k {
		profile |= flagDebugf
	}
	k = false
	if _, k = logger.(Infoer); k {
		profile |= flagInfo
	}
	k = false
	if _, k = logger.(ArgsInfoer); k {
		profile |= flagInfoa
	}
	k = false
	if _, k = logger.(FormatInfoer); k {
		profile |= flagInfof
	}
	k = false
	if _, k = logger.(Warninger); k {
		profile |= flagWarning
	}
	k = false
	if _, k = logger.(ArgsWarninger); k {
		profile |= flagWarninga
	}
	k = false
	if _, k = logger.(FormatWarninger); k {
		profile |= flagWarningf
	}
	k = false
	if _, k = logger.(Warner); k {
		profile |= flagWarn
	}
	k = false
	if _, k = logger.(ArgsWarner); k {
		profile |= flagWarna
	}
	k = false
	if _, k = logger.(FormatWarner); k {
		profile |= flagWarnf
	}
	k = false
	if _, k = logger.(Errorer); k {
		profile |= flagError
	}
	k = false
	if _, k = logger.(ArgsErrorer); k {
		profile |= flagErrora
	}
	k = false
	if _, k = logger.(FormatErrorer); k {
		profile |= flagErrorf
	}
	k = false
	if _, k = logger.(Criticaler); k {
		profile |= flagCritical
	}
	k = false
	if _, k = logger.(ArgsCriticaler); k {
		profile |= flagCriticala
	}
	k = false
	if _, k = logger.(FormatCriticaler); k {
		profile |= flagCriticalf
	}
	k = false
	if _, k = logger.(Criticer); k {
		profile |= flagCritical
	}
	k = false
	if _, k = logger.(ArgsCriticer); k {
		profile |= flagCriticala
	}
	k = false
	if _, k = logger.(FormatCriticer); k {
		profile |= flagCriticalf
	}
	k = false
	if _, k = logger.(Criter); k {
		profile |= flagCritical
	}
	k = false
	if _, k = logger.(ArgsCriter); k {
		profile |= flagCriticala
	}
	k = false
	if _, k = logger.(FormatCriter); k {
		profile |= flagCriticalf
	}
	k = false
	if _, k = logger.(Noticer); k {
		profile |= flagNotice
	}
	k = false
	if _, k = logger.(ArgsNoticer); k {
		profile |= flagNoticea
	}
	k = false
	if _, k = logger.(FormatNoticer); k {
		profile |= flagNoticef
	}

	return
}

func getFuncPointerFromFlag(flag uint64, logger interface{}) (fp interface{}) {
	switch flag {
	case flagDebug:
		fp = (logger.(Debugger)).Debug
	case flagDebuga:
		fp = (logger.(ArgsDebugger)).Debug
	case flagDebugf:
		fp = (logger.(FormatDebugger)).Debugf
	case flagInfo:
		fp = (logger.(Infoer)).Info
	case flagInfoa:
		fp = (logger.(ArgsInfoer)).Info
	case flagInfof:
		fp = (logger.(FormatInfoer)).Infof
	case flagWarning:
		fp = (logger.(Warninger)).Warning
	case flagWarninga:
		fp = (logger.(ArgsWarninger)).Warning
	case flagWarningf:
		fp = (logger.(FormatWarninger)).Warningf
	case flagWarn:
		fp = (logger.(Warner)).Warn
	case flagWarna:
		fp = (logger.(ArgsWarner)).Warn
	case flagWarnf:
		fp = (logger.(FormatWarner)).Warnf
	case flagError:
		fp = (logger.(Errorer)).Error
	case flagErrora:
		fp = (logger.(ArgsErrorer)).Error
	case flagErrorf:
		fp = (logger.(FormatErrorer)).Errorf
	case flagCritical:
		fp = (logger.(Criticaler)).Critical
	case flagCriticala:
		fp = (logger.(ArgsCriticaler)).Critical
	case flagCriticalf:
		fp = (logger.(FormatCriticaler)).Criticalf
	case flagCritic:
		fp = (logger.(Criticer)).Critic
	case flagCritica:
		fp = (logger.(ArgsCriticer)).Critic
	case flagCriticf:
		fp = (logger.(FormatCriticer)).Criticf
	case flagCrit:
		fp = (logger.(Criter)).Crit
	case flagCrita:
		fp = (logger.(ArgsCriter)).Crit
	case flagCritf:
		fp = (logger.(FormatCriter)).Critf
	case flagNotice:
		fp = (logger.(Noticer)).Notice
	case flagNoticea:
		fp = (logger.(ArgsNoticer)).Notice
	case flagNoticef:
		fp = (logger.(FormatNoticer)).Noticef
	}

	return
}

// BindRoutes will bind the existing logger methods in the injected logger with
// the given map. This map tells which of the DepLog methods should be routed
// to which injected logger methods, or ignore if they don't exist.
// see SetupDefaultRoutes(...) as an example
func BindRoutes(l *DepLog, levels []uint64) {
	// link the diff log methods if they exist
	links := map[uint64]uint64{}
	for _, lvl := range levels {
		link := lvl & l.profile
		if link > 0 {
			links[lvl] = link
		}
	}

	// binds the function pointer target to the desired level
	// if there exists more than one flag in the link, choose the first one
	for lvl, fpt := range links {
		flag := nextMethodProfile(fpt, 0)
		fp := getFuncPointerFromFlag(flag, l.injectedLogger)
		l.routes[lvl] = fp
	}
}

// SetupDefaultRoutes confiure the DepLog to only use logger methods that
// exists in the injected logger. If they do not exist, logging is ignored.
func SetupDefaultRoutes(l *DepLog) {
	// each key or index, represents an method in DepLog
	// order does not matter, as each method uses their own flag to find
	// which function pointer they can call in the routes map
	levels := []uint64{
		/*info*/ flagInfo | flagInfoa | flagPrint | flagPrinta,
		/*infof*/ flagInfof | flagPrintf,
		/*warn*/ flagWarn | flagWarning | flagWarna | flagWarninga,
		/*warnf*/ flagWarnf | flagWarningf,
		/*debug*/ flagDebug | flagDebuga,
		/*debugf*/ flagDebugf,
		/*err*/ flagError | flagErrora,
		/*errf*/ flagErrorf,
		/*crit*/ flagCrit | flagCritic | flagCritical | flagCrita | flagCritica | flagCriticala,
		/*critf*/ flagCritf | flagCriticf | flagCriticalf,
	}

	BindRoutes(l, levels)
}
