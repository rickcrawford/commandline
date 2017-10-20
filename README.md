Command Line Boilerplate
------------------------

This is a golang project that shows usage of Viper and Cobra for a simple command line application. This application 
supports SIG_HUP and SIG_TERM signals to restart and exit the application.


Sample usage:

```bash
$ go run main.go -h

____ ____ _    _    ___  __   ___    __   ____ __   ____
| __\|   ||\/\ |\/\ |  \ | \|\|  \   | |  |___\| \|\| __\
| \__| . ||   \|   \| . \|  \|| . \  | |__| /  |  \||  ]_
|___/|___/|/v\/|/v\/|/\_/|/\_/|___/  |___/|/   |/\_/|___/

Usage:
  commandline [flags]

Flags:
      --config string                 /path/to/config.yml
  -h, --help                          help for commandline
  -q, --quiet                         Quiet mode. Do not display banner messages
      --shutdown-wait-time duration   Shutdown wait time (default 500ms)

```

Here is an example sending `^C`

```bash
$ go run main.go --shutdown-wait-time 10s

____ ____ _    _    ___  __   ___    __   ____ __   ____
| __\|   ||\/\ |\/\ |  \ | \|\|  \   | |  |___\| \|\| __\
| \__| . ||   \|   \| . \|  \|| . \  | |__| /  |  \||  ]_
|___/|___/|/v\/|/v\/|/\_/|/\_/|___/  |___/|/   |/\_/|___/


Started with process id: 33495
2017/10/20 09:26:23 Start!
2017/10/20 09:26:23 myApplicationMain
^C2017/10/20 09:26:26 myApplicationMain: shutdown started...
2017/10/20 09:26:31 myApplicationMain: shutdown completed.
2017/10/20 09:26:36 Exit!
```

There are dependencies on the following git projects:

* [fsnotify](https://github.com/fsnotify/fsnotify)
* [cobra](https://github.com/spf13/cobra)
* [viper](https://github.com/spf13/viper)

