package crypto

import (
	"fmt"
	"github.com/obscuronet/go-obscuro/integration/common/testlog"
	"golang.org/x/exp/rand"
	"golang.org/x/exp/slices"
	"testing"
)

const N = 1000

func Test_dataEncryptionServiceImpl_Encrypt(t1 *testing.T) {
	d := NewDataEncryptionService(testlog.Logger())

	plaintext := make([]byte, 512*1024)
	_, err := rand.Read(plaintext)
	if err != nil {
		panic(err)
	}

	var gains = float64(0)
	for i := 0; i < N; i++ {
		encryptedBlob, err := d.Encrypt(plaintext)
		if err != nil {
			panic(err)
		}

		byteFrequency := make(map[byte]int)
		for _, b := range encryptedBlob {
			byteFrequency[b] = byteFrequency[b] + 1
		}

		freqs := make([]int, 0)
		for _, v := range byteFrequency {
			freqs = append(freqs, v)
		}

		max := slices.Max(freqs)
		avg := len(encryptedBlob) / 256

		// calculate the gain as the difference between the numbers of the best byte compared to 0
		gain := float64(max-byteFrequency[byte(0)]) * 100 / float64(avg)
		gains += gain
	}

	fmt.Printf("GAIN: %f", gains/N)

}
