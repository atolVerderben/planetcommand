package plancom

import (
	assets "github.com/atolVerderben/planetcommand/plancom/internal"
	"github.com/atolVerderben/tentsuyu"
)

func loadAudio() *tentsuyu.AudioPlayer {
	audioPlayer, _ := tentsuyu.NewAudioPlayer()
	audioPlayer.AddSongFromBytes("BGM", assets.FALLINGORGAN_MP3)
	//audioPlayer.AddSongFromFile("BGM", "assets/audio/FallingOrgan.mp3") //Space Station LOOP

	//Sound Effects
	exVolume := 0.6
	audioPlayer.AddSoundEffectFromBytes("explosion1", assets.RETRO_EXPLOSION_02_OGG, exVolume)
	audioPlayer.AddSoundEffectFromBytes("explosion2", assets.RETRO_EXPLOSION_03_OGG, exVolume)
	audioPlayer.AddSoundEffectFromBytes("explosion3", assets.RETRO_EXPLOSION_05_OGG, exVolume)
	audioPlayer.AddSoundEffectFromBytes("explosion4", assets.RETRO_EXPLOSION_04_OGG, exVolume)
	audioPlayer.AddSoundEffectFromBytes("explosion5", assets.RETRO_DIE_02_OGG, exVolume)
	//audioPlayer.AddSoundEffectFromFile("explosion1", "assets/audio/explosion/retro_explosion_02.ogg", exVolume)
	//audioPlayer.AddSoundEffectFromFile("explosion2", "assets/audio/explosion/retro_explosion_03.ogg", exVolume)
	//audioPlayer.AddSoundEffectFromFile("explosion3", "assets/audio/explosion/retro_explosion_05.ogg", exVolume)
	//audioPlayer.AddSoundEffectFromFile("explosion4", "assets/audio/explosion/retro_explosion_04.ogg", exVolume)
	//audioPlayer.AddSoundEffectFromFile("explosion5", "assets/audio/explosion/retro_die_02.ogg", exVolume)

	//audioPlayer.AddSoundEffectFromFile("overheat", "assets/audio/beep_02.ogg", 1.0)
	audioPlayer.AddSoundEffectFromBytes("overheat", assets.BEEP_02_OGG, 1.0)
	audioPlayer.AddSoundEffectFromBytes("hit-planet", assets.RETRO_DIE_03_OGG, 0.8)
	//audioPlayer.AddSoundEffectFromFile("hit-planet", "assets/audio/explosion/retro_die_03.ogg", 0.8)

	audioPlayer.AddSoundEffectFromBytes("blaster", assets.RETRO_LASER_01_OGG, 0.7)
	audioPlayer.AddSoundEffectFromBytes("blaster2", assets.RETRO_LASER_02_OGG, 0.7)
	//audioPlayer.AddSoundEffectFromFile("blaster", "assets/audio/retro_laser_01.ogg", 0.7)
	//audioPlayer.AddSoundEffectFromFile("blaster2", "assets/audio/retro_laser_02.ogg", 0.7)

	audioPlayer.AddSoundEffectFromFile("ufo-in", "assets/audio/enemies/ufo/synth_misc_03.ogg", 0.8)
	audioPlayer.AddSoundEffectFromFile("ufo-hit", "assets/audio/enemies/ufo/retro_coin_02.ogg", 1.0)
	audioPlayer.AddSoundEffectFromFile("ufo-die", "assets/audio/enemies/ufo/power_up_04.ogg", 1.0)
	return audioPlayer
}
