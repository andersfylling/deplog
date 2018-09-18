# DepLog
[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)[![forthebadge](https://forthebadge.com/images/badges/for-you.svg)](https://forthebadge.com)

## Health
| Branch       | Build status  | Code climate | Go Report Card |
| ------------ |:-------------:|:---------------:|:-------------:|
| master     | [![CircleCI](https://circleci.com/gh/andersfylling/deplog/tree/master.svg?style=shield)](https://circleci.com/gh/andersfylling/deplog/tree/master) | [![Maintainability](https://api.codeclimate.com/v1/badges/91efae0f1e83a0c5fcf1/maintainability)](https://codeclimate.com/github/andersfylling/deplog/maintainability) | [![Go Report Card](https://goreportcard.com/badge/github.com/andersfylling/deplog)](https://goreportcard.com/report/github.com/andersfylling/deplog) |

## About
As I've been writing other libraries/modules, I noticed the need to support injecting a logger such that the log files can be somewhat consistent for a user. However, there seems to be many different loggers out there and no standard interface. This project is not another logger, it is a wrapper. A wrapper for any kind of logger out there (although, it needs more testing).

DepLog identifies the single logging methods that are commonly used for any logger injected. It then creates a profile which is used to link the DepLog logger methods to the methods of the injected logger. It also supports writing profiles for any logger out there and customize behaviour.

The goal of this project is to remove the habit of forcing users of libraries to utilize one specific logger, and let them choose which ever they prefer, while still keeping a general yet expressive logger interface for the library developers.

## Setup
> **Note:* that DepLog does not have a exported package variable to store your logger instance at. If you do want this, I recommend creating your own sub-package called logger and keep all your logger configuration in there.


```GoLang
// let's pretend the injected logger only has .Debug as a method
// this will create a profile showing it has a Debug method only
// profile == deplog.FlagDebug
// this profile is used later in rerouting

// create a deplogger instance which has mapped the default logger behaviour
logger, err := NewDepLogger(someRandomLoggerInstance)
if err != nil {
  panic(err) // someRandomLoggerInstance was nil
}
deplog.SetupDefaultRoutes(logger)

logger.Debug("test") // will trigger someRandomLoggerInstance.Debug(...)
logger.Info("yep") // not called, as the injected logger is missing Info method

// let's route the Info method to the .Debug method
logger.Route(deplog.FlagDebug, FlagInfo, FlagDebug)
logger.Debug("this prints!") // triggers someRandomLoggerInstance.Debug(...)
```
This might be a little confusing. It's important to know that calling NewDepLogger will initialize a default relationship, where methods missing in the injected logger, are simply ignored. However, as we saw, we can route deplog methods to another logger method as desired per profile using the DepLog.Route option.

DepLog.Route takes tree parameters:
 1. profile: the profile for a specific logger
 2. src: the method or methods to be rerouted
 3. dest: the method to which `src` are routed

The reason I used deplog.FlagDebug, is because it is the specific profile for any logger that only has a .Debug method.

You can also create more advanced behaviour using the deplog.BindRoutes(...) which takes the DepLog instance with an array of logger levels (info, debug, etc.). It is here you decided which logger level should be vound to which logger method. Here I show the default network. The first line of the map tells us that when we run a DepLog.Info method, it should be routed to first found logger method that is either: .Info(string), .Info(...interface{}), .Print(string), .Print(...interface{}).
While DepLog.Debugf is only routed to .Debugf(string, ...interface{})
```GoLang
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
deplog.BindRoutes(deplogger, deplogger.profile, levels)
```
> **Note:* the use of deplogger.profile in the BindRoutes function, causes the network to bind to any kind of profile, and is not limited to say just a logrus profile.

## When should I use which logger level?
Either refer to the way the different loggers utilise them, or check out similar discussion in Java.
https://stackoverflow.com/questions/5817738/how-to-use-log-levels-in-java

I know this is not Java, but the logging concept of levels is can be reused as they are based in social understanding of the words and not uniquely defined by a programming language.

## TODO
The goal is also to create network profiles for each of the most commonly used loggers out there (logrus, go-logger, google/logger, etc.), and make sure these route the methods in a intuitive way as possible so you can use DepLog right out of the box.
