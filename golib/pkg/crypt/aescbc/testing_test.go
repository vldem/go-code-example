package aescbc

import "testing"

type aescbcTestFixture struct {
	key1                     string
	key2                     string
	iv                       string
	plainText                string
	securedTextBase64        string
	securedTextHex           string
	securedTextBase64NotSafe string
	securedTextHexNotSafe    string
	isSafeMode               bool
}

func setUp(t *testing.T) aescbcTestFixture {
	var fixture aescbcTestFixture
	fixture.key1 = `ynbUGdgBeTxYaNYhM+FeSgS2BA+9dS5kZ0NOzMhZRq0=`
	fixture.key2 = `8rtPClzS4M5BjE47mr1VqDrmjfF+LsDOur0mAmuDaEUPNIuAZKLdVeiqq7RtOdmadMTO0oCobx5do/5ib884qA==`
	fixture.iv = `xJPIKJecyJvc0qzoRe7hgw==`
	fixture.plainText = `Complex2Pass.`
	fixture.securedTextBase64 = `xJPIKJecyJvc0qzoRe7hg3Tof8RhJW9Fe1GmlK+5Hzgvmmf0kdaEH6zeln+eyiXT72sqPmWDQDDpb54+knDNyO0MHNG+FpBEIme46E+YmzaRKW5qn4QyNYM6kWkTbiFP`
	fixture.securedTextHex = `c493c828979cc89bdcd2ace845eee18374e87fc461256f457b51a694afb91f382f9a67f491d6841facde967f9eca25d3ef6b2a3e65834030e96f9e3e9270cdc8ed0c1cd1be1690442267b8e84f989b3691296e6a9f843235833a9169136e214f`
	fixture.securedTextBase64NotSafe = `xJPIKJecyJvc0qzoRe7hg5EpbmqfhDI1gzqRaRNuIU8=`
	fixture.securedTextHexNotSafe = `c493c828979cc89bdcd2ace845eee18391296e6a9f843235833a9169136e214f`
	fixture.isSafeMode = true

	return fixture
}

// func (f *usersRepoFixtures) tearDown() {

// }
