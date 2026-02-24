package genid

import "github.com/matoous/go-nanoid/v2"

const alphabet = "123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const nanoIDLength = 10

func GenerateNanoID() string {
	id, _ := gonanoid.Generate(alphabet, nanoIDLength)
	return id
}
