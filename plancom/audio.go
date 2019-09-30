package plancom

import (
	"github.com/atolVerderben/tentsuyu"
)

func loadAudio() *tentsuyu.AudioPlayer {
	audioPlayer, _ := tentsuyu.NewAudioPlayer()
	//audioPlayer.AddSongFromBytes("BGM", assets.BG_MUSIC_WAV)
	audioPlayer.AddSongFromFile("BGM", "assets/audio/FallingOrgan.mp3") //Space Station LOOP

	//Sound Effects
	exVolume := 0.6
	//audioPlayer.AddSoundEffectFromBytes("explosion1", assets.EXPLOSION1_WAV, exVolume)
	//audioPlayer.AddSoundEffectFromBytes("explosion2", assets.EXPLOSION2_WAV, exVolume)
	//audioPlayer.AddSoundEffectFromBytes("explosion3", assets.EXPLOSION3_WAV, exVolume)
	//audioPlayer.AddSoundEffectFromBytes("explosion4", assets.EXPLOSION4_WAV, exVolume)
	//audioPlayer.AddSoundEffectFromBytes("explosion5", assets.EXPLOSION5_WAV, exVolume)
	audioPlayer.AddSoundEffectFromFile("explosion1", "assets/audio/explosion/retro_explosion_02.ogg", exVolume)
	audioPlayer.AddSoundEffectFromFile("explosion2", "assets/audio/explosion/retro_explosion_03.ogg", exVolume)
	audioPlayer.AddSoundEffectFromFile("explosion3", "assets/audio/explosion/retro_explosion_05.ogg", exVolume)
	audioPlayer.AddSoundEffectFromFile("explosion4", "assets/audio/explosion/retro_explosion_04.ogg", exVolume)
	audioPlayer.AddSoundEffectFromFile("explosion5", "assets/audio/explosion/retro_die_02.ogg", exVolume)

	audioPlayer.AddSoundEffectFromFile("overheat", "assets/audio/beep_02.ogg", 1.0)

	//audioPlayer.AddSoundEffectFromBytes("hit-planet", assets.HIT_PLANET_WAV, 0.4)
	audioPlayer.AddSoundEffectFromFile("hit-planet", "assets/audio/explosion/retro_die_03.ogg", 0.8)

	//audioPlayer.AddSoundEffectFromBytes("blaster", assets.BLASTER_WAV, 0.3)
	//audioPlayer.AddSoundEffectFromBytes("blaster2", assets.BLASTER2_WAV, 0.3)
	audioPlayer.AddSoundEffectFromFile("blaster", "assets/audio/retro_laser_01.ogg", 0.7)
	audioPlayer.AddSoundEffectFromFile("blaster2", "assets/audio/retro_laser_02.ogg", 0.7)
	return audioPlayer
}
