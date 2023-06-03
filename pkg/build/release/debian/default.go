package debian

const DefaultTargetMountpoint = "build/mount"

var Default = &Debian{
	Suite:  "bullseye",
	Mirror: "http://deb.debian.org/debian",
}
