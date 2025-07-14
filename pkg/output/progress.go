package output

import (
	"fmt"
	"sync/atomic"
	"time"
)

// ProgressSpinner displays a spinning progress indicator
// Extracted and adapted from firescan's progress display logic
type ProgressSpinner struct {
	message    string
	isRunning  *int32
	done       chan bool
}

// NewSpinner creates a new progress spinner
func NewSpinner(message string) *ProgressSpinner {
	var running int32
	return &ProgressSpinner{
		message:   message,
		isRunning: &running,
		done:      make(chan bool),
	}
}

// Start begins the spinner animation
func (p *ProgressSpinner) Start() {
	atomic.StoreInt32(p.isRunning, 1)
	go p.spin()
}

// Stop stops the spinner animation
func (p *ProgressSpinner) Stop() {
	if atomic.LoadInt32(p.isRunning) == 1 {
		atomic.StoreInt32(p.isRunning, 0)
		p.done <- true
		// Clear the line
		fmt.Printf("\r%80s\r", "")
	}
}

// UpdateMessage updates the spinner message
func (p *ProgressSpinner) UpdateMessage(message string) {
	p.message = message
}

// spin runs the spinner animation
func (p *ProgressSpinner) spin() {
	spinners := []rune{'|', '/', '-', '\\'}
	i := 0
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-p.done:
			return
		case <-ticker.C:
			if atomic.LoadInt32(p.isRunning) == 1 {
				fmt.Printf("\r[%s%c%s] %s", CyanColor, spinners[i%len(spinners)], Reset, p.message)
				i++
			}
		}
	}
}

// ProgressCounter displays progress with counters
type ProgressCounter struct {
	message     string
	current     *int64
	total       int64
	found       *int32
	isRunning   *int32
	done        chan bool
}

// NewCounter creates a new progress counter
func NewCounter(message string, total int64) *ProgressCounter {
	var current int64
	var found int32
	var running int32
	
	return &ProgressCounter{
		message:   message,
		current:   &current,
		total:     total,
		found:     &found,
		isRunning: &running,
		done:      make(chan bool),
	}
}

// Start begins the progress counter display
func (p *ProgressCounter) Start() {
	atomic.StoreInt32(p.isRunning, 1)
	go p.display()
}

// Stop stops the progress counter
func (p *ProgressCounter) Stop() {
	if atomic.LoadInt32(p.isRunning) == 1 {
		atomic.StoreInt32(p.isRunning, 0)
		p.done <- true
		// Clear the line
		fmt.Printf("\r%80s\r", "")
	}
}

// Increment increments the current counter
func (p *ProgressCounter) Increment() {
	atomic.AddInt64(p.current, 1)
}

// IncrementFound increments the found counter
func (p *ProgressCounter) IncrementFound() {
	atomic.AddInt32(p.found, 1)
}

// GetCurrent returns the current count
func (p *ProgressCounter) GetCurrent() int64 {
	return atomic.LoadInt64(p.current)
}

// GetFound returns the found count
func (p *ProgressCounter) GetFound() int32 {
	return atomic.LoadInt32(p.found)
}

// display runs the progress display
func (p *ProgressCounter) display() {
	spinners := []rune{'|', '/', '-', '\\'}
	i := 0
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-p.done:
			return
		case <-ticker.C:
			if atomic.LoadInt32(p.isRunning) == 1 {
				currentChecked := atomic.LoadInt64(p.current)
				currentFound := atomic.LoadInt32(p.found)
				fmt.Printf("\r[%s%c%s] %s [Checked: %d/%d | Found: %d]", 
					CyanColor, spinners[i%len(spinners)], Reset, p.message, currentChecked, p.total, currentFound)
				i++
			}
		}
	}
}

// SimpleProgress displays a simple percentage progress bar
func SimpleProgress(current, total int64, message string) {
	if total == 0 {
		return
	}
	
	percentage := float64(current) / float64(total) * 100
	fmt.Printf("\r%s: %.1f%% (%d/%d)", message, percentage, current, total)
}