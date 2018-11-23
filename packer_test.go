package main

import (
	"basetools"
	"testing"
)

func TestPWDCreater(t *testing.T) {
	pwd := newDirPwd("/home/aliy/study/alal/")
	var P = pwd.GetPackedPWD()

	basetools.AssertEqual(t, 5, len(P))
	basetools.AssertEqual(t, P[0].PWD, "/")
	basetools.AssertEqual(t, P[0].Name, "ROOT")

	basetools.AssertEqual(t, P[1].PWD, "/home")
	basetools.AssertEqual(t, P[1].Name, "home")

	basetools.AssertEqual(t, P[2].PWD, "/home/aliy")
	basetools.AssertEqual(t, P[2].Name, "aliy")

	basetools.AssertEqual(t, P[3].PWD, "/home/aliy/study")
	basetools.AssertEqual(t, P[3].Name, "study")

	basetools.AssertEqual(t, P[4].PWD, "/home/aliy/study/alal")
	basetools.AssertEqual(t, P[4].Name, "alal")
}

func TestPWDCreaterBad(t *testing.T) {
	var pwd = newDirPwd("")
	var P = pwd.GetPackedPWD()

	basetools.AssertEqual(t, 1, len(P))
	basetools.AssertEqual(t, P[0].PWD, "/")
	basetools.AssertEqual(t, P[0].Name, "ROOT")

	pwd = newDirPwd("/")
	P = pwd.GetPackedPWD()

	basetools.AssertEqual(t, 1, len(P))
	basetools.AssertEqual(t, P[0].PWD, "/")
	basetools.AssertEqual(t, P[0].Name, "ROOT")

	pwd = newDirPwd("/home")
	P = pwd.GetPackedPWD()

	basetools.AssertEqual(t, 2, len(P))
	basetools.AssertEqual(t, P[0].PWD, "/")
	basetools.AssertEqual(t, P[0].Name, "ROOT")

	basetools.AssertEqual(t, P[1].PWD, "/home")
	basetools.AssertEqual(t, P[1].Name, "home")

}

func TestSomeBasic(t *testing.T) {
	var w, h, wh = 0, 0, 0

	w = 0x33
	h = 0x44
	wh = ((w & 0xffff) << 16) | (h & 0xffff)
	basetools.AssertEqual(t, wh, 0x00330044)

	tw := (wh & 0xffff0000) >> 16
	th := wh & 0x0000ffff
	basetools.AssertEqual(t, tw, w)
	basetools.AssertEqual(t, th, h)

}
