package cmd

import (
	"fmt"
	"os"
	"rommi/lib/audio/pa"
	"rommi/lib/train"

	"github.com/ThingiverseIO/console"
	"github.com/ThingiverseIO/uuid"
	"github.com/spf13/cobra"
)

func init() {
	trainCmd.AddCommand(trainAudioCmd)
	trainAudioCmd.AddCommand(trainAudioShowCmd)
	trainAudioCmd.AddCommand(trainAudioPlayCmd)
	trainAudioCmd.AddCommand(trainAudioRecordCmd)
	trainAudioCmd.AddCommand(trainAudioInteractiveCmd)
}

var trainAudioCmd = &cobra.Command{
	Use:   "audio",
	Short: "Create, update and modify rommi's training audio records.",
}

var trainAudioShowCmd = &cobra.Command{
	Use:   "show [MODELPATH]",
	Short: "Show all audio in the training database.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		path, err := os.Getwd()
		if err != nil {
			return
		}
		if len(args) != 0 {
			path = args[0]
		}
		cmd.Println("Opening Model Directory:", path)
		t, err := train.Open(path)
		if err != nil {
			return
		}
		cmd.Println("Success")
		cmd.Println("")
		printAudio(t)
		return
	},
}

var trainAudioPlayCmd = &cobra.Command{
	Use:   "play [MODELPATH]",
	Short: "Record audio for the given UUID.",
	RunE:  runTrainAudioPlay,
}

func runTrainAudioPlay(cmd *cobra.Command, args []string) error {
	path, err := os.Getwd()
	if err != nil {
		console.Println(err)
		return nil
	}
	if len(args) > 0 {
		path = args[0]
	}
	console.Println("Opening Model Directory:", path)
	t, err := train.Open(path)
	if err != nil {
		console.Println(err)
		return nil
	}
	console.Println("Success")

	speaker, aborted := selectSpeaker(t)
	if aborted {
		return nil
	}
	condition, aborted := selectCondition(t, speaker)
	if aborted {
		return nil
	}

	ids := t.AvailableAudio(speaker, condition)
	ops := map[string]interface{}{}
	for _, id := range ids {
		s, _ := t.GetSentence(id)
		ops[fmt.Sprintf("%s:%s", id, s)] = id
	}

	_, iid, aborted := console.AskOptionValue("Please select a sentence", ops)
	if aborted {
		return nil
	}
	id := iid.(uuid.UUID)

	a, err := t.GetAudio(speaker, condition, id)
	if err != nil {
		console.Println(err)
		return nil
	}

	cmd.Println("Length of recorded WAVE is:", a.Size())
	cmd.Println("SampleRate of recorded WAVE is:", a.Rate())
	cmd.Println("Nr of channels of recorded WAVE is:", a.Channels())
	if !console.AskYesOrNo("Play File?", true) {
		return nil
	}
	player, err := pa.NewPlayer()
	if err != nil {
		console.Println(err)
		return nil
	}
	console.Println("Playing")
	player.Play(a)
	if err != nil {
		console.Println(err)
		return nil
	}
	return nil
}

var trainAudioRecordCmd = &cobra.Command{
	Use:   "record [MODELPATH]",
	Short: "Record audio for training.",
	RunE:  runTrainAudioRecord,
}

func runTrainAudioRecord(cmd *cobra.Command, args []string) error {
	path, err := os.Getwd()
	if err != nil {
		console.Println(err)
		return nil
	}
	if len(args) > 0 {
		path = args[0]
	}
	console.Println("Opening Model Directory:", path)
	t, err := train.Open(path)
	if err != nil {
		console.Println(err)
		return nil
	}
	console.Println("Success")

	speaker, aborted := selectSpeaker(t)
	if aborted {
		return nil
	}
	condition, aborted := selectCondition(t, speaker)
	if aborted {
		return nil
	}
	sentence, id, aborted := selectSentence(t)
	if aborted {
		return nil
	}
	triggerword := t.GetTriggerWord()

	console.Println("SUMMARY")
	console.Println("=======")
	console.Println("Speaker:", speaker)
	console.Println("Condition:", condition)
	console.Println("Triggerword:", triggerword)
	console.Println("SentenceId:", id)
	console.Println("Sentence:", sentence)
	if t.HasAudio(speaker, condition, id) {
		console.Println("Warning: Audio Already Exist. Exiting Audio Will Be Overwritten.")
	}
	if !console.AskYesOrNo("Proceed?", false) {
		return nil
	}
	wav, err := recordSentence(t, id)
	if err != nil {
		console.Println(err)
		return nil
	}
	cmd.Println("Length of recorded WAVE is:", len(wav))

	if console.AskYesOrNo("Save?", true) {
		console.Println("Saving")
		err = t.SaveAudio(speaker, condition, id, wav)
		if err != nil {
			console.Println("Error saving audio:", err)
			return nil
		}
		cmd.Println("Success")
	}
	return nil
}

