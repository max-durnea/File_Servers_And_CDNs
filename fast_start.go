package main
import 
(
	"os/exec"
)

func processVideoForFastStart(filePath string) (string,error){
	output_filePath := filePath + ".processing"
	cmd:=exec.Command("ffmpeg","-i",filePath,"-c","copy","-movflags","faststart","-f","mp4",output_filePath)
	cmd.Run()
	return output_filePath,nil
}