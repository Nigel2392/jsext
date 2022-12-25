package router

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// Router regex delimiters
const (
	RT_PATH_VAR_PREFIX = "<<"
	RT_PATH_VAR_SUFFIX = ">>"
	RT_PATH_VAR_DELIM  = ":"
)

// Router regex types.
const (
	NameInt    = "int"
	NameString = "string"
	NameSlug   = "slug"
	NameUUID   = "uuid"
	NameAny    = "any"
	NameHex    = "hex"
)

// Router regex patterns
const (
	// Match any character
	RT_PATH_REGEX_ANY = ".+"
	// Match any number
	RT_PATH_REGEX_NUM = "[0-9]+"
	// Match any string
	RT_PATH_REGEX_STR = "[a-zA-Z]+"
	// Match any hex number
	RT_PATH_REGEX_HEX = "[0-9a-fA-F]+"
	// Match any UUID
	RT_PATH_REGEX_UUID = "[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}"
	// Match any alphanumeric string
	RT_PATH_REGEX_ALPHANUMERIC = "[0-9a-zA-Z_-]+"
)

// Routes to be registered in the router
type Route struct {
	Name              string
	name              string
	Path              string
	Callable          func(v Vars, u *url.URL)
	RegexUrl          string
	skipTrailingSlash bool
	Children          []*Route
}

func (r *Route) String() string {
	return "Route: " + r.name + " -> " + r.Path
}

func (r *Route) getRoute(name string) *Route {
	for _, route := range r.Children {
		if route.name == name {
			return route
		}
	}
	for _, route := range r.Children {
		var rt = route.getRoute(name)
		if rt != nil {
			return rt
		}
	}
	return nil
}

// Register a new route as a child of this route
// If the path does not start with a slash, and r.Path does not end in one
// then a slash is added to the path
//
//	-> Route{Path: "/api", Children: []*Route{Path: "posts/"}}
//
// Will result in the path "/api/posts"
func (r *Route) Register(name, path string, callable func(v Vars, u *url.URL)) *Route {
	if r.skipTrailingSlash && len(path) > 1 {
		path = strings.TrimSuffix(path, "/")
	}
	if !strings.HasPrefix(path, "/") && !strings.HasSuffix(r.Path, "/") {
		path = "/" + path
	} else if strings.HasPrefix(path, "/") && strings.HasSuffix(r.Path, "/") {
		path = path[1:]
	}
	path = r.Path + path
	name = r.name + ":" + name

	for _, route := range r.Children {
		if route.name == name {
			panic("Route already exists: " + name)
		}
	}

	var showNameSlice = strings.Split(name, ":")
	var showName = showNameSlice[len(showNameSlice)-1]

	var route = &Route{Name: showName, name: name, Path: path, Callable: callable, skipTrailingSlash: r.skipTrailingSlash}
	r.Children = append(r.Children, route)
	return route
}

// If the path matches the route, return true and the named capture groups
// If capture group is not named, returns $1, $2, etc.
func (r *Route) Match(path string) (bool, *Route, Vars) {
	var rpath = r.regexr("")
	var rex = regexp.MustCompile(rpath)
	var m = rex.FindStringSubmatch(path)

	// Get named capture groups
	var vars = make(Vars, len(m))
	var subNames = rex.SubexpNames()
	if len(subNames) != len(m) {
		return false, nil, nil
	}
	for i, name := range rex.SubexpNames() {
		if i != 0 && name != "" {
			vars[name] = m[i]
		}
	}
	if len(m) > 0 && m[0] == path {
		return true, r, vars
	}

	for _, child := range r.Children {
		if ok, nr, vars := child.Match(path); ok {
			return ok, nr, vars
		}
	}
	return false, nil, nil
}

