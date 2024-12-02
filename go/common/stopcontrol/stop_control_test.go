package stopcontrol

import (
	"testing"
	"time"
)

func TestStopControl_Stop(t *testing.T) {
	sc := New()
	sc.Stop()

	if !sc.IsStopping() {
		t.Error("Expected IsStopping to return true after Stop, but got false")
	}

	// Ensure it's safe to call Stop multiple times
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Expected no panic when calling Stop multiple times")
			}
		}()
		sc.Stop()
	}()
}

func TestStopControl_IsStopping(t *testing.T) {
	sc := New()

	if sc.IsStopping() {
		t.Error("Expected IsStopping to return false initially, but got true")
	}

	sc.Stop()

	if !sc.IsStopping() {
		t.Error("Expected IsStopping to return true after Stop, but got false")
	}
}

func TestStopControl_Done(t *testing.T) {
	sc := New()

	select {
	case <-sc.Done():
		t.Error("Expected Done channel to be blocking initially")
	case <-time.After(50 * time.Millisecond): // Allow a small delay to check the non-blocking state
	}

	sc.Stop()

	select {
	case _, ok := <-sc.Done():
		if ok {
			t.Error("Expected Done channel to be closed after Stop")
		}
	case <-time.After(50 * time.Millisecond):
		t.Error("Expected Done channel to be closed immediately after Stop")
	}
}

func TestStopControl_OnStop(t *testing.T) {
	sc := New()

	called := false
	sc.OnStop(func() {
		called = true
	})

	if called {
		t.Error("Expected callback to not be called before Stop")
	}

	sc.Stop()

	// wait 50ms to ensure the callback is called
	time.Sleep(50 * time.Millisecond)

	if !called {
		t.Error("Expected callback to be called after Stop")
	}
}
