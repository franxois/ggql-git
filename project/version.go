package project

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
)

type Version struct {
	Release                 string
	Major, Minor, Patch, Rc int
	IsRc                    bool
}

func (v Version) FullVer() string {

	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.IsRc {
		s = fmt.Sprintf("%s-%d", s, v.Rc)
	}
	return s
}

type byVersions []Version

func (a byVersions) Len() int      { return len(a) }
func (a byVersions) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byVersions) Less(i, j int) bool {

	if a[i].Major != a[j].Major {
		return a[i].Major < a[j].Major
	}
	if a[i].Minor != a[j].Minor {
		return a[i].Minor < a[j].Minor
	}
	if a[i].Patch != a[j].Patch {
		return a[i].Patch < a[j].Patch
	}

	if a[i].IsRc && a[j].IsRc {
		return a[i].Rc < a[j].Rc
	}

	return a[i].IsRc && !a[j].IsRc
}

func tagListToVersions(tags []string) []Version {

	var versions []Version

	for _, versionTxt := range tags {

		r := regexp.MustCompile(`v(?P<fullVer>(?P<release>(?P<major>[\d]+)\.(?P<minor>[\d]+))\.(?P<patch>[\d]+))(-(?P<rc>[\d]*))*`)
		matches := r.FindStringSubmatch(versionTxt)

		if len(matches) >= 8 {

			v := Version{Release: matches[2]}

			major, err := strconv.Atoi(matches[3])

			if err == nil {
				v.Major = major
			}

			minor, err := strconv.Atoi(matches[4])
			if err == nil {
				v.Minor = minor
			}
			patch, err := strconv.Atoi(matches[5])
			if err == nil {
				v.Patch = patch
			}

			if matches[7] != "" {
				rc, err := strconv.Atoi(matches[7])
				if err == nil {
					v.IsRc = true
					v.Rc = rc
					versions = append(versions, v)
				}
			} else {
				versions = append(versions, v)
			}
		}
	}

	sort.Sort(byVersions(versions))
	return versions
}