// Generate the regex path for the router to match.
func (r *Route) regexr(path string) string {
	if r.RegexUrl == "" {
		var path = path + r.Path

		if path == "/" {
			r.RegexUrl = path
			return r.RegexUrl
		}

		var hasPrefixSlash = strings.HasPrefix(path, "/")
		var hasTrailingSlash = strings.HasSuffix(path, "/")
		if hasPrefixSlash {
			path = path[1:]
		}
		if hasTrailingSlash {
			path = path[:len(path)-1]
		}

		var parts = strings.Split(path, "/")
		for i, part := range parts {
			parts[i] = toRegex(part)
		}
		r.RegexUrl = strings.Join(parts, "/")
		if hasPrefixSlash {
			r.RegexUrl = "/" + r.RegexUrl
		}
		if hasTrailingSlash {
			r.RegexUrl = r.RegexUrl + "/"
		}
	}
	return r.RegexUrl
}

// Format the url based on the arguments given.
// Panics if route accepts more arguments than are given.
func (r *Route) URL(args ...any) string {
	var path = r.Path

	// If the length of the path is less than the length of the pre/suffix and the delimiter
	// then there are no variables in the path
	if len(path) <= len(RT_PATH_VAR_DELIM)+len(RT_PATH_VAR_PREFIX)+len(RT_PATH_VAR_SUFFIX) {
		return path
	}
	// Remove the first and last slash if they exist
	var hasPrefixSlash = strings.HasPrefix(path, "/")
	var hasTrailingSlash = strings.HasSuffix(path, "/")
	if hasPrefixSlash {
		path = path[1:]
	}
	if hasTrailingSlash {
		path = path[:len(path)-1]
	}
	// Split the path into parts
	var parts = strings.Split(path, "/")
	// Replace the parts that are variables with the arguments
	for i, part := range parts {
		if strings.HasPrefix(part, RT_PATH_VAR_PREFIX) && strings.HasSuffix(part, RT_PATH_VAR_SUFFIX) {
			if len(args) == 0 {
				panic("not enough arguments for URL: " + r.name)
			}
			var arg = args[0]
			args = args[1:]
			parts[i] = fmt.Sprintf("%v", arg)
		}
	}
	// Join the parts back together
	path = strings.Join(parts, "/")
	// Add the slashes back if they were there
	if hasPrefixSlash {
		path = "/" + path
	}
	if hasTrailingSlash {
		path = path + "/"
	}
	return path
}

// Convert a string to a regex string with a capture group.
func toRegex(str string) string {
	if !strings.HasPrefix(str, RT_PATH_VAR_PREFIX) || !strings.HasSuffix(str, RT_PATH_VAR_SUFFIX) {
		return str
	}
	str = strings.TrimPrefix(str, RT_PATH_VAR_PREFIX)
	str = strings.TrimSuffix(str, RT_PATH_VAR_SUFFIX)
	var parts = strings.Split(str, RT_PATH_VAR_DELIM)
	if len(parts) == 1 {
		return "(?P<" + parts[0] + ">" + typToRegx(parts[0]) + ")"
	} else if len(parts) != 2 {
		return str
	}
	var groupName = parts[0]
	var typ = parts[1]
	return "(?P<" + groupName + ">" + typToRegx(typ) + ")"
}

// Convert a type (string) to a regex for use in capture groups.
func typToRegx(typ string) string {
	// regex for raw is: raw(REGEX)
	var hasRaw string = strings.ToLower(typ)
	if strings.HasPrefix(hasRaw, "raw(") && strings.HasSuffix(hasRaw, ")") {
		return hasRaw[4 : len(hasRaw)-1]
	}
	switch typ {
	case NameInt:
		return RT_PATH_REGEX_NUM
	case NameString, NameSlug:
		return RT_PATH_REGEX_ALPHANUMERIC
	case NameUUID:
		return RT_PATH_REGEX_UUID
	case NameAny:
		return RT_PATH_REGEX_ANY
	case NameHex:
		return RT_PATH_REGEX_HEX
	default:
		return RT_PATH_REGEX_STR
	}
}
