package main_test

import (
	"net/http"
	"testing"
)

func BenchmarkGetEvents(b *testing.B) {
	url := "http://localhost:8080/events"
	for i := 0; i < b.N; i++ {
		resp, err := http.Get(url)
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()
	}
}

func BenchmarkGetEventsParallel(b *testing.B) {
	url := "http://localhost:8080/events"
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := http.Get(url)
			if err != nil {
				b.Fatal(err)
			}
			resp.Body.Close()
		}
	})
}

func BenchmarkGetEventDetails(b *testing.B) {
	url := "http://localhost:8080/events/sand-castle"
	for i := 0; i < b.N; i++ {
		resp, err := http.Get(url)
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()
	}
}

func BenchmarkGetEventDetailsParallel(b *testing.B) {
	url := "http://localhost:8080/events/sand-castle"
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := http.Get(url)
			if err != nil {
				b.Fatal(err)
			}
			resp.Body.Close()
		}
	})
}
