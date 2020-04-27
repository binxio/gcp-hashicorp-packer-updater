package main

import "regexp"

type PublicImageProject struct {
	Project       string
	FamilyPattern *regexp.Regexp
}

type PublicImageProjects []PublicImageProject

var PublicProjects PublicImageProjects = PublicImageProjects{
	PublicImageProject{"centos-cloud", regexp.MustCompile("^centos.*$")},
	PublicImageProject{"google-containers", regexp.MustCompile("^container-vm.*$")},
	PublicImageProject{"coreos-cloud", regexp.MustCompile("^coreos*$")},
	PublicImageProject{"cos-cloud", regexp.MustCompile("^cos.*$")},
	PublicImageProject{"debian-cloud", regexp.MustCompile("^debian.*$")},
	PublicImageProject{"rhel-cloud", regexp.MustCompile("^rhel.*$")},
	PublicImageProject{"suse-cloud", regexp.MustCompile("^sles.*$")},
	PublicImageProject{"ubuntu-os-cloud", regexp.MustCompile("^ubuntu.*$")},
	PublicImageProject{"windows-cloud", regexp.MustCompile("^windows.*$")},
	PublicImageProject{"goog-vmruntime-images", regexp.MustCompile("^gae-builder-vm.*$")},
	PublicImageProject{"opensuse-cloud", regexp.MustCompile("^.*$")}}

func (m *PublicImageProjects) FindProjectForName(name string) *PublicImageProject {
	if name == "" {
		return nil
	}
	for _, p := range *m {
		if p.FamilyPattern.MatchString(name) {
			return &p
		}
	}
	return nil
}
