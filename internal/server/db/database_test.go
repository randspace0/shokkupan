package db

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/randspace0/shokkupan/internal/models"
)

func TestConcurrentAuthReadsAndStatsWrites(t *testing.T) {
	d, err := Init(filepath.Join(t.TempDir(), "shokku.db"))
	if err != nil {
		t.Fatal(err)
	}
	pool, err := d.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { pool.Close() })
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	d = d.WithContext(ctx)
	const clientID = "signin-test-client"
	if err := d.Create(&models.GithubApp{ClientId: clientID}).Error; err != nil {
		t.Fatal(err)
	}

	const workers, samples = 8, 25
	start := make(chan struct{})
	results := make(chan error, workers)
	for worker := 0; worker < workers; worker++ {
		go func() {
			<-start
			for sample := 0; sample < samples; sample++ {
				if err := d.Create(&models.AppStatsSample{AppName: "signin-test"}).Error; err != nil {
					results <- fmt.Errorf("recording stats: %w", err)
					return
				}
				var app models.GithubApp
				if err := d.First(&app).Error; err != nil {
					results <- fmt.Errorf("reading GitHub metadata: %w", err)
					return
				}
				if app.ClientId != clientID {
					results <- fmt.Errorf("client ID = %q, want %q", app.ClientId, clientID)
					return
				}
			}
			results <- nil
		}()
	}
	close(start)
	for worker := 0; worker < workers; worker++ {
		if err := <-results; err != nil {
			t.Error(err)
		}
	}
	var count int64
	if err := d.Model(&models.AppStatsSample{}).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != workers*samples {
		t.Fatalf("recorded samples = %d, want %d", count, workers*samples)
	}
}
