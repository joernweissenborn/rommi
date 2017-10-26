package audio

type Recorder interface {
	OpenRecordStream() (audio *AudioStream, err error)
	StartRecording() (err error)
	StopRecording() (err error)
}
