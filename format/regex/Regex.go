package regex

import "regexp"

var Regex_IPAddress = regexp.MustCompile(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.){3}(25[0-5]|(2[0-4]|1\d|[1-9]|)\d)(:([1-9]|[1-9]\d|[1-9]\d\d|[1-9]\d\d\d|[1-5]\d\d\d\d|6[0-4]\d\d\d|65[0-4]\d\d|655[0-2]\d|6553[0-6])){0,1}$`)
