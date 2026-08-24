package assemble

import "testing"

func TestProxyAllowedSchemes(t *testing.T) {
	for _, proxy := range []string{
		"http://127.0.0.1:8080",
		"https://127.0.0.1:8080",
		"socks5://127.0.0.1:9050",
		"sock5://127.0.0.1:9050",
	} {
		if !proxyAllowed(proxy) {
			t.Fatalf("proxyAllowed(%q) = false, want true", proxy)
		}
	}
}

func TestBuildProxyAndClientConfiguresKnownProxyDialers(t *testing.T) {
	for _, proxy := range []string{
		"http://127.0.0.1:8080",
		"https://127.0.0.1:8080",
		"socks5://127.0.0.1:9050",
		"sock5://127.0.0.1:9050",
	} {
		client := BuildProxyAndClient(proxy, 0, false, false)
		if client.Dial == nil {
			t.Fatalf("BuildProxyAndClient(%q) left Dial nil, want proxy dialer", proxy)
		}
	}
}

func TestBuildProxyAndClientForceUnknownProxyContinuesDirect(t *testing.T) {
	client := BuildProxyAndClient("ftp://127.0.0.1:21", 0, false, true)
	if client.Dial != nil {
		t.Fatal("BuildProxyAndClient unknown proxy with force configured Dial, want direct connection")
	}
}

func TestProxyUnknownAccepted(t *testing.T) {
	called := false
	ok := confirmUnknownProxy(false, func() bool {
		called = true
		return true
	})

	if !ok {
		t.Fatal("confirmUnknownProxy accepted = false, want true")
	}
	if !called {
		t.Fatal("confirmUnknownProxy did not ask for confirmation")
	}
}

func TestProxyUnknownRejected(t *testing.T) {
	called := false
	ok := confirmUnknownProxy(false, func() bool {
		called = true
		return false
	})

	if ok {
		t.Fatal("confirmUnknownProxy rejected = true, want false")
	}
	if !called {
		t.Fatal("confirmUnknownProxy did not ask for confirmation")
	}
}

func TestProxyForceSkipsConfirmation(t *testing.T) {
	ok := confirmUnknownProxy(true, func() bool {
		t.Fatal("confirmUnknownProxy should not ask when forceProxy is true")
		return false
	})

	if !ok {
		t.Fatal("confirmUnknownProxy force = false, want true")
	}
}