var trainAudioInteractiveCmd = &cobra.Command{
	Use:   "interactive [MODELPATH]",
	Short: "Record audio for training interactively.",
	RunE:  runTrainAudioInteractive,
}

func runTrainAudioInteractive(cmd *cobra.Command, args []string) error {
	path, err := os.Getwd()
	if err != nil {
		console.Println(err)
		return nil
	}
	if len(args) > 0 {
		path = args[0]
	}
	console.Println("Opening Model Directory:", path)
	t, err := train.Open(path)
	if err != nil {
		console.Println(err)
		return nil
	}
	console.Println("Success")

	var (
		aborted   bool
		speaker   string
		condition string
	)
	for !aborted {
		if speaker == "" {
			speaker, aborted = selectSpeaker(t)
			if aborted {
				return nil
			}
			continue
		}
		if condition == "" {
			condition, aborted = selectCondition(t, speaker)
			if aborted {
				return nil
			}
			continue
		}

		console.Printline()
		console.Println("selected speaker: ", speaker)
		console.Println("selected condition: ", condition)
		console.Printline()

		var action string
		action, aborted = console.AskOption("select action", "play", "record")
		console.Printline()
		if aborted {
			return nil
		}

		switch action {
		case "play":
			action, aborted = console.AskOption("what to play", "all", "single")
			console.Printline()
			if aborted {
				aborted = false
				continue
			}
			switch action {
			case "single":
				s, id, aborted := selectSentence(t)
				if aborted {
					return nil
				}
				console.Printf("Playing sentence %s: '%s'\n", id, s)
				a, err := t.GetAudio(speaker, condition, id)
				if err != nil {
					console.Println(err)
					return nil
				}

				cmd.Println("Length of recorded WAVE is:", a.Size())
				cmd.Println("SampleRate of recorded WAVE is:", a.Rate())
				cmd.Println("Nr of channels of recorded WAVE is:", a.Channels())
				if !console.AskYesOrNo("Play?", true) {
					continue
				}
				player, err := pa.NewPlayer()
				if err != nil {
					console.Println(err)
					return nil
				}
				console.Println("Playing")
				player.Play(a)
				if err != nil {
					console.Println(err)
					return nil
				}
			case "all":
				ids := t.GetAllIds()
				all := len(ids)
				for i, id := range ids {
					if i != 0 {
						if !console.AskYesOrNo("Proceed?", true) {
							break
						}
						s, _ := t.GetSentence(id)
						console.Printf("Playing sentence %d of %d  %s: '%s'\n", i+1, all, id, s)
						a, err := t.GetAudio(speaker, condition, id)
						if err != nil {
							console.Println(err)
							return nil
						}

						cmd.Println("Length of recorded WAVE is:", a.Size())
						cmd.Println("SampleRate of recorded WAVE is:", a.Rate())
						cmd.Println("Nr of channels of recorded WAVE is:", a.Channels())
						if !console.AskYesOrNo("Play?", true) {
							continue
						}
						player, err := pa.NewPlayer()
						if err != nil {
							console.Println(err)
							return nil
						}
						console.Println("Playing")
						player.Play(a)
						if err != nil {
							console.Println(err)
							return nil
						}
					}
				}
			}
		case "record":
			action, aborted = console.AskOption("what to record", "all", "single", "missing")
			console.Printline()
			if aborted {
				aborted = false
				continue
			}
			switch action {
			case "single":
				s, id, aborted := selectSentence(t)
				if aborted {
					return nil
				}
				console.Printf("Recording for sentence %s: '%s'\n", id, s)
				if t.HasAudio(speaker, condition, id) {
					console.Println("Warning: Audio Already Exist. Exiting Audio Will Be Overwritten.")
				}
				if !console.AskYesOrNo("Record?", true) {
					continue
				}
				wav, err := recordSentence(t, id)
				if err != nil {
					console.Println(err)
					return nil
				}
				cmd.Println("Length of recorded WAVE is:", len(wav))

				if console.AskYesOrNo("Save?", true) {
					console.Println("Saving")
					err = t.SaveAudio(speaker, condition, id, wav)
					if err != nil {
						console.Println("Error saving audio:", err)
						return nil
					}
					cmd.Println("Success")
				}
			case "all":
				ids := t.GetAllIds()
				all := len(ids)
				for i, id := range ids {
					if i != 0 {
						if !console.AskYesOrNo("Proceed?", true) {
							break
						}
					}
					s, _ := t.GetSentence(id)
					console.Printf("Recording for sentence %d of %d  %s: '%s'\n", i+1, all, id, s)
					if t.HasAudio(speaker, condition, id) {
						console.Println("Warning: Audio Already Exist. Exiting Audio Will Be Overwritten.")
					}
					if !console.AskYesOrNo("Record?", true) {
						continue
					}
					wav, err := recordSentence(t, id)
					if err != nil {
						console.Println(err)
						return nil
					}
					cmd.Println("Length of recorded WAVE is:", len(wav))

					if console.AskYesOrNo("Save?", true) {
						console.Println("Saving")
						err = t.SaveAudio(speaker, condition, id, wav)
						if err != nil {
							console.Println("Error saving audio:", err)
							return nil
						}
						cmd.Println("Success")
					}
				}
				console.Println("All Sentences Recorded")
			case "missing":
				ids := t.GetAllIds()
				var mids []uuid.UUID
				for _, id := range ids {
					if !t.HasAudio(speaker, condition, id) {
						mids = append(mids, id)
					}
				}
				all := len(mids)
				for i, id := range mids {
					if i != 0 {
						if !console.AskYesOrNo("Proceed?", true) {
							break
						}
					}
					s, _ := t.GetSentence(id)
					console.Printf("Recording for sentence %d of %d  %s: '%s'\n", i+1, all, id, s)
					if t.HasAudio(speaker, condition, id) {
						console.Println("Warning: Audio Already Exist. Exiting Audio Will Be Overwritten.")
					}
					if !console.AskYesOrNo("Record?", true) {
						continue
					}
					wav, err := recordSentence(t, id)
					if err != nil {
						console.Println(err)
						return nil
					}
					cmd.Println("Length of recorded WAVE is:", len(wav))

					if console.AskYesOrNo("Save?", true) {
						console.Println("Saving")
						err = t.SaveAudio(speaker, condition, id, wav)
						if err != nil {
							console.Println("Error saving audio:", err)
							return nil
						}
						cmd.Println("Success")
					}
				}
				console.Println("All Sentences Recorded")
			}
		}

	}
	return nil
}

