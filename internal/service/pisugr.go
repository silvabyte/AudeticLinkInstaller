package service

//https://github.com/PiSugar/PiSugar/wiki/PiSugar-Power-Manager-(Software)#commands-of-controlling-pisugar-server-systemd-service

// // install pisugar
// // todo: this is interactive, so we need to run it manually
// cfg.Progress.UpdateMessage("Downloading Pisugar Power Manager...")
// cmd = exec.Command("wget", "https://cdn.pisugar.com/release/pisugar-power-manager.sh")
// if err := execWithLogging(cmd); err != nil {
// 	return fmt.Errorf("failed to install pisugar power manager: %w", err)
// }

// cfg.Progress.UpdateMessage("Installing Pisugar Power Manager...")
// cmd = exec.Command("bash", "pisugar-power-manager.sh", "-c", "release")
// if err := execWithLogging(cmd); err != nil {
// 	return fmt.Errorf("failed to install pisugar power manager: %w", err)
// }

// //echo "set_button_enable double 1" | nc -q 0 127.0.0.1 8423
// cfg.Progress.UpdateMessage("Setting enable double button on Pisugar Power Manager...")
// cmd = exec.Command("sh", "-c", "echo 'set_button_enable double 1' | nc -q 0 127.0.0.1 8423")
// if err := execWithLogging(cmd); err != nil {
// 	return fmt.Errorf("failed to set pisugar power manager: %w", err)
// }

// //echo "set_button_shell double curl -X POST http://0.0.0.0:8481/record/toggle" | nc -q 0 127.0.0.1 8423
// cfg.Progress.UpdateMessage("Setting toggle record button on Pisugar Power Manager...")
// cmd = exec.Command("sh", "-c", "echo 'set_button_shell double curl -X POST http://0.0.0.0:8481/record/toggle' | nc -q 0 127.0.0.1 8423")
// if err := execWithLogging(cmd); err != nil {
// 	return fmt.Errorf("failed to set pisugar power manager: %w", err)
// }
