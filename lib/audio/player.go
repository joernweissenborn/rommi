package audio

type Player interface {
	Play(audio Audio) error
}