func selectSentence(t *train.Train) (sentence string, id uuid.UUID, aborted bool) {

	ops := map[string]interface{}{}
	for _, id := range t.GetAllIds() {
		s,_:= t.GetSentence(id)
		ops[s] = id
	}
	sentence, iid, aborted := console.AskOptionValue("Please select a sentence to record", ops)
	if aborted {
		return
	}
	id = iid.(uuid.UUID)
	return
}

func selectSpeaker(t *train.Train) (speaker string, aborted bool) {
	speakers := t.GetSpeakers()
	speaker = "New"
	if len(speakers) != 0 {
		if speaker, aborted = console.AskOption(
			"Please select a speaker or create a new one", append(speakers, "New Speaker")...); aborted {
			return
		}
	}

	if speaker == "New" {
		speaker = console.AskString("Please enter name for new speaker (empty to abort): ")
	}
	if speaker == "" {
		aborted = true
	}
	return
}

func selectCondition(t *train.Train, speaker string) (condition string, aborted bool) {
	conditions := t.GetConditions(speaker)
	condition = "New"
	if len(conditions) != 0 {
		if condition, aborted = console.AskOption(
			"Please select a condition or create a new one", append(conditions, "New Speaker")...); aborted {
			return
		}
	}

	if condition == "New" {
		condition = console.AskString("Please enter name for new condition (empty to abort): ")
	}
	if condition == "" {
		aborted = true
	}
	return
}

func recordSentence(t *train.Train, id uuid.UUID) (wav []byte, err error) {
	triggerword := t.GetTriggerWord()
	sentence,_ := t.GetSentence(id)
	console.Println("Please Record The Following Sentence:")
	console.Println("")
	console.Printf("\t%s %s\n", triggerword, sentence)
	console.Println("")
	wav, err = record()
	return
}

func printAudio(t *train.Train) {
	// if len(db) == 0 {
	//         fmt.Println("No Audio In Database")
	//         return
	// }
	fmt.Println("Audio in Database")
	fmt.Println("=====================")
	fmt.Println("")
	for _, speaker := range t.GetSpeakers() {
			fmt.Printf("Speaker '%s':\n", speaker)
		for _, condition := range t.GetConditions(speaker) {
				fmt.Printf("\tCondition '%s':\n", condition)
				ids := t.AvailableAudio(speaker, condition)
				for _, id := range ids {
					s, ok := t.GetSentence(id)
					if !ok {
						s = "Unknown Sentence"
					}
					fmt.Printf("\t\tID: %s Sentence: %s\n", id, s)
				}
			}
		}
	}
