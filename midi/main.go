package main

import (
	"fmt"
	"math"
	"time"

	"github.com/hajimehoshi/oto/v2"
)

const (
	sampleRate          = 44100
	quarterNoteDuration = 0.5 // Duration of a quarter note in seconds
)

// Function to generate a sine wave
func generateSineWave(frequency float64, duration float64) []byte {
	samples := int(sampleRate * duration)
	buf := make([]byte, samples*2) // 2 bytes per sample (16 bits)
	for i := 0; i < samples; i++ {
		sample := int16(math.Sin(2*math.Pi*frequency*float64(i)/sampleRate) * 32767)
		buf[i*2] = byte(sample & 0xff)
		buf[i*2+1] = byte((sample >> 8) & 0xff)
	}
	return buf
}

// Function to play a tone
func playTone(frequency float64, duration float64) error {
	ctx, err, _ := oto.NewContext(sampleRate, 2, 2, 4096)
	if err != nil {
		return fmt.Errorf("error creating audio context: %v", err)
	}
	defer ctx.Close()

	// Generate and play the sine wave
	wave := generateSineWave(frequency, duration)
	player := ctx.NewPlayer()
	defer player.Close()

	_, err = player.Write(wave)
	if err != nil {
		return fmt.Errorf("error playing tone: %v", err)
	}

	time.Sleep(time.Duration(duration * float64(time.Second))) // Wait for the duration of the note
	return nil
}

func noteToFrequency(note string) float64 {
	frequencies := map[string]float64{
		"C5": 523.25,
		"D5": 587.33,
		"E5": 659.25,
		"F5": 698.46,
		"G5": 783.99,
		"A5": 880.00,
		"B5": 987.77,
	}
	return frequencies[note]
}

func main() {
	// Example MIDI notation (to be defined correctly)
	parts := []string{"C5q", "D5q", "E5q", "F5q", "G5q"}
	// Play the notes
	for _, part := range parts[3:] { // Ignore the first 3 elements
		note := part[:2]    // Note (e.g., C5)
		duration := part[2] // Duration (e.g., q)
		// Determine the frequency of the note
		frequency := noteToFrequency(note)

		// Define the duration of the note
		var noteDuration float64
		switch duration {
		case 'q': // quarter note
			noteDuration = quarterNoteDuration
		// You can add more cases for other durations, if necessary
		default:
			noteDuration = quarterNoteDuration // default
		}

		// Play the note
		err := playTone(frequency, noteDuration)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("Playback completed.")
}
